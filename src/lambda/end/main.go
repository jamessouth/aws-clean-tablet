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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

type livePlayerList []struct {
	PlayerID string `dynamodbav:"playerid"`
	Name     string `dynamodbav:"name"`
	ConnID   string `dynamodbav:"connid"`
	Color    string `dynamodbav:"color"`
	Score    int    `dynamodbav:"score"`
	Index    string `dynamodbav:"index"`
}

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
		tableName = os.Getenv("tableName")
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		auth      = req.RequestContext.Authorizer.(map[string]interface{})
		id, name  = auth["principalId"].(string), auth["username"].(string)
		body      struct {
			Action, Gameno string
			Leader         bool
		}
	)

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "CONNECT"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#G": "game",
			"#L": "leader",
			"#P": "playing",
			"#C": "color",
			"#I": "index",
			"#R": "returning",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":g": &types.AttributeValueMemberS{Value: ""},
			":f": &types.AttributeValueMemberBOOL{Value: false},
			":c": &types.AttributeValueMemberS{Value: "transparent"},
			":t": &types.AttributeValueMemberBOOL{Value: true},
		},
		UpdateExpression: aws.String("SET #G = :g, #L = :f, #P = :f, #C = :c, #I = :g, #R = :t"),
	})
	callErr(err)

	if body.Leader {

		di, err := ddbsvc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "LIVEGME"},
				"sk": &types.AttributeValueMemberS{Value: body.Gameno},
			},
			TableName:    aws.String(tableName),
			ReturnValues: types.ReturnValueAllOld,
		})
		callErr(err)

		// var game struct {
		// 	Players livePlayerList
		// 	Ids map[string]string
		// }
		// err = attributevalue.UnmarshalMap(di.Attributes, &game)
		// if err != nil {
		// 	return callErr(err)
		// }

		// ui2, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		// 	Key:                      gameItemKey,
		// 	TableName:                aws.String(tableName),
		// 	ExpressionAttributeNames: exAttrNms,
		// 	ExpressionAttributeValues: map[string]types.AttributeValue{
		// 		":t": &types.AttributeValueMemberBOOL{Value: true},
		// 	},
		// 	UpdateExpression: aws.String("SET #P.#I.#R = :t"),
		// 	ReturnValues:     types.ReturnValueAllNew,
		// })

		// callErr(err)

	}

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
}

func callErr(err error) {
	if err != nil {
		var transCxldErr *types.TransactionCanceledException
		if errors.As(err, &transCxldErr) {
			fmt.Printf("put item error777, %v\n",
				transCxldErr.CancellationReasons)
		}

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

	}
}
