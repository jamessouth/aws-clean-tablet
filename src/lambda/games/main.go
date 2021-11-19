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

type player struct {
	Name   string `json:"name"`
	ConnID string `json:"connid"`
	Ready  bool   `json:"ready"`
	Color  string `json:"color,omitempty"`
	Score  int    `json:"score"`
	Answer answer `json:"answer"`
}

type gameout struct {
	Pk       string `json:"pk,omitempty"`
	No       string `json:"no"`
	Ready    bool   `json:"ready"`
	Starting bool   `json:"starting,omitempty"`
	Loading  bool   `json:"loading,omitempty"`

	Players playerList `json:"players"`
}

type connin struct {
	GSI1SK string `json:"gsi1sk"`
}

type insertConnPayload struct {
	ListGames gameOutList `json:"listGms"`
	ConnID    string      `json:"connID"`
}

type modifyConnPayload struct {
	ModConnGm string `json:"modConn"`
}

type insertGamePayload struct {
	AddGame gameout `json:"addGame"`
}

type modifyListGamePayload struct {
	ModListGame gameout `json:"mdLstGm"`
}

type modifyLiveGamePayload struct {
	ModLiveGame gameout `json:"mdLveGm"`
}

type removeGamePayload struct {
	RemoveGame gameout `json:"rmvGame"`
}

type lessFunc func(p1, p2 *player) int

type multiSorter struct {
	players []player
	less    []lessFunc
}

func (ms *multiSorter) Sort(players []player) {
	ms.players = players
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.players)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.players[i], ms.players[j] = ms.players[j], ms.players[i]
}

