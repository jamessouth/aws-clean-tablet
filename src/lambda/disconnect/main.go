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

// GameItemAttrs holds values to be put in db
type GameItemAttrs struct {
	Players []string `dynamodbav:":p,stringset"` //name + connid
}

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

	k, err := attributevalue.MarshalMap(Key{
		Pk: "CONN",
		Sk: req.RequestContext.ConnectionID,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "cant find table name"))
	}

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	op, err := svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       k,
		TableName: aws.String(tableName),

		// ExpressionAttributeNames:  map[string]string{},
		// ExpressionAttributeValues: map[string]types.AttributeValue{},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,

		ReturnValues: types.ReturnValueAllOld,
	})

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

	var connItem string
	err = attributevalue.Unmarshal(op.Attributes["game"], &connItem)
	if err != nil {
		fmt.Println("del item unmarshal err", err)
	}

	// fmt.Printf("conn item: %v\n", connItem)
	// fmt.Printf("conn item2: %v\n", &connItem)
	// fmt.Println("join unmr2", connItem, &connItem, *connItem)

	if len(connItem) > 0 {
		// fmt.Printf("conn item: %v\n", connItem)
		gameKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: connItem,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 7, %v", err))
		}

		gameAttrs, err := attributevalue.MarshalMap(GameItemAttrs{
			Players: []string{
				auth["username"].(string) + "#" + req.RequestContext.ConnectionID,
			},
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 8, %v", err))
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			Key:       gameKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("contains(Color, :v_sub)"),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
			},
			ExpressionAttributeValues: gameAttrs,

			UpdateExpression: aws.String("DELETE #PL :p"),
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
		Body:              fmt.Sprintf("cap used: %v", &op.ConsumedCapacity.CapacityUnits),
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
