package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

const maxPlayersPerGame = 8

// Key holds values to be put in db
type Key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// Player holds values to be put in db
type Player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready,omitempty"`
	Color  string `dynamodbav:"color,omitempty"`
}

type game struct {
	Pk       string            `dynamodbav:"pk"`
	Sk       string            `dynamodbav:"sk"`
	Starting bool              `dynamodbav:"starting"`
	Leader   string            `dynamodbav:"leader"`
	Loading  bool              `dynamodbav:"loading"`
	Ready    bool              `dynamodbav:"ready"`
	Players  map[string]Player `dynamodbav:"players"`
}

type body struct {
	Game, Type string
	Value      bool
}

// ConnItemAttrs holds vals for db
type ConnItemAttrs struct {
	Game string `dynamodbav:":g"`
	Zero *int   `dynamodbav:":z,omitempty"`
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	// fmt.Println("lobbbbbby", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.NewFromConfig(cfg)
	// svc2 := lamb.NewFromConfig(cfg)

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	id := auth["principalId"].(string)
	name := auth["username"].(string)

	var body body
	var gameno string

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	if body.Game != "new" {
		gameno = body.Game
	} else {
		gameno = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	gameItemKey, err := attributevalue.MarshalMap(Key{
		Pk: "GAME",
		Sk: gameno,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal gik Record, %v", err))
	}

	connItemKey, err := attributevalue.MarshalMap(Key{
		Pk: "CONN#" + id,
		Sk: name,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal cik Record 3, %v", err))
	}

	marshalledFalse, err := attributevalue.Marshal(false)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal false, %v", err))
	}

	marshalledTrue, err := attributevalue.Marshal(true)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal true, %v", err))
	}

	if body.Type == "join" {

		zero := 0
		connAttrs, err := attributevalue.MarshalMap(ConnItemAttrs{
			Game: gameno,
			Zero: &zero,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 4, %v", err))
		}

		player := Player{
			Name:   name,
			ConnID: req.RequestContext.ConnectionID,
			Ready:  false,
			Color:  "",
		}

		marshalledPlayersMap, err := attributevalue.Marshal(map[string]Player{
			id: player,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal map Record 22, %v", err))
		}

		marshalledMaxPlayers, err := attributevalue.Marshal(maxPlayersPerGame)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal max Record 22, %v", err))
		}

		marshalledPlayer, err := attributevalue.Marshal(player)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal indiv Record 22, %v", err))
		}

		updateConnInput := types.Update{
			Key:                 connItemKey,
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("size (#IG) = :z"),
			ExpressionAttributeNames: map[string]string{
				"#IG": "game",
			},
			ExpressionAttributeValues: connAttrs,
			UpdateExpression:          aws.String("SET #IG = :g"),
		}

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_exists(#PL) AND size (#PL) < :m"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
							"#ID": id,
							// "#ST": "starting",#ST = :f
							// "#LO": "loading",#LO = :f
							"#RD": "ready",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":m": marshalledMaxPlayers,
							":f": marshalledFalse,
							":p": marshalledPlayer,
						},
						UpdateExpression: aws.String("SET #PL.#ID = :p, #RD = :f"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		callErr(err)

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_not_exists(#PL)"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
							"#ST": "starting",
							"#LO": "loading",
							"#LE": "leader",
							"#RD": "ready",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":p": marshalledPlayersMap,
							":e": &types.AttributeValueMemberS{Value: ""},
							":f": marshalledFalse,
						},
						UpdateExpression: aws.String("SET #PL = :p, #RD = :f, #ST = :f, #LO = :f, #LE = :e"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		callErr(err)

	} else if body.Type == "leave" {

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			Key:       gameItemKey,
			TableName: aws.String(tableName),

			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
			},

			UpdateExpression: aws.String("REMOVE #PL.#ID"),
		})
		callErr(err)

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			Key:       connItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#IG": "game",
			},

			ExpressionAttributeValues: map[string]types.AttributeValue{

				":g": &types.AttributeValueMemberS{Value: ""},
			},
			UpdateExpression: aws.String("SET #IG = :g"),
		})
		callErr(err)

		ui3, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			Key:                 gameItemKey,
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("#LE = :c"),
			ExpressionAttributeNames: map[string]string{
				"#LE": "leader",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{

				":c": &types.AttributeValueMemberS{Value: name + req.RequestContext.ConnectionID},
				":e": &types.AttributeValueMemberS{Value: ""},
			},
			UpdateExpression: aws.String("SET #LE = :e"),
			ReturnValues:     types.ReturnValueAllNew,
		})
		callErr(err)

		var game game
		err = attributevalue.UnmarshalMap(ui3.Attributes, &game)
		if err != nil {
			fmt.Println("joingame leave unmarshal err", err)
		}

		if len(game.Players) > 2 {

			if game.Ready {

				if game.Leader == "" {

					for k, v := range game.Players {

						fmt.Printf("%s, %v, %+v", "ui", k, v)

						_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
							Key:       gameItemKey,
							TableName: aws.String(tableName),
							ExpressionAttributeNames: map[string]string{
								"#LE": "leader",
							},
							ExpressionAttributeValues: map[string]types.AttributeValue{
								":l": &types.AttributeValueMemberS{Value: v.Name + v.ConnID},
							},
							UpdateExpression: aws.String("SET #LE = :l"),
						})
						callErr(err)
						break

					}

				}

			} else {
				callFunction(game, gameItemKey, tableName, marshalledTrue, ctx, svc)

			}

		} else {
			_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key:       gameItemKey,
				TableName: aws.String(tableName),
				ExpressionAttributeNames: map[string]string{
					"#RD": "ready",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":f": marshalledFalse,
				},
				UpdateExpression: aws.String("SET #RD = :f"),
			})

			callErr(err)
		}

	} else if body.Type == "ready" {

		if body.Value {

			ui, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key:       gameItemKey,
				TableName: aws.String(tableName),
				ExpressionAttributeNames: map[string]string{
					"#PL": "players",
					"#ID": id,
					"#RD": "ready",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":t": marshalledTrue,
				},
				UpdateExpression: aws.String("SET #PL.#ID.#RD = :t"),
				ReturnValues:     types.ReturnValueAllNew,
			})

			callErr(err)

			var game game
			err = attributevalue.UnmarshalMap(ui.Attributes, &game)
			if err != nil {
				fmt.Println("del item unmarshal err", err)
			}
			if len(game.Players) > 2 {

				callFunction(game, gameItemKey, tableName, marshalledTrue, ctx, svc)

			}

		} else {

			_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key:       gameItemKey,
				TableName: aws.String(tableName),
				ExpressionAttributeNames: map[string]string{
					"#PL": "players",
					"#ID": id,
					"#RD": "ready",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":f": marshalledFalse,
				},
				UpdateExpression: aws.String("SET #PL.#ID.#RD = :f, #RD = :f"),
			})
			callErr(err)

			_, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

				Key:                 gameItemKey,
				TableName:           aws.String(tableName),
				ConditionExpression: aws.String("#LE = :c"),
				ExpressionAttributeNames: map[string]string{
					"#LE": "leader",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{

					":c": &types.AttributeValueMemberS{Value: name + req.RequestContext.ConnectionID},
					":e": &types.AttributeValueMemberS{Value: ""},
				},
				UpdateExpression: aws.String("SET #LE = :e"),
			})
			callErr(err)

		}

	} else {
		fmt.Println("other lobby")
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

func callFunction(gm game, gik map[string]types.AttributeValue, tn string, mt types.AttributeValue, ctx context.Context, svc *dynamodb.Client) {
	readyCount := 0

	for k, v := range gm.Players {

		fmt.Printf("%s, %v, %+v", "uicf", k, v)

		if v.Ready {
			readyCount++
			if readyCount == len(gm.Players) {
				var uii dynamodb.UpdateItemInput
				if gm.Leader == "" {
					uii = dynamodb.UpdateItemInput{
						Key:       gik,
						TableName: aws.String(tn),
						ExpressionAttributeNames: map[string]string{

							"#RD": "ready",
							"#LE": "leader",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":t": mt,
							":l": &types.AttributeValueMemberS{Value: v.Name + v.ConnID},
						},
						UpdateExpression: aws.String("SET #RD = :t, #LE = :l"),
					}
				} else {
					uii = dynamodb.UpdateItemInput{
						Key:       gik,
						TableName: aws.String(tn),
						ExpressionAttributeNames: map[string]string{
							"#RD": "ready",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":t": mt,
						},
						UpdateExpression: aws.String("SET #RD = :t"),
					}
				}
				_, err := svc.UpdateItem(ctx, &uii)

				callErr(err)
			}
		}
	}
}

func callErr(err error) {
	if err != nil {
		var transCxldErr *types.TransactionCanceledException
		if errors.As(err, &transCxldErr) {
			fmt.Printf("put item error777, %v\n",
				transCxldErr.CancellationReasons)
		}

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
