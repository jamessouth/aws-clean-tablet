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
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
	Color  string `dynamodbav:"color,omitempty"`
}

type gameout struct {
	No      string     `json:"no"`
	Leader  string     `json:"leader,omitempty"`
	Players playerList `json:"players"`
}

type connin struct {
	GSI1SK string `json:"gsi1sk"`
}

// type gamein struct {
// 	No      string            `json:"no"`
// 	Leader  string            `json:"leader,omitempty"`
// 	Players map[string]player `json:"players"`
// }

type insertConnPayload struct {
	Games []gameout `json:"games"`
	Type  string    `json:"type"`
}

type modifyConnPayload struct {
	Ingame      string `json:"ingame"`
	Leadertoken string `json:"leadertoken"`
	Type        string `json:"type"`
}

type insertGamePayload struct {
	Games gameout `json:"games"`
	Type  string  `json:"type"`
}

func (p playerList) sortByName() playerList {
	sort.Slice(p, func(i, j int) bool {
		return p[i].Name > p[j].Name
	})

	return p
}

func (pm playerMap) getPlayersSlice() (res playerList) {
	for _, v := range pm {
		res = append(res, v)
	}

	return
}

type gamesList []gamein
type connsList []connin
type playerList []player
type playerMap map[string]player

func (gl gamesList) mapGames() (res []gameout) {
	for _, g := range gl {
		res = append(res, gameout{
			No:      g.Sk,
			Leader:  g.Leader,
			Players: g.Players.getPlayersSlice().sortByName(),
		})
	}

	return
}

func FromDynamoDBEventAVMap(from map[string]events.DynamoDBAttributeValue) (to map[string]types.AttributeValue, err error) {
	to = make(map[string]types.AttributeValue, len(from))
	for field, value := range from {
		to[field], err = FromDynamoDBEventAV(value)
		if err != nil {
			return nil, err
		}
	}

	return to, nil
}

func FromDynamoDBEventAV(from events.DynamoDBAttributeValue) (types.AttributeValue, error) {
	switch from.DataType() {

	case events.DataTypeBoolean:
		return &types.AttributeValueMemberBOOL{Value: from.Boolean()}, nil

	case events.DataTypeString:
		return &types.AttributeValueMemberS{Value: from.String()}, nil

	case events.DataTypeMap:
		values, err := FromDynamoDBEventAVMap(from.Map())
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberM{Value: values}, nil

	default:
		return nil, fmt.Errorf("unknown AttributeValue union member, %T", from)
	}
}

type answer struct {
	PlayerID, Answer string
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
	Pk       string    `dynamodbav:"pk"`
	Sk       string    `dynamodbav:"sk"`
	Starting bool      `dynamodbav:"starting"`
	Leader   string    `dynamodbav:"leader"`
	Loading  bool      `dynamodbav:"loading"`
	Players  playerMap `dynamodbav:"players"`
	Answers  []answer  `dynamodbav:"answers"`
}

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	// fmt.Println("reqqqq", req)
	for _, rec := range req.Records {
		tableName := strings.Split(rec.EventSourceArn, "/")[1]
		ni := rec.Change.NewImage
		fmt.Printf("%s: %+v\n", "new db ni", ni)

		item, err := FromDynamoDBEventAVMap(ni)
		if err != nil {
			fmt.Println("item unmarshal err", err)
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
				fmt.Println("eppppppppppppppppp", ep)
				return ep, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
			config.WithEndpointResolver(customResolver),
		)
		if err != nil {
			fmt.Println("cfg err")
		}

		apigwsvc := apigatewaymanagementapi.NewFromConfig(cfg)
		ddbsvc := dynamodb.NewFromConfig(cfg)

		recType := item["pk"].(*types.AttributeValueMemberS).Value[:4]

		if recType == "CONN" {

			var connRecord connItem
			err = attributevalue.UnmarshalMap(item, &connRecord)
			if err != nil {
				fmt.Println("item unmarshal err", err)
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

				}

				var games gamesList
				err = attributevalue.UnmarshalListOfMaps(gamesResults.Items, &games)
				if err != nil {
					fmt.Println("query unmarshal err", err)
				}

				payload, err := json.Marshal(insertConnPayload{
					Games: games.mapGames(),
					Type:  "games",
				})
				if err != nil {
					fmt.Println("error marshalling", err)
				}

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connRecord.GSI1SK), Data: payload}

				_, err = apigwsvc.PostToConnection(ctx, &conn)
				if err != nil {

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

				}
			} else if rec.EventName == dynamodbstreams.OperationTypeModify && !connRecord.Playing {

				payload, err := json.Marshal(modifyConnPayload{
					Ingame:      connRecord.Game,
					Leadertoken: connRecord.Sk + "_" + connRecord.GSI1SK,
					Type:        "user",
				})
				if err != nil {
					fmt.Println("error marshalling payload", err)
				}

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connRecord.GSI1SK), Data: payload}

				_, err = apigwsvc.PostToConnection(ctx, &conn)
				if err != nil {

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

				}

			}

		} else if recType == "GAME" {

			var gameRecord gamein
			err = attributevalue.UnmarshalMap(item, &gameRecord)
			if err != nil {
				fmt.Println("item unmarshal err", err)
			}

			fmt.Printf("%s%+v\n", "gammmmme ", gameRecord)

			if rec.EventName == dynamodbstreams.OperationTypeInsert {

				payload, err := json.Marshal(insertGamePayload{
					Games: gameout{
						No:      gameRecord.Sk,
						Leader:  gameRecord.Leader,
						Players: gameRecord.Players.getPlayersSlice(),
					},
					Type: "games",
				})
				if err != nil {
					fmt.Println("error marshalling payload", err)
				}

				connsParams := dynamodb.QueryInput{
					TableName:              aws.String(tableName),
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

				connsResults, err := ddbsvc.Query(ctx, &connsParams)
				if err != nil {

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

				}

				var conns connsList
				err = attributevalue.UnmarshalListOfMaps(connsResults.Items, &conns)
				if err != nil {
					fmt.Println("query unmarshal err", err)
				}

				for _, v := range conns {

					conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.GSI1SK), Data: payload}

					_, err = apigwsvc.PostToConnection(ctx, &conn)
					if err != nil {

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

					}

				}

			} else if rec.EventName == dynamodbstreams.OperationTypeModify {
				if strings.HasPrefix(item["pk"].String(), "CONN") {

				} else if strings.HasPrefix(item["pk"].String(), "GAME") {

					if item["loading"].Boolean() {

						if len(item["answers"].List()) == 0 || len(item["answers"].List()) == 8 {

						}

					}

				} else {
					fmt.Println("other modify", item)
				}
			}

		} else {
			fmt.Println("other record type", item)
		}

	}

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
