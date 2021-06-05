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

type gamein struct {
	No      string            `json:"no"`
	Leader  string            `json:"leader,omitempty"`
	Players map[string]player `json:"players"`
}

type insertConnPayload struct {
	Games []gameout `json:"games"`
	Type  string    `json:"type"`
}

// type insertGamePayload struct {
// 	Games gameout `json:"games"`
// 	Type  string  `json:"type"`
// }

func (p playerList) sortByName() playerList {
	sort.Slice(p, func(i, j int) bool {
		return p[i].Name > p[j].Name
	})

	return p
}

func getPlayersSlice(m map[string]player) (res playerList) {
	for _, v := range m {
		res = append(res, v)
	}

	return
}

type gamesList []gamein
type playerList []player

func (gl gamesList) mapGames() (res []gameout) {
	for _, g := range gl {
		res = append(res, gameout{
			No:      g.No,
			Leader:  g.Leader,
			Players: getPlayersSlice(g.Players).sortByName(),
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

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	// fmt.Println("reqqqq", req)
	for _, rec := range req.Records {
		tableName := strings.Split(rec.EventSourceArn, "/")[1]
		item := rec.Change.NewImage
		fmt.Printf("%s: %+v\n", "new db item", item)

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

		if rec.EventName == dynamodbstreams.OperationTypeInsert {
			if strings.HasPrefix(item["pk"].String(), "CONN") {

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
				// if err != nil {
				// 	panic(fmt.Sprintf("failed to marshal query input, %v", err))
				// }
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

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(item["GSI1SK"].String()), Data: payload}

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

			} else if strings.HasPrefix(item["pk"].String(), "GAME") {

				// var players map[string]player

				game, _ := FromDynamoDBEventAVMap(item)

				var gamein gamein
				err = attributevalue.UnmarshalMap(game, &gamein)
				if err != nil {
					fmt.Println("item unmarshal err", err)
				}

				fmt.Printf("%s%+v\n", "gammmmme ", gamein)

				// err = attributevalue.UnmarshalMap(, &players)
				// if err != nil {
				// 	fmt.Println("query unmarshal err", err)
				// }

				// payload, err := json.Marshal(insertGamePayload{
				// 	Games: gameout{
				// 		No:      item["sk"].String(),
				// 		Leader:  item["leader"].String(),
				// 		Players: getPlayersSlice(),
				// 	},
				// 	Type: "games",
				// })
				// if err != nil {
				// 	fmt.Println("error marshalling", err)
				// }

			} else {
				fmt.Println("other game")
			}
			// }
			// }
		}

	}
	// r := bufio.NewReader(&buffer)
	// log := make([]byte, 0, 1024)
	// for {
	// 	n, err := io.ReadFull(r, log[:cap(log)])
	// 	log = log[:n]
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != io.ErrUnexpectedEOF {
	// 			fmt.Fprintln(os.Stderr, err)
	// 			break
	// 		}
	// 	}

	// 	fmt.Printf("read %d bytes: ", n)

	// 	fmt.Println(string(log))
	// }

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
