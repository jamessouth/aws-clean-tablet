package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

const connect string = "CONNECT"

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		bod    = req.Body
		region = strings.Split(req.RequestContext.DomainName, ".")[2]
	)

	if len(bod) > 75 { //TODO replace with observed value
		return callErr(errors.New("improper json input - too long"))
	}

	fmt.Println("end", bod, len(bod))

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return callErr(err)
	}

	var (
		ddbsvc              = dynamodb.NewFromConfig(cfg)
		auth                = req.RequestContext.Authorizer.(map[string]interface{})
		id, name, tableName = auth["principalId"].(string), auth["username"].(string), auth["tableName"].(string)
	)

	_, err = ddbsvc.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":      &types.AttributeValueMemberS{Value: connect},
			"sk":      &types.AttributeValueMemberS{Value: id},
			"game":    &types.AttributeValueMemberS{Value: ""},
			"name":    &types.AttributeValueMemberS{Value: name},
			"playing": &types.AttributeValueMemberBOOL{Value: false},
			"color":   &types.AttributeValueMemberS{Value: "transparent"},
			"connID":  &types.AttributeValueMemberS{Value: req.RequestContext.ConnectionID},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		return callErr(err)
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

func callErr(err error) (events.APIGatewayProxyResponse, error) {
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

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusBadRequest,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, err
}
