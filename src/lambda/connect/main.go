package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/smithy-go"
)

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go | build-lambda-zip.exe -o main.zip main / sam local invoke ConnectFunction -e ./event.json

// ConnItem holds values to be put in db
// type ConnItem struct {
// 	Pk     string `dynamodbav:"pk"`     //'CONN' + uuid
// 	Sk     string `dynamodbav:"sk"`     //name
// 	Game   string `dynamodbav:"game"`   //game no or blank
// 	GSI1PK string `dynamodbav:"GSI1PK"` //'CONN'
// 	GSI1SK string `dynamodbav:"GSI1SK"` //conn id
// }

// ConnItem2 holds values to be put in db
type ConnItem2 struct {
	Pk     string `json:"pk"`     //'CONN' + uuid
	Sk     string `json:"sk"`     //name
	Game   string `json:"game"`   //game no or blank
	GSI1PK string `json:"GSI1PK"` //'CONN'
	GSI1SK string `json:"GSI1SK"` //conn id
}

// StatItem holds values to be put in db
// type StatItem struct {
// 	Pk     string `dynamodbav:"pk"`     //'STAT' + uuid
// 	Sk     string `dynamodbav:"sk"`     //name
// 	GSI1PK string `dynamodbav:"GSI1PK"` //'STAT'
// 	GSI1SK string `dynamodbav:"GSI1SK"` //wins
// }

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

	// svc := dynamodb.NewFromConfig(cfg)
	svc2 := sfn.NewFromConfig(cfg)

	// svc.Handlers.Send.PushFront(func(r *request.Request) {
	// 	r.HTTPRequest.Header.Set("CustomHeader", fmt.Sprintf("%d", 10))
	// })
	// auth := req.RequestContext.Authorizer.(map[string]interface{})

	// auth["principalId"].(string)
	// auth["username"].(string)

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	id := auth["principalId"].(string)
	name := auth["username"].(string)

	// connItem, err := attributevalue.MarshalMap()
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to marshal Record, %v", err))
	// }

	// connid, err := attributevalue.Marshal("CONN#" + id)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to marshal Record 10, %v", err))
	// }

	// statid, err := attributevalue.Marshal("STAT#" + id)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to marshal Record 11, %v", err))
	// }

	// statItem, err := attributevalue.MarshalMap(StatItem{
	// 	Pk:     "STAT#" + id,
	// 	Sk:     name,
	// 	GSI1PK: "STAT",
	// 	GSI1SK: "0",
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
	// }

	// tableName, ok := os.LookupEnv("tableName")
	// if !ok {
	// 	panic(fmt.Sprintf("%v", "cant find table name"))
	// }

	// connItemInput := dynamodb.PutItemInput{
	// 	TableName: aws.String(tableName),
	// 	Item:      connItem,
	// 	// ExpressionAttributeValues: map[string]types.AttributeValue{
	// 	// 	":id": connid,
	// 	// },
	// 	ConditionExpression: aws.String("attribute_not_exists(pk)"),
	// }

	// err = panicProtectedPut(ctx, svc, &connItemInput)

	// if err != nil {
	// 	// fmt.Println("poi", err)
	// 	var condCheckErr *types.ConditionalCheckFailedException
	// 	if errors.As(err, &condCheckErr) {
	// 		fmt.Printf("connection already exists, not putting, %v\n", condCheckErr.ErrorMessage())

	// 	} else {

	// 		// To get any API error
	// 		var apiErr smithy.APIError
	// 		if errors.As(err, &apiErr) {
	// 			fmt.Printf("db error 1, Code: %v, Message: %v",
	// 				apiErr.ErrorCode(), apiErr.ErrorMessage())
	// 		}

	// 	}
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode:        http.StatusBadRequest,
	// 		Headers:           map[string]string{"Content-Type": "application/json"},
	// 		MultiValueHeaders: map[string][]string{},
	// 		Body:              "baddd",
	// 		IsBase64Encoded:   false,
	// 	}, err

	// }
	// statItemInput := dynamodb.PutItemInput{
	// 	TableName: aws.String(tableName),
	// 	Item:      statItem,
	// 	// ExpressionAttributeNames: map[string]string,
	// 	// ExpressionAttributeValues: map[string]types.AttributeValue{
	// 	// 	":id": statid,
	// 	// },
	// 	ConditionExpression: aws.String("attribute_not_exists(pk)"),
	// }

	// err = panicProtectedPut(ctx, svc, &statItemInput)

	// if err != nil {
	// 	// fmt.Println("poi", err)
	// 	var condCheckErr *types.ConditionalCheckFailedException
	// 	if errors.As(err, &condCheckErr) {
	// 		fmt.Printf("stat already exists, not putting, %v\n", condCheckErr.ErrorMessage())

	// 	} else {

	// 		// To get any API error
	// 		var apiErr smithy.APIError
	// 		if errors.As(err, &apiErr) {
	// 			fmt.Printf("db error 2, Code: %v, Message: %v",
	// 				apiErr.ErrorCode(), apiErr.ErrorMessage())
	// 		}

	// 		return events.APIGatewayProxyResponse{
	// 			StatusCode:        http.StatusOK,
	// 			Headers:           map[string]string{"Content-Type": "application/json"},
	// 			MultiValueHeaders: map[string][]string{},
	// 			Body:              "",
	// 			IsBase64Encoded:   false,
	// 		}, err

	// 	}

	// }

	smarn, ok := os.LookupEnv("smarn")
	if !ok {
		panic(fmt.Sprintf("%v", "cant find smarn"))
	}

	seii, err := json.Marshal(ConnItem2{
		Pk:     "CONN#" + id,
		Sk:     name,
		Game:   "",
		GSI1PK: "CONN",
		GSI1SK: req.RequestContext.ConnectionID,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	sei := sfn.StartExecutionInput{
		StateMachineArn: aws.String(smarn),
		Input:           aws.String(string(seii)),
		Name:            aws.String("bill1"),
		// TraceHeader:     new(string),
	}

	_, err = svc2.StartExecution(ctx, &sei)

	if err != nil {
		// fmt.Println("poi", err)
		var condCheckErr *types.ConditionalCheckFailedException
		if errors.As(err, &condCheckErr) {
			fmt.Printf("connection already exists, not putting, %v\n", condCheckErr.ErrorMessage())

		} else {

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("db error 1, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}

		}
		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusBadRequest,
			Headers:           map[string]string{"Content-Type": "application/json"},
			MultiValueHeaders: map[string][]string{},
			Body:              "baddd222",
			IsBase64Encoded:   false,
		}, err

	}

	// fmt.Println("smmmm", se)

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, nil
}

// func panicProtectedPut(ctx context.Context, svc *dynamodb.Client, pii *dynamodb.PutItemInput) error {
// 	fmt.Println("panicProtectedPut called", pii)
// 	defer func() {
// 		recover()
// 	}()
// 	_, err := svc.PutItem(ctx, pii)

// 	return err
// }

func main() {
	lambda.Start(handler)
}
