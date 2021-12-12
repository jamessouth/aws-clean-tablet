package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	"github.com/aws/smithy-go"
)

type connItem struct {
	Pk      string `dynamodbav:"pk"`      //'CONNECT#' + uuid
	Sk      string `dynamodbav:"sk"`      //name
	Game    string `dynamodbav:"game"`    //game no or blank
	Playing bool   `dynamodbav:"playing"` //playing or not
	GSI1PK  string `dynamodbav:"GSI1PK"`  //'CONNECT'
	GSI1SK  string `dynamodbav:"GSI1SK"`  //conn id
}

type connin struct {
	GSI1SK string
}

type connsList []connin

//-----------------------------------------------------------------------------

type listPlayer struct {
	Name   string `json:"name" dynamodbav:"name"`
	ConnID string `json:"connid" dynamodbav:"connid"`
	Ready  bool   `json:"ready" dynamodbav:"ready"`
}
type listPlayerMap map[string]listPlayer

type fromDBListGame struct {
	Pk      string        `dynamodbav:"pk"` //'LISTGME'
	Sk      string        `dynamodbav:"sk"` //no
	Ready   bool          `dynamodbav:"ready"`
	Players listPlayerMap `dynamodbav:"players"`
}

type listPlayerList []listPlayer

type fromDBListGameList []fromDBListGame

func (gl fromDBListGameList) mapListGames() (res toFEListGameList) {
	res = make(toFEListGameList, 0)

	for _, g := range gl {
		res = append(res, toFEListGame{
			No:      g.Sk,
			Ready:   g.Ready,
			Players: g.Players.getListPlayersSlice().sort(nameList),
		})
	}

	return
}

type toFEListGameList []toFEListGame

type toFEListGame struct {
	No      string         `json:"no"`
	Ready   bool           `json:"ready"`
	Players listPlayerList `json:"players"`
}

func (pm listPlayerMap) getListPlayersSlice() (res listPlayerList) {
	res = make(listPlayerList, 0)

	for _, v := range pm {
		res = append(res, v)
	}

	return
}

//-------------------------------------------------------------------------------

type answer struct {
	PlayerID string `json:"playerid"`
	Answer   string `json:"answer"`
}

type livePlayer struct {
	Name        string `json:"name"`
	ConnID      string `json:"connid"`
	Color       string `json:"color"`
	Score       int    `json:"score"`
	Answer      answer `json:"answer"`
	HasAnswered bool   `json:"hasAnswered"`
}

type livePlayerList []livePlayer

type livePlayerMap map[string]livePlayer

type toFELiveGame struct {
	No           string         `json:"no"`
	CurrentWord  string         `json:"currentWord"`
	PreviousWord string         `json:"previousWord"`
	Players      livePlayerList `json:"players"`
	AnswersCount int            `json:"answersCount"`
}

type fromDBLiveGame struct {
	Sk           string        `dynamodbav:"sk"`
	Players      livePlayerMap `dynamodbav:"players"`
	CurrentWord  string        `dynamodbav:"currentWord"`
	PreviousWord string        `dynamodbav:"previousWord"`
	AnswersCount int           `dynamodbav:"answersCount"`
	SendToFront  bool          `dynamodbav:"sendToFront"`
}

type fromDBLiveGameLite struct {
	Sk      string        `dynamodbav:"sk"`
	Players livePlayerMap `dynamodbav:"players"`
	// CurrentWord  string        `dynamodbav:"currentWord"`
	// PreviousWord string        `dynamodbav:"previousWord"`
	// AnswersCount int           `dynamodbav:"answersCount"`
	// SendToFront  bool          `dynamodbav:"sendToFront"`
}

func (pm livePlayerMap) getLivePlayersSlice() (res livePlayerList) {
	res = make(livePlayerList, 0)

	for _, v := range pm {
		res = append(res, v)
	}

	return
}

type insertConnPayload struct {
	ListGames toFEListGameList `json:"listGms"`
	ConnID    string           `json:"connID"`
}

type modifyConnPayload struct {
	ModConnGm string `json:"modConn"`
}

type insertGamePayload struct {
	AddGame toFEListGame `json:"addGame"`
}

type modifyListGamePayload struct {
	ModListGame toFEListGame `json:"mdLstGm"`
}

type modifyLiveGamePayload struct {
	ModLiveGame toFELiveGame `json:"mdLveGm"`
}

func (p modifyLiveGamePayload) MarshalJSON() ([]byte, error) {
	if p.ModLiveGame.AnswersCount == len(p.ModLiveGame.Players) {
		return []byte("null"), nil
	}
	if p.ModLiveGame.AnswersCount > 0 && p.ModLiveGame.AnswersCount < len(p.ModLiveGame.Players) {
		for _, p := range p.ModLiveGame.Players {
			if p.Answer.Answer != "" {
				p.Answer.Answer = ""
				p.HasAnswered = true
			}
		}
	}

	return json.Marshal(p)
}

