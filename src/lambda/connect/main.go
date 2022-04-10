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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

// type statItem struct {
// 	Pk     string `dynamodbav:"pk"`     //'STAT#' + uuid
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
		return callErr(err)
	}

	var (
		tableName = os.Getenv("tableName")
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		auth      = req.RequestContext.Authorizer.(map[string]interface{})
		id, name  = auth["principalId"].(string), auth["username"].(string)
	)

	_, err = ddbsvc.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":      &types.AttributeValueMemberS{Value: "CONNECT#" + id},
			"sk":      &types.AttributeValueMemberS{Value: name},
			"game":    &types.AttributeValueMemberS{Value: ""},
			"playing": &types.AttributeValueMemberBOOL{Value: false},
			"leader":  &types.AttributeValueMemberBOOL{Value: false},
			"color":   &types.AttributeValueMemberS{Value: "transparent"},
			"index":   &types.AttributeValueMemberS{Value: ""},
			"GSI1PK":  &types.AttributeValueMemberS{Value: "CONNECT"},
			"GSI1SK":  &types.AttributeValueMemberS{Value: req.RequestContext.ConnectionID},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		return callErr(err)
	}

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "STAT#" + id},
			"sk": &types.AttributeValueMemberS{Value: name},
		},

		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "GSI1PK",
			"#S": "GSI1SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p": &types.AttributeValueMemberS{Value: "STAT"},
			":s": &types.AttributeValueMemberN{Value: "0"},
		},
		UpdateExpression: aws.String("SET #P = :p ADD #S :s"),
		// ReturnValues:     types.ReturnValueAllNew,
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