func (ms *multiSorter) Less(i, j int) bool {
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

var name = func(a, b *player) int {
	if a.Name > b.Name {
		return -1
	}

	return 1
}

var score = func(a, b *player) int {
	if a.Score < b.Score {
		return -1
	}
	if a.Score > b.Score {
		return 1
	}

	return 0
}

func (p playerList) sort(fs ...lessFunc) playerList {
	OrderedBy(fs...).Sort(p)

	return p
}

func (pm playerMap) getPlayersSlice() (res playerList) {
	res = make(playerList, 0)

	for _, v := range pm {
		res = append(res, v)
	}

	return
}

type gameInList []gamein
type gameOutList []gameout
type connsList []connin
type playerList []player
type playerMap map[string]player

func (gl gameInList) mapGames() (res gameOutList) {
	res = make(gameOutList, 0)

	for _, g := range gl {
		res = append(res, gameout{
			Pk:       "",
			No:       g.Sk,
			Ready:    g.Ready,
			Starting: false,
			Loading:  false,

			Players: g.Players.getPlayersSlice().sort(name),
		})
	}

	return
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

type answer struct {
	PlayerID string `json:"playerid"`
	Answer   string `json:"answer"`
}

type connItem struct {
	Pk      string `dynamodbav:"pk"`      //'CONN#' + uuid
	Sk      string `dynamodbav:"sk"`      //name
	Game    string `dynamodbav:"game"`    //game no or blank
	Playing bool   `dynamodbav:"playing"` //playing or not
	GSI1PK  string `dynamodbav:"GSI1PK"`  //'CONN'
	GSI1SK  string `dynamodbav:"GSI1SK"`  //conn id
}

type gamein struct {
	Pk       string `dynamodbav:"pk"`
	Sk       string `dynamodbav:"sk"`
	Starting bool   `dynamodbav:"starting"`
	Ready    bool   `dynamodbav:"ready"`
	Loading  bool   `dynamodbav:"loading"`

	Players      playerMap `dynamodbav:"players"`
	AnswersCount int       `dynamodbav:"answersCount"`
	SendToFront  bool      `dynamodbav:"sendToFront"`
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

		if rec.EventName == dynamodbstreams.OperationTypeRemove {
			oi := rec.Change.OldImage
			fmt.Printf("%s: %+v\n", "old db oi", oi)

			// item, err := FromDynamoDBEventAVMap(oi)
			// if err != nil {
			// 	return callErr(err)
			// }

			// var connRecord connItem
			// err = attributevalue.UnmarshalMap(item, &connRecord)
			// if err != nil {
			// 	return callErr(err)
			// }

			// fmt.Printf("%s%+v\n", "connrecord ", connRecord)

			// _, err = apigwsvc.DeleteConnection(ctx, &apigatewaymanagementapi.DeleteConnectionInput{
			// 	ConnectionId: aws.String(connRecord.GSI1SK),
			// })
			// if err != nil {
			// 	return callErr(err)
			// }

			return getReturnValue(http.StatusOK), nil
		}

		ddbcfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
		)
		if err != nil {
			return callErr(err)
		}

		ddbsvc := dynamodb.NewFromConfig(ddbcfg)

		recType := item["pk"].(*types.AttributeValueMemberS).Value[:4]

		if recType == "CONN" {

			var connRecord connItem
			err = attributevalue.UnmarshalMap(item, &connRecord)
			if err != nil {
				return callErr(err)
			}

			fmt.Printf("%s%+v\n", "connrecord ", connRecord)

			if rec.EventName == dynamodbstreams.OperationTypeInsert {

				gamesParams := dynamodb.QueryInput{
					TableName:              aws.String(tableName),
					ScanIndexForward:       aws.Bool(false),
					KeyConditionExpression: aws.String("pk = :gm"),
					FilterExpression:       aws.String("#ST = :st"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":gm": &types.AttributeValueMemberS{Value: "GAME"},
						":st": &types.AttributeValueMemberBOOL{Value: false},
					},
					ExpressionAttributeNames: map[string]string{
						"#ST": "starting",
					},
				}

				gamesResults, err := ddbsvc.Query(ctx, &gamesParams)
				if err != nil {
					return callErr(err)
				}

				var games gameInList
				err = attributevalue.UnmarshalListOfMaps(gamesResults.Items, &games)
				if err != nil {
					return callErr(err)
				}

				// Tipe:  "listGames",
				payload, err := json.Marshal(insertConnPayload{
					ListGames: games.mapGames(),
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

			} else if rec.EventName == dynamodbstreams.OperationTypeModify && !connRecord.Playing {

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

			} else {
				fmt.Println("remove conn ", connRecord)
			}

		} else if recType == "GAME" {

			var gameRecord gamein
			err = attributevalue.UnmarshalMap(item, &gameRecord)
			if err != nil {
				return callErr(err)
			}

			fmt.Printf("%s%+v\n", "gammmmme ", gameRecord)

			// if len(gameRecord.Answers) > 0 && len(gameRecord.Answers) < len(gameRecord.Players) {
			// 	return getReturnValue(http.StatusOK), nil
			// }

			if rec.EventName == dynamodbstreams.OperationTypeInsert {

				err = sendGamesToConns(ctx, ddbsvc, apigwsvc, gameRecord, tableName, "add")
				if err != nil {
					return callErr(err)
				}

			} else if rec.EventName == dynamodbstreams.OperationTypeModify {

				if gameRecord.SendToFront {

					if gameRecord.Loading {
						gp := modifyGamePayload{
							ModGame: gameout{
								Pk:       "",
								No:       gameRecord.Sk,
								Ready:    gameRecord.Ready,
								Starting: gameRecord.Starting,
								Loading:  gameRecord.Loading,

								Players: gameRecord.Players.getPlayersSlice().sort(score, name),
							},
						}

						payload, err := json.Marshal(gp)
						if err != nil {
							return callErr(err)
						}

						for _, v := range gp.ModGame.Players {

							conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

							_, err = apigwsvc.PostToConnection(ctx, &conn)
							if err != nil {
								return callErr(err)
							}

						}

					} else {
						if gameRecord.Starting {

							err = sendGamesToConns(ctx, ddbsvc, apigwsvc, gameRecord, tableName, "rem")
							if err != nil {
								return callErr(err)
							}

						} else {
							err = sendGamesToConns(ctx, ddbsvc, apigwsvc, gameRecord, tableName, "mod")
							if err != nil {
								return callErr(err)
							}

						}
					}
				}

			} else {
				fmt.Println("remove game ", gameRecord)
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

func getGamePayload(g gamein, opt string) (payload []byte, err error) {

	pl := gameout{
		Pk:       "",
		No:       g.Sk,
		Ready:    g.Ready,
		Starting: g.Starting,
		Loading:  g.Loading,

		Players: g.Players.getPlayersSlice().sort(name),
	}

	if opt == "add" {
		payload, err = json.Marshal(insertGamePayload{
			AddGame: pl,
		})
	} else if opt == "mod" {
		payload, err = json.Marshal(modifyGamePayload{
			ModGame: pl,
		})
	} else if opt == "rem" {
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
			":cn": &types.AttributeValueMemberS{Value: "CONN"},
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

func sendGamesToConns(ctx context.Context, ddbsvc *dynamodb.Client, apigwsvc *apigatewaymanagementapi.Client, gr gamein, tn, opt string) error {

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
