package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

// Key holds values to be put in db
type Key struct {
	Pk string `dynamodbav:"pk"` //GAME
	Sk string `dynamodbav:"sk"` //game no
}

// GameItemAttrs holds values to be put in db
type GameItemAttrs struct {
	Players []string `dynamodbav:":p,stringset"` //name + connid
}

type body struct {
	Game string
}

// ConnItemAttrs holds vals for db
type ConnItemAttrs struct {
	Game string `dynamodbav:":g"`
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	// fmt.Println("lobbbbbby", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.NewFromConfig(cfg)

	connItemKey, err := attributevalue.MarshalMap(Key{
		Pk: "CONN",
		Sk: req.RequestContext.ConnectionID,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 3, %v", err))
	}

	op3, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       connItemKey,
		TableName: aws.String(tableName),

		// ConsistentRead:           new(bool),
		// ExpressionAttributeNames: map[string]string{"#PL": "players"},
		// ProjectionExpression:     new(string),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})
	// fmt.Println("op", op)
	if err != nil {

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
	var connItem ConnItemAttrs
	err = attributevalue.UnmarshalMap(op3.Item, &connItem)
	if err != nil {
		fmt.Println("get item unmarshal err")
	}

	if len(connItem.Game) > 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "application/json"},

			Body:            fmt.Sprintf("cap used: %v", &op3.ConsumedCapacity.CapacityUnits),
			IsBase64Encoded: false,
		}, nil
	}

	var gameno string
	var body body
	// var ga map[string]types.AttributeValue
	// var ue string

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	// fmt.Printf("body: %v\n", body.Game)

	if body.Game == "new" {
		gameno = fmt.Sprintf("%d", time.Now().UnixNano())
	} else {
		gameno = body.Game
	}

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	ga, err := attributevalue.MarshalMap(GameItemAttrs{
		Players: []string{
			auth["username"].(string) + "#" + req.RequestContext.ConnectionID,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
	}

	ca, err := attributevalue.MarshalMap(ConnItemAttrs{
		InGame: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 4, %v", err))
	}

	gk, err := attributevalue.MarshalMap(Key{
		Pk: "GAME",
		Sk: gameno,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	ck, err := attributevalue.MarshalMap(Key{
		Pk: "CONN",
		Sk: req.RequestContext.ConnectionID,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 5, %v", err))
	}

	op, err := svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{

		TransactItems: []types.TransactWriteItem{
			{
				Update: &types.Update{

					Key:       gk,
					TableName: aws.String(tableName),
					// ConditionExpression: aws.String("contains(Color, :v_sub)"),
					ExpressionAttributeNames: map[string]string{
						"#PL": "players",
					},
					ExpressionAttributeValues: ga,

					UpdateExpression: aws.String("ADD #PL :p"),
				},
			},
			{
				Update: &types.Update{

					Key:       ck,
					TableName: aws.String(tableName),
					// ConditionExpression: aws.String("contains(Color, :v_sub)"),
					ExpressionAttributeNames: map[string]string{
						"#IG": "ingame",
					},
					ExpressionAttributeValues: ca,

					UpdateExpression: aws.String("SET #IG = :ig"),
				},
			},
		},
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
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},

		Body:            fmt.Sprintf("cap used: %v", &op.ConsumedCapacity),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
