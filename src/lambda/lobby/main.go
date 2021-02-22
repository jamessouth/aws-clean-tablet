package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// GameItem holds values to be put in db
type GameItem struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

func handler(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("%s: %+v\n", "lobbbbbby", req)

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("session init error")
	}

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.New(sess, aws.NewConfig())

	gameno := fmt.Sprintf("%d", time.Now().UnixNano())

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	g, err := dynamodbattribute.MarshalMap(GameItem{
		Pk:   "GAME#" + gameno,
		Sk:   req.RequestContext.ConnectionID,
		Name: auth["username"].(string),
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
		Item:                   g,
		ReturnConsumedCapacity: aws.String("TOTAL"),
	})
	// fmt.Println("op", op)
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
