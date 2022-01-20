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
	Color   string `dynamodbav:"color"`   //player's color
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
		pls := g.Players.getListPlayersSlice()
		pls.sortByName()
		res = append(res, toFEListGame{
			No:      g.Sk,
			Ready:   g.Ready,
			Players: pls,
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

type livePlayer struct {
	Name        string `json:"name"`
	ConnID      string `json:"connid"`
	Color       string `json:"color"`
	Score       int    `json:"score"`
	Index       int    `json:"index"`
	Answer      string `json:"answer"`
	HasAnswered bool   `json:"hasAnswered"`
}

type livePlayerList []livePlayer

type liveGame struct {
	Sk           string         `json,dynamodbav:"sk"`
	Players      livePlayerList `json,dynamodbav:"players"`
	CurrentWord  string         `json,dynamodbav:"currentWord"`
	PreviousWord string         `json,dynamodbav:"previousWord"`
	AnswersCount int            `json,dynamodbav:"answersCount"`
	ShowAnswers  bool           `json,dynamodbav:"showAnswers"`
}

type insertConnPayload struct {
	ListGames toFEListGameList `json:"listGms"`
	ConnID    string           `json:"connID"`
}

type modifyConnPayload struct {
	ModConnGm string `json:"modConn"`
	Color     string `json:"color"`
}

type insertGamePayload struct {
	AddGame toFEListGame `json:"addGame"`
}

type modifyListGamePayload struct {
	ModListGame toFEListGame `json:"mdLstGm"`
}

type modifyLiveGamePayload struct {
	ModLiveGame liveGame
}

// https://go.dev/play/p/CvniMWPoLKG
func (p modifyLiveGamePayload) MarshalJSON() ([]byte, error) {
	if p.ModLiveGame.AnswersCount == len(p.ModLiveGame.Players) {
		return []byte(`null`), nil
	}

	if p.ModLiveGame.AnswersCount > 0 {
		for i, pl := range p.ModLiveGame.Players {
			if pl.HasAnswered {
				pl.Answer = ""
				p.ModLiveGame.Players[i] = pl
			}
		}
	}

	m, err := json.Marshal(p.ModLiveGame)
	if err != nil {
		return m, err
	}

	return []byte(fmt.Sprintf("{%q:%s}", "mdLveGm", m)), nil
}

type removeGamePayload struct {
	RemoveGame toFEListGame `json:"rmvGame"`
}

func (players listPlayerList) sortByName() {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})
}

func (players livePlayerList) sortByAnswerThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Answer != players[j].Answer:
			return players[i].Answer < players[j].Answer
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) sortByScoreThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Score != players[j].Score:
			return players[i].Score > players[j].Score
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) addIndex() livePlayerList {
	for i, p := range players {
		p.Index = i
		players[i] = p
	}

	return players
}

func FromDynamoDBEventAVMap(m map[string]events.DynamoDBAttributeValue) (res map[string]types.AttributeValue, err error) {
	// fmt.Println("av map: ", m)
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
	// fmt.Println("av list: ", l)
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
	// fmt.Println("av type: ", av, av.DataType())
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
	for i, rec := range req.Records {
		// fmt.Println("rekkkk", req.Records, len(req.Records))
		fmt.Println("reccc: ", i, rec, len(req.Records))

		tableName := strings.Split(rec.EventSourceArn, "/")[1]

		var ni map[string]events.DynamoDBAttributeValue

		if rec.EventName == dynamodbstreams.OperationTypeRemove {
			ni = rec.Change.OldImage
		} else {
			ni = rec.Change.NewImage
		}
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

				// if connRecord.Playing {
				// 	fmt.Println("mod conn playing in a game", connRecord)

				// } else {

				// Tipe:   "modifyConn",
				payload, err := json.Marshal(modifyConnPayload{
					ModConnGm: connRecord.Game,
					Color:     connRecord.Color,
				})

				if err != nil {
					return callErr(err)
				}

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connRecord.GSI1SK), Data: payload}

				_, err = apigwsvc.PostToConnection(ctx, &conn)
				if err != nil {
					return callErr(err)
				}
				// }

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

			payload, err := getGamePayload(listGameRecord, opt)
			if err != nil {
				return callErr(err)
			}

			conns, err := unmarshalConns(ctx, ddbsvc, tableName)
			if err != nil {
				return callErr(err)
			}

			for _, v := range conns {

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.GSI1SK), Data: payload}

				_, err := apigwsvc.PostToConnection(ctx, &conn)
				if err != nil {
					return callErr(err)
				}

			}

		} else if recType == "LIVEGME" {

			if rec.EventName == dynamodbstreams.OperationTypeInsert || rec.EventName == dynamodbstreams.OperationTypeModify {

				var gameRecord liveGame
				err = attributevalue.UnmarshalMap(item, &gameRecord)
				if err != nil {
					return callErr(err)
				}

				fmt.Printf("%s%+v\n", "live gammmmme ", gameRecord)

				// if gameRecord.SendToFront {
				pls := gameRecord.Players.addIndex()

				// if gameRecord.AnswersCount == len(gameRecord.Players) {
				// 	return getReturnValue(http.StatusOK), nil
				// } else

				// if gameRecord.AnswersCount > 0 {
				// 	pls.sortByScoreThenName()
				// } else {
				// 	pls.sortByAnswerThenName()
				// }

				if gameRecord.ShowAnswers {
					pls.sortByAnswerThenName()
				} else {
					pls.sortByScoreThenName()
				}

				gp := modifyLiveGamePayload{
					ModLiveGame: liveGame{
						Sk:           gameRecord.Sk,
						Players:      pls,
						CurrentWord:  gameRecord.CurrentWord,
						PreviousWord: gameRecord.PreviousWord,
						AnswersCount: gameRecord.AnswersCount,
						ShowAnswers:  gameRecord.ShowAnswers,
					},
				}

				payload, err := json.Marshal(gp)
				if err != nil {
					return callErr(err)
				}

				for _, v := range pls {

					conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

					_, err = apigwsvc.PostToConnection(ctx, &conn)
					if err != nil {
						return callErr(err)
					}

				}

				// }

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
		pls := g.Players.getListPlayersSlice()
		pls.sortByName()
		pl := toFEListGame{
			No:      g.Sk,
			Ready:   g.Ready,
			Players: pls,
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