type removeGamePayload struct {
	RemoveGame toFEListGame `json:"rmvGame"`
}

type lessFuncList func(p1, p2 *listPlayer) int

type multiSorterList struct {
	players listPlayerList
	less    []lessFuncList
}

func (ms *multiSorterList) Sort(players listPlayerList) {
	ms.players = players
	sort.Sort(ms)
}

func OrderedByList(less ...lessFuncList) *multiSorterList {
	return &multiSorterList{
		less: less,
	}
}

func (ms *multiSorterList) Len() int {
	return len(ms.players)
}

func (ms *multiSorterList) Swap(i, j int) {
	ms.players[i], ms.players[j] = ms.players[j], ms.players[i]
}

func (ms *multiSorterList) Less(i, j int) bool {
	for _, k := range ms.less {
		switch k(&ms.players[i], &ms.players[j]) {
		case 1:
			return true
		case -1:
			return false
		}
	}

	return true
}

func (p listPlayerList) sort(fs ...lessFuncList) listPlayerList {
	OrderedByList(fs...).Sort(p)

	return p
}

var nameList = func(a, b *listPlayer) int {
	if a.Name > b.Name {
		return -1
	}

	return 1
}

type lessFuncLive func(p1, p2 *livePlayer) int

type multiSorterLive struct {
	players livePlayerList
	less    []lessFuncLive
}

func (ms *multiSorterLive) Sort(players livePlayerList) {
	ms.players = players
	sort.Sort(ms)
}

func OrderedByLive(less ...lessFuncLive) *multiSorterLive {
	return &multiSorterLive{
		less: less,
	}
}

func (ms *multiSorterLive) Len() int {
	return len(ms.players)
}

func (ms *multiSorterLive) Swap(i, j int) {
	ms.players[i], ms.players[j] = ms.players[j], ms.players[i]
}

func (ms *multiSorterLive) Less(i, j int) bool {
	for _, k := range ms.less {
		switch k(&ms.players[i], &ms.players[j]) {
		case 1:
			return true
		case -1:
			return false
		}
	}

	return true
}

func (p livePlayerList) sort(fs ...lessFuncLive) livePlayerList {
	OrderedByLive(fs...).Sort(p)

	return p
}

var nameLive = func(a, b *livePlayer) int {
	if a.Name > b.Name {
		return -1
	}

	return 1
}

var score = func(a, b *livePlayer) int {
	if a.Score < b.Score {
		return -1
	}
	if a.Score > b.Score {
		return 1
	}

	return 0
}

func FromDynamoDBEventAVMap(m map[string]events.DynamoDBAttributeValue) (res map[string]types.AttributeValue, err error) {
	fmt.Println("av map: ", m)
	res = make(map[string]types.AttributeValue, len(m))

	for k, v := range m {
		res[k], err = FromDynamoDBEventAV(v)
		if err != nil {
			return nil, err
		}
	}

	return
}

func FromDynamoDBEventAVList(l []events.DynamoDBAttributeValue) (res []types.AttributeValue, err error) {
	fmt.Println("av list: ", l)
	res = make([]types.AttributeValue, len(l))

	for i, v := range l {
		res[i], err = FromDynamoDBEventAV(v)
		if err != nil {
			return nil, err
		}
	}

	return
}

func FromDynamoDBEventAV(av events.DynamoDBAttributeValue) (types.AttributeValue, error) {
	fmt.Println("av type: ", av, av.DataType())
	switch av.DataType() {

	case events.DataTypeBoolean: // 1
		return &types.AttributeValueMemberBOOL{Value: av.Boolean()}, nil

	case events.DataTypeList: // 3
		values, err := FromDynamoDBEventAVList(av.List())
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberL{Value: values}, nil

	case events.DataTypeMap: // 4
		values, err := FromDynamoDBEventAVMap(av.Map())
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberM{Value: values}, nil

	case events.DataTypeNumber: // 5
		return &types.AttributeValueMemberN{Value: av.Number()}, nil

	case events.DataTypeNull: // 7
		return &types.AttributeValueMemberNULL{Value: av.IsNull()}, nil

	case events.DataTypeString: // 8
		return &types.AttributeValueMemberS{Value: av.String()}, nil

	default:
		return nil, fmt.Errorf("unknown AttributeValue union member, %T", av)
	}
}

