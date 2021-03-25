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
	Ready  bool   `dynamodbav:"ready"`
}

// GameItemAttrs holds values to be put in db
type GameItemAttrs struct {
	Players map[string]Player `dynamodbav:":p"`
}

// MaxSize int      `dynamodbav:":maxsize,omitempty"`

type body struct {
	Game, Type string
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

		marshaledID, err := attributevalue.Marshal(id)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal marshaledID, %v", err))
		}

		marshaledMaxSize, err := attributevalue.Marshal(maxPlayersPerGame)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal marshaledMaxSize, %v", err))
		}

		// Players: []string{auth["username"].(string) + "#" + },
		// MaxSize: maxPlayersPerGame,
		player, err := attributevalue.MarshalMap(Player{
			Name:   name,
			ConnID: req.RequestContext.ConnectionID,
			Ready:  false,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal player, %v", err))
		}

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

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{

			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{

						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":maxsize": marshaledMaxSize,
						},

						UpdateExpression: aws.String("ADD #PL :p"),
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

			var intServErr *types.InternalServerError
			if errors.As(err, &intServErr) {
				fmt.Printf("put item error, %v",
					intServErr.ErrorMessage())
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("db error, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}
	} else {
		connItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "CONN#" + id,
			Sk: name,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 3, %v", err))
		}

		gameAttrs, err := attributevalue.MarshalMap(GameItemAttrs{
			Players: []string{auth["username"].(string) + "#" + req.RequestContext.ConnectionID},
			// MaxSize: maxPlayersPerGame,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
		}

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

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{

			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{

						Key:       gameItemKey,
						TableName: aws.String(tableName),
						// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
						},
						ExpressionAttributeValues: gameAttrs,

						UpdateExpression: aws.String("DELETE #PL :p"),
					},
				},
				{
					Update: &types.Update{

						Key:       connItemKey,
						TableName: aws.String(tableName),
						// ConditionExpression: aws.String("size (#IG) = :zero"),
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
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},

		Body:            "",
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
