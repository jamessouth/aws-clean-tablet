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

// const maxPlayersPerGame = 8

// Key holds values to be put in db
type Key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// Player holds values to be put in db
type Player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
	Color  string `dynamodbav:"color"`
}

type body struct {
	Game, Type string
	Value      bool
}

// ConnItemAttrs holds vals for db
type ConnItemAttrs struct {
	Game string `dynamodbav:":g"`
	Zero *int   `dynamodbav:":zero,omitempty"`
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

	if body.Type == "join" {

		connItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "CONN#" + id,
			Sk: name,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 3, %v", err))
		}

		// marshaledID, err := attributevalue.Marshal(id)
		// if err != nil {
		// 	panic(fmt.Sprintf("failed to marshal marshaledID, %v", err))
		// }

		// marshaledMaxSize, err := attributevalue.Marshal(maxPlayersPerGame)
		// if err != nil {
		// 	panic(fmt.Sprintf("failed to marshal marshaledMaxSize, %v", err))
		// }

		// Players: []string{auth["username"].(string) + "#" + },
		// MaxSize: maxPlayersPerGame,
		// player, err := attributevalue.Marshal()
		// if err != nil {
		// 	panic(fmt.Sprintf("failed to marshal player, %v", err))
		// }

		z := 0
		connAttrs, err := attributevalue.MarshalMap(ConnItemAttrs{
			Game: gameno,
			Zero: &z,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 4, %v", err))
		}

		gameItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: gameno,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record, %v", err))
		}

		// GameItemAttrs holds values to be put in db
		// type GameItemAttrs struct {
		// 	Players map[string]Player `dynamodbav:"Players"`
		// 	MaxSize int               `dynamodbav:":maxsize,omitempty"`
		// 	Player  Player            `dynamodbav:":player"`
		// }

		att1, err := attributevalue.Marshal(map[string]Player{
			id: {
				Name:   name,
				ConnID: req.RequestContext.ConnectionID,
				Ready:  false,
				Color:  "",
			},
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}
		// att2, err := attributevalue.Marshal(maxPlayersPerGame)
		// if err != nil {
		// 	panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		// }
		att3, err := attributevalue.Marshal(Player{
			Name:   name,
			ConnID: req.RequestContext.ConnectionID,
			Ready:  false,
			Color:  "",
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{

			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{

						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_exists(#PL)"), //and size less than max
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
							"#ID": id,
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							// ":p": att1,
							// ":maxsize": att2,
							":player": att3,
						},

						UpdateExpression: aws.String("SET #PL.#ID = :player"),
					},
				},

				{
					Update: &types.Update{

						Key:                 connItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("size (#IG) = :zero"),
						ExpressionAttributeNames: map[string]string{
							"#IG": "game",
						},
						ExpressionAttributeValues: connAttrs,

						UpdateExpression: aws.String("SET #IG = :g"),
					},
				},
			},
			// ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		// fmt.Println("op", op)
		if err != nil {

			var intServErr *types.TransactionCanceledException
			if errors.As(err, &intServErr) {
				fmt.Printf("put item error777, %v\n",
					intServErr.CancellationReasons)
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				// fmt.Println(err.Error(), apiErr.Error())
				fmt.Printf("db error777, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{

			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{

						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_not_exists(#PL)"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
							// "#ID": id,
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":p": att1,
							// ":maxsize": att2,
							// ":player": att3,
						},

						UpdateExpression: aws.String("SET #PL = :p"),
					},
				},

				{
					Update: &types.Update{

						Key:                 connItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("size (#IG) = :zero"),
						ExpressionAttributeNames: map[string]string{
							"#IG": "game",
						},
						ExpressionAttributeValues: connAttrs,

						UpdateExpression: aws.String("SET #IG = :g"),
					},
				},
			},
			// ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		// fmt.Println("op", op)
		if err != nil {

			var intServErr *types.TransactionCanceledException
			if errors.As(err, &intServErr) {
				fmt.Printf("put item error888, %v\n",
					intServErr.CancellationReasons)
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				// fmt.Println(err.Error(), apiErr.Error())
				fmt.Printf("db error888, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}

		att4, err := attributevalue.Marshal(false)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
			ExpressionAttributeNames: map[string]string{
				// "#PL": "players",
				"#ST": "starting",
				"#RD": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":r": att4,
				// ":maxsize": att2,
				// ":player": att3,
			},

			UpdateExpression: aws.String("SET #RD = :r, #ST = :r"),
			// ReturnValues:     types.ReturnValueAllNew,
		})
		// fmt.Println("op", op)
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

	} else if body.Type == "leave" {
		connItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "CONN#" + id,
			Sk: name,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 3, %v", err))
		}

		// gameAttrs, err := attributevalue.MarshalMap(GameItemAttrs{
		// 	Players: []string{auth["username"].(string) + "#" + req.RequestContext.ConnectionID},
		// 	// MaxSize: maxPlayersPerGame,
		// })
		// if err != nil {
		// 	panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
		// }

		connAttrs, err := attributevalue.MarshalMap(ConnItemAttrs{
			Game: "",
			// Zero: 0,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 4, %v", err))
		}

		gameItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: gameno,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record, %v", err))
		}

		ui2, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
			},

			UpdateExpression: aws.String("REMOVE #PL.#ID"),
			ReturnValues:     types.ReturnValueAllNew,
		})

		if err != nil {

			var intServErr *types.InternalServerError
			if errors.As(err, &intServErr) {
				fmt.Printf("put item error 1122, %v",
					intServErr.ErrorMessage())
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("db error 1112222, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}
		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       connItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#IG": "game",
			},
			ExpressionAttributeValues: connAttrs,

			UpdateExpression: aws.String("SET #IG = :g"),
		})

		if err != nil {

			var intServErr *types.InternalServerError
			if errors.As(err, &intServErr) {
				fmt.Printf("put item error 1122, %v",
					intServErr.ErrorMessage())
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("db error 1112222, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}

		var game map[string]Player
		err = attributevalue.Unmarshal(ui2.Attributes["players"], &game)
		if err != nil {
			fmt.Println("del item unmarshal err", err)
		}
		readyCount := 0
		readyBool := false
		for k, v := range game {

			fmt.Printf("%s, %v, %+v", "ui", k, v)

			if v.Ready {
				readyCount++
			}
		}
		if len(game) > 2 && readyCount == len(game) {
			readyBool = true
		}

		att3, err := attributevalue.Marshal(readyBool)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
			ExpressionAttributeNames: map[string]string{
				// "#PL": "players",
				// "#ID": id,
				"#RD": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":r": att3,
				// ":maxsize": att2,
				// ":player": att3,
			},

			UpdateExpression: aws.String("SET #RD = :r"),
			// ReturnValues:     types.ReturnValueAllNew,
		})
		// fmt.Println("op", op)
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

	} else if body.Type == "ready" {
		gameItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: gameno,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record, %v", err))
		}

		att2, err := attributevalue.Marshal(body.Value)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}

		ui, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
				"#RD": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":r": att2,
				// ":maxsize": att2,
				// ":player": att3,
			},

			UpdateExpression: aws.String("SET #PL.#ID.#RD = :r"),
			ReturnValues:     types.ReturnValueAllNew,
		})
		// fmt.Println("op", op)
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

		var game map[string]Player
		err = attributevalue.Unmarshal(ui.Attributes["players"], &game)
		if err != nil {
			fmt.Println("del item unmarshal err", err)
		}
		readyCount := 0
		readyBool := false
		var leaderID, leaderName string

		for k, v := range game {

			fmt.Printf("%s, %v, %+v", "ui", k, v)

			if v.Ready {
				readyCount++
				if readyCount == len(game) && len(game) > 2 {
					readyBool = true
					leaderID, leaderName = k, v.Name
				}
			}
		}
		connItemKey2, err := attributevalue.MarshalMap(Key{
			Pk: "CONN#" + leaderID,
			Sk: leaderName,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record connitemkey2, %v", err))
		}
		att5, err := attributevalue.Marshal(true)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record att5, %v", err))
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       connItemKey2,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#LE": "leader",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":l": att5,
			},

			UpdateExpression: aws.String("SET #LE = :l"),
		})

		if err != nil {

			var intServErr *types.InternalServerError
			if errors.As(err, &intServErr) {
				fmt.Printf("put item error 1122, %v",
					intServErr.ErrorMessage())
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("db error 1112222, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}

		att3, err := attributevalue.Marshal(readyBool)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
			ExpressionAttributeNames: map[string]string{
				// "#PL": "players",
				// "#ID": id,
				"#RD": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":r": att3,
				// ":maxsize": att2,
				// ":player": att3,
			},

			UpdateExpression: aws.String("SET #RD = :r"),
			// ReturnValues:     types.ReturnValueAllNew,
		})
		// fmt.Println("op", op)
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

		// if readyBool {
		// 	mj, err := json.Marshal(lambdaInput{
		// 		Game:   game,
		// 		ApiId:  req.RequestContext.APIID,
		// 		Stage:  req.RequestContext.Stage,
		// 		Region: reg,
		// 	})
		// 	if err != nil {
		// 		fmt.Println("game item marshal err", err)
		// 	}

		// 	ii := lamb.InvokeInput{
		// 		FunctionName: aws.String("ct-playJS"),
		// 		// ClientContext:  new(string),
		// 		// InvocationType: "",
		// 		// LogType:        "",
		// 		Payload: mj,
		// 		// Qualifier:      new(string),
		// 	}

		// 	li, err := svc2.Invoke(ctx, &ii)

		// 	fmt.Printf("\n%s, %+v\n", "liii", *li)
		// 	// fmt.Println(*li.FunctionError, li.Payload)
		// 	q := *li
		// 	z := q.FunctionError
		// 	x := q.Payload
		// 	// fmt.Println(*z, x)

		// 	if z != nil {
		// 		fmt.Println(*z, x)
		// 	}

		// 	if err != nil {

		// 		var intServErr *types.InternalServerError
		// 		if errors.As(err, &intServErr) {
		// 			fmt.Printf("get item error, %v",
		// 				intServErr.ErrorMessage())
		// 		}

		// 		// To get any API error
		// 		var apiErr smithy.APIError
		// 		if errors.As(err, &apiErr) {
		// 			fmt.Printf("db error, Code: %v, Message: %v",
		// 				apiErr.ErrorCode(), apiErr.ErrorMessage())
		// 		}

		// 	}

		// }

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