func getReturnValue(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}
}

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	// fmt.Println("reqqqq", req)
	for _, rec := range req.Records {
		fmt.Println("reccc: ", rec)

		tableName := strings.Split(rec.EventSourceArn, "/")[1]
		ni := rec.Change.NewImage
		fmt.Printf("%s: %+v\n", "new db ni", ni)

		item, err := FromDynamoDBEventAVMap(ni)
		if err != nil {
			return callErr(err)
		}

		apiid, ok := os.LookupEnv("CT_APIID")
		if !ok {
			panic(fmt.Sprintf("%v", "can't find api id"))
		}

		stage, ok := os.LookupEnv("CT_STAGE")
		if !ok {
			panic(fmt.Sprintf("%v", "can't find stage"))
		}
		endpoint := "https://" + apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage

		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == apigatewaymanagementapi.ServiceID && region == rec.AWSRegion {
				ep := aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: rec.AWSRegion,
				}

				return ep, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

		apigwcfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
			config.WithEndpointResolver(customResolver),
		)
		if err != nil {
			return callErr(err)
		}

		apigwsvc := apigatewaymanagementapi.NewFromConfig(apigwcfg)

		ddbcfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
		)
		if err != nil {
			return callErr(err)
		}

		ddbsvc := dynamodb.NewFromConfig(ddbcfg)

		recType := item["pk"].(*types.AttributeValueMemberS).Value[:7]

		if recType == "CONNECT" {

			var connRecord connItem
			err = attributevalue.UnmarshalMap(item, &connRecord)
			if err != nil {
				return callErr(err)
			}

			fmt.Printf("%s%+v\n", "connrecord ", connRecord)

			if rec.EventName == dynamodbstreams.OperationTypeInsert {

				listGamesParams := dynamodb.QueryInput{
					TableName:              aws.String(tableName),
					ScanIndexForward:       aws.Bool(false),
					KeyConditionExpression: aws.String("pk = :gm"),
					// FilterExpression:       aws.String("#ST = :st"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":gm": &types.AttributeValueMemberS{Value: "LISTGME"},
						// ":st": &types.AttributeValueMemberBOOL{Value: false},
					},
				}

				listGamesResults, err := ddbsvc.Query(ctx, &listGamesParams)
				if err != nil {
					return callErr(err)
				}

				var listGames fromDBListGameList
				err = attributevalue.UnmarshalListOfMaps(listGamesResults.Items, &listGames)
				if err != nil {
					return callErr(err)
				}

				// Tipe:  "listGames",
				payload, err := json.Marshal(insertConnPayload{
					ListGames: listGames.mapListGames(),
					ConnID:    connRecord.GSI1SK,
				})
				if err != nil {
					return callErr(err)
				}

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connRecord.GSI1SK), Data: payload}

				_, err = apigwsvc.PostToConnection(ctx, &conn)
				if err != nil {
					return callErr(err)
				}

			} else if rec.EventName == dynamodbstreams.OperationTypeModify {

				if connRecord.Playing {
					fmt.Println("mod conn playing in a game", connRecord)

				} else {

					// Tipe:   "modifyConn",
					payload, err := json.Marshal(modifyConnPayload{
						ModConnGm: connRecord.Game,
					})

					if err != nil {
						return callErr(err)
					}

					conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connRecord.GSI1SK), Data: payload}

					_, err = apigwsvc.PostToConnection(ctx, &conn)
					if err != nil {
						return callErr(err)
					}
				}

			} else {
				oi := rec.Change.OldImage
				fmt.Printf("%s: %+v\n", "remove conn oi", oi)
			}

		} else if recType == "LISTGME" {

			var listGameRecord fromDBListGame
			err = attributevalue.UnmarshalMap(item, &listGameRecord)
			if err != nil {
				return callErr(err)
			}

			fmt.Printf("%s%+v\n", "list gammmmme ", listGameRecord)
			var opt string
			switch rec.EventName {
			case dynamodbstreams.OperationTypeInsert:
				opt = "add"
			case dynamodbstreams.OperationTypeModify:
				opt = "mod"
			default:
				oi := rec.Change.OldImage
				fmt.Printf("%s: %+v\n", "remove list game oi", oi)
				opt = "rem"
			}

			err = sendGamesToConns(ctx, ddbsvc, apigwsvc, listGameRecord, tableName, opt)
			if err != nil {
				return callErr(err)
			}

		} else if recType == "LIVEGME" {

			if rec.EventName == dynamodbstreams.OperationTypeInsert {

				var gameRecordLite fromDBLiveGameLite
				err = attributevalue.UnmarshalMap(item, &gameRecordLite)
				if err != nil {
					return callErr(err)
				}

				fmt.Printf("%s%+v\n", "live gammmmme lite ", gameRecordLite)

				// send to front

			} else if rec.EventName == dynamodbstreams.OperationTypeModify {

				var gameRecord fromDBLiveGame
				err = attributevalue.UnmarshalMap(item, &gameRecord)
				if err != nil {
					return callErr(err)
				}

				fmt.Printf("%s%+v\n", "live gammmmme ", gameRecord)

				if gameRecord.SendToFront {

					// if gameRecord.AnswersCount > 0 {

					// }
					// else if gameRecord.AnswersCount == len(gameRecord.Players) {

					// }

					gp := modifyLiveGamePayload{
						ModLiveGame: toFELiveGame{
							No:           gameRecord.Sk,
							CurrentWord:  gameRecord.CurrentWord,
							PreviousWord: gameRecord.PreviousWord,
							Players:      gameRecord.Players.getLivePlayersSlice().sort(score, nameLive),
							AnswersCount: gameRecord.AnswersCount,
						},
					}

					payload, err := json.Marshal(gp)
					if err != nil {
						return callErr(err)
					}

					for _, v := range gp.ModLiveGame.Players {

						conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

						_, err = apigwsvc.PostToConnection(ctx, &conn)
						if err != nil {
							return callErr(err)
						}

					}

				}

			} else {
				oi := rec.Change.OldImage
				fmt.Printf("%s: %+v\n", "remove live game oi", oi)
			}

		} else {
			fmt.Println("other record type", item)
		}

	}

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
}

