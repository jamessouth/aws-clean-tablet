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

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go | build-lambda-zip.exe -o main.zip main / sam local invoke ConnectFunction -e ./event.json

// ConnItem holds values to be put in db
type ConnItem struct {
	Pk   string `dynamodbav:"pk"`   //'CONN'
	Sk   string `dynamodbav:"sk"`   //conn id
	Game string `dynamodbav:"game"` //game no or blank
}

// StatItem holds values to be put in db
type StatItem struct {
	Pk     string `dynamodbav:"pk"`     //uuid
	Sk     string `dynamodbav:"sk"`     //name
	GSI1PK string `dynamodbav:"GSI1PK"` //'STAT'
	GSI1SK int    `dynamodbav:"GSI1SK"` //wins
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}
	// logger := aws.NewDefaultLogger()

	// sess.Handlers.Send.PushFront(func(r *request.Request) {
	// 	logger.Log(fmt.Sprintf("Request: %s /%v, Payload: %s",
	// 		r.ClientInfo.ServiceName, r.Operation, r.Params))
	// })

	// .WithLogLevel(aws.LogDebugWithHTTPBody)

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.NewFromConfig(cfg)

	// svc.Handlers.Send.PushFront(func(r *request.Request) {
	// 	r.HTTPRequest.Header.Set("CustomHeader", fmt.Sprintf("%d", 10))
	// })
	// auth := req.RequestContext.Authorizer.(map[string]interface{})

	// auth["principalId"].(string)
	// auth["username"].(string)

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	connItem, err := attributevalue.MarshalMap(ConnItem{
		Pk:   "CONN",
		Sk:   req.RequestContext.ConnectionID,
		Game: "",
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	id := auth["principalId"].(string)

	mid, err := attributevalue.Marshal(id)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 6, %v", err))
	}

	statItem, err := attributevalue.MarshalMap(StatItem{
		Pk:     id,
		Sk:     auth["username"].(string),
		GSI1PK: "STAT",
		GSI1SK: 0,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "cant find table name"))
	}

	op, err := svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:              aws.String(tableName),
		Item:                   connItem,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
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
	pii := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      statItem,
		// ExpressionAttributeNames: map[string]string,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": mid,
		},
		ConditionExpression:    aws.String(":id <> pk"),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	err = panicProtectedPut(ctx, svc, &pii)

	if err != nil {
		// fmt.Println("poi", err)
		var condCheckErr *types.ConditionalCheckFailedException
		if errors.As(err, &condCheckErr) {
			fmt.Printf("item with this pk exists, not putting, %v\n", condCheckErr.ErrorMessage())

		} else {

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("db error, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

			return events.APIGatewayProxyResponse{
				StatusCode:        http.StatusBadRequest,
				Headers:           map[string]string{"Content-Type": "application/json"},
				MultiValueHeaders: map[string][]string{},
				Body:              fmt.Sprintf("cap used: %v, %v", &op.ConsumedCapacity.CapacityUnits, &op.ConsumedCapacity.CapacityUnits),
				IsBase64Encoded:   false,
			}, err

		}

	}

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              fmt.Sprintf("cap used: %v, %v", &op.ConsumedCapacity.CapacityUnits, &op.ConsumedCapacity.CapacityUnits),
		IsBase64Encoded:   false,
	}, nil
}

func panicProtectedPut(ctx context.Context, svc *dynamodb.Client, pii *dynamodb.PutItemInput) error {
	fmt.Println("panicProtectedPut called")
	defer func() {
		recover()
	}()
	_, err := svc.PutItem(ctx, pii)

	return err
}

func main() {
	lambda.Start(handler)
}
