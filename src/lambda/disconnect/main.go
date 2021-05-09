package main

import (
	"context"
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

type game struct {
	Pk       string            `dynamodbav:"pk"`
	Sk       string            `dynamodbav:"sk"`
	Starting bool              `dynamodbav:"starting"`
	Leader   string            `dynamodbav:"leader"`
	Loading  bool              `dynamodbav:"loading"`
	Players  map[string]Player `dynamodbav:"players"`
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

	op, err := svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       connItemKey,
		TableName: aws.String(tableName),

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

		ui2, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
				"#LE": "leader",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":e": &types.AttributeValueMemberS{Value: ""},
			},
			UpdateExpression: aws.String("REMOVE #PL.#ID SET #LE = :e"),
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

		callFunction(ui2.Attributes, gameItemKey, tableName, ctx, svc)

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

func callFunction(rv, gik map[string]types.AttributeValue, tn string, ctx context.Context, svc *dynamodb.Client) {
	var gm game
	err := attributevalue.UnmarshalMap(rv, &gm)
	if err != nil {
		fmt.Println("unmarshal err", err)
	}

	if len(gm.Players) < 3 {
		return
	}

	readyCount := 0
	for k, v := range gm.Players {

		fmt.Printf("%s, %v, %+v", "uicf", k, v)

		if v.Ready {
			readyCount++
			if readyCount == len(gm.Players) {
				time.Sleep(1000 * time.Millisecond)
				_, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
					Key:       gik,
					TableName: aws.String(tn),
					ExpressionAttributeNames: map[string]string{
						"#LE": "leader",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":l": &types.AttributeValueMemberS{Value: v.Name + "_" + v.ConnID},
					},
					UpdateExpression: aws.String("SET #LE = :l"),
				})

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
		} else {
			return
		}
	}
}
