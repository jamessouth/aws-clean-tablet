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

// ConnItem holds values to be put in db
type ConnItem struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
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

	k, err := attributevalue.MarshalMap(ConnItem{
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

	op, err := svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:              aws.String(tableName),
		Key:                    k,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
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
