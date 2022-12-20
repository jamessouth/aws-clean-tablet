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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/smithy-go"
)

const connect string = "CONNECT"

func getReturnValue(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}
}

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
		ddbsvc        = dynamodb.NewFromConfig(cfg)
		sfnsvc        = sfn.NewFromConfig(cfg)
		auth          = req.RequestContext.Authorizer.(map[string]interface{})
		id, tableName = auth["principalId"].(string), auth["tableName"].(string)
		et            struct {
			Endtoken string
		}
	)

	ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: connect},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#T": "endtoken",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":t": &types.AttributeValueMemberS{Value: ""},
		},
		UpdateExpression: aws.String("SET #T = :t"), //TODO - use remove?
		ReturnValues:     types.ReturnValueAllOld,
	})
	if err != nil {
		return callErr(err)
	}

	err = attributevalue.UnmarshalMap(ui.Attributes, &et)
	if err != nil {
		return callErr(err)
	}

	stsi := sfn.SendTaskSuccessInput{
		Output:    aws.String("\"\""),
		TaskToken: aws.String(et.Endtoken),
	}

	_, err = sfnsvc.SendTaskSuccess(ctx, &stsi)
	if err != nil {
		return callErr(err)
	}

	return getReturnValue(http.StatusOK), nil
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

	return getReturnValue(http.StatusBadRequest), err
}
