package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

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
// type GameItemAttrs struct {
// 	Players []string `dynamodbav:":p,stringset"` //name + connid
// }

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	// fmt.Printf("%s: %+v\n", "request", req)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.NewFromConfig(cfg)

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	id := auth["principalId"].(string)
	name := auth["username"].(string)

	connItemKey, err := attributevalue.MarshalMap(Key{
		Pk: "CONN#" + id,
		Sk: name,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "cant find table name"))
	}

	// auth := req.RequestContext.Authorizer.(map[string]interface{})

	op, err := svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       connItemKey,
		TableName: aws.String(tableName),

		// ExpressionAttributeNames:  map[string]string{},
		// ExpressionAttributeValues: map[string]types.AttributeValue{},
		// ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,

		ReturnValues: types.ReturnValueAllOld,
	})

	if err != nil {

		var intServErr *types.InternalServerError
		if errors.As(err, &intServErr) {
			fmt.Printf("delete item error, %v",
				intServErr.ErrorMessage())
		}

		// To get any API error
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("db error, Code: %v, Message: %v",
				apiErr.ErrorCode(), apiErr.ErrorMessage())
		}

	}

	var game string
	err = attributevalue.Unmarshal(op.Attributes["game"], &game)
	if err != nil {
		fmt.Println("del item unmarshal err", err)
	}

	if len(game) > 0 {

		gameItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: game,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 7, %v", err))
		}

		// gameAttrs, err := attributevalue.MarshalMap(GameItemAttrs{
		// 	Players: []string{
		// 		auth["username"].(string) + "#" + req.RequestContext.ConnectionID,
		// 	},
		// })
		// if err != nil {
		// 	panic(fmt.Sprintf("failed to marshal Record 8, %v", err))
		// }

		op2, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("(attribute_exists(#PL) AND size (#PL) < :maxsize) OR attribute_not_exists(#PL)"),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
			},
			// ExpressionAttributeValues: gameAttrs,

			UpdateExpression: aws.String("REMOVE #PL.#ID"),
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
		err = attributevalue.Unmarshal(op2.Attributes["players"], &game)
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