func getGamePayload(g fromDBListGame, opt string) (payload []byte, err error) {

	if opt == "add" {
		pl := toFEListGame{
			No:      g.Sk,
			Players: g.Players.getListPlayersSlice(),
		}
		payload, err = json.Marshal(insertGamePayload{
			AddGame: pl,
		})
	} else if opt == "mod" {
		pl := toFEListGame{
			No:      g.Sk,
			Ready:   g.Ready,
			Players: g.Players.getListPlayersSlice().sort(nameList),
		}
		payload, err = json.Marshal(modifyListGamePayload{
			ModListGame: pl,
		})
	} else if opt == "rem" {
		pl := toFEListGame{
			No: g.Sk,
		}
		payload, err = json.Marshal(removeGamePayload{
			RemoveGame: pl,
		})

	} else {
		return nil, errors.New("invalid payload option provided")
	}

	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	return payload, nil
}

func getConnsInput(tn string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(tn),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :cn"),
		FilterExpression:       aws.String("#PL = :f"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cn": &types.AttributeValueMemberS{Value: "CONNECT"},
			":f":  &types.AttributeValueMemberBOOL{Value: false},
		},
		ExpressionAttributeNames: map[string]string{
			"#PL": "playing",
		},
	}
}

func getConnsItems(ctx context.Context, ddbsvc *dynamodb.Client, tn string) ([]map[string]types.AttributeValue, error) {

	connsParams := getConnsInput(tn)

	connsResults, err := ddbsvc.Query(ctx, connsParams)
	if err != nil {
		return nil, fmt.Errorf("query for conns failed: %w", err)
	}

	return connsResults.Items, nil
}

func unmarshalConns(ctx context.Context, ddbsvc *dynamodb.Client, tn string) (connsList, error) {

	ci, err := getConnsItems(ctx, ddbsvc, tn)
	if err != nil {
		return nil, fmt.Errorf("could not get conns: %w", err)
	}

	var conns connsList
	err = attributevalue.UnmarshalListOfMaps(ci, &conns)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal conns: %w", err)
	}

	return conns, nil

}

func sendGamesToConns(ctx context.Context, ddbsvc *dynamodb.Client, apigwsvc *apigatewaymanagementapi.Client, gr fromDBListGame, tn, opt string) error {

	payload, err := getGamePayload(gr, opt)
	if err != nil {
		return fmt.Errorf("could not get payload: %w", err)
	}

	conns, err := unmarshalConns(ctx, ddbsvc, tn)
	if err != nil {
		return fmt.Errorf("could not get connections: %w", err)
	}

	for _, v := range conns {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.GSI1SK), Data: payload}

		_, err := apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return fmt.Errorf("could not post to connection %s: %w", v.GSI1SK, err)
		}

	}

	return nil
}

func callErr(err error) (events.APIGatewayProxyResponse, error) {

	var intServErr *types.InternalServerError
	if errors.As(err, &intServErr) {
		fmt.Printf("get item error, %v",
			intServErr.ErrorMessage())
	}

	// To get any API error
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("db error, Code: %v, Message: %v",
			apiErr.ErrorCode(), apiErr.ErrorMessage())
	}

	return getReturnValue(http.StatusBadRequest), err

}
