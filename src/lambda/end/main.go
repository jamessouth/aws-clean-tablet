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

const (
	connect string = "CONNECT"
)

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		bod = req.Body
		reg = strings.Split(req.RequestContext.DomainName, ".")[2]
	)

	if len(bod) > 75 { //TODO replace with observed value
		callErr(errors.New("improper json input - too long"))
	}

	fmt.Println("end", bod, len(bod))

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		callErr(err)
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
		UpdateExpression: aws.String("SET #T = :t"),
		ReturnValues:     types.ReturnValueAllOld,
	})
	if err != nil {
		callErr(err)
	}

	err = attributevalue.UnmarshalMap(ui.Attributes, &et)
	if err != nil {
		callErr(err)
	}

	stsi := sfn.SendTaskSuccessInput{
		Output:    aws.String("\"\""),
		TaskToken: aws.String(et.Endtoken),
	}

	_, err = sfnsvc.SendTaskSuccess(ctx, &stsi)
	if err != nil {
		callErr(err)
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

func callErr(err error) {
	if err != nil {

		// To get any API error
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("db error, Code: %v, Message: %v",
				apiErr.ErrorCode(), apiErr.ErrorMessage())
		}

	}
}
