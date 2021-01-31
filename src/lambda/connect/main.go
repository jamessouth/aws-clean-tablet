package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go / build-lambda-zip.exe -o main.zip main / sam local invoke ConnectFunction -e ../event.json

// Item holds values to be put in db
type Item struct {
	Pk     string `json:"pk"`
	Sk     string `json:"sk"`
	GSI1Pk string `json:"gsi1pk"`
	GSI1Sk string `json:"gsi1sk"`
	ConnID string `json:"connid"`
}

func handler(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("session init error")
	}
	logger := aws.NewDefaultLogger()

	sess.Handlers.Send.PushFront(func(r *request.Request) {
		logger.Log(fmt.Sprintf("Request: %s /%v, Payload: %s",
			r.ClientInfo.ServiceName, r.Operation, r.Params))
	})

	svc := dynamodb.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint("http://192.168.4.27:8000"))

	// svc.Handlers.Send.PushFront(func(r *request.Request) {
	// 	r.HTTPRequest.Header.Set("CustomHeader", fmt.Sprintf("%d", 10))
	// })
	auth := req.RequestContext.Authorizer.(map[string]interface{})

	i, err := dynamodbattribute.MarshalMap(Item{
		Pk:     auth["principalId"].(string),
		Sk:     "GAME#",
		GSI1Pk: auth["username"].(string),
		GSI1Sk: "STATS#",
		ConnID: req.RequestContext.ConnectionID,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "cant find table name"))
	}

	op, err := svc.PutItem(&dynamodb.PutItemInput{
		TableName:              aws.String(tableName),
		Item:                   i,
		ReturnConsumedCapacity: aws.String("TOTAL"),
	})
	fmt.Println("op", op)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		// return
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              fmt.Sprintf("cap used: %v", op.ConsumedCapacity),
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
