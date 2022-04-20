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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

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
	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	var (
		tableName = aws.String(os.Getenv("tableName"))
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		// auth      = req.RequestContext.Authorizer.(map[string]interface{})
		// id, name  = auth["principalId"].(string), auth["username"].(string)
		body struct {
			Action, Info string
		}
	)

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}
	fmt.Printf("%s%+v\n", "bod ", body)

	leadersResults, err := ddbsvc.Query(ctx, &dynamodb.QueryInput{
		TableName:              tableName,
		KeyConditionExpression: aws.String("pk = :s"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: "STAT"},
		},
	})
	if err != nil {
		return callErr(err)
	}

	var leaders []struct {
		Pk, Sk, Name             string
		Wins, TotalPoints, Games int
	}
	err = attributevalue.UnmarshalListOfMaps(leadersResults.Items, &leaders)
	callErr(err)
	fmt.Printf("%s%+v\n", "res ", leaders)

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
