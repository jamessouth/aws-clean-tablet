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
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
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
	var connItem string
	err = attributevalue.Unmarshal(op3.Item["game"], &connItem)
	if err != nil {
		fmt.Println("get item unmarshal err", err)
	}

	fmt.Println("join unmr", connItem)

	if len(connItem) > 0 {
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

	gameAttrs, err := attributevalue.MarshalMap(GameItemAttrs{
		Players: []string{
			auth["username"].(string) + "#" + req.RequestContext.ConnectionID,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
	}

	connAttrs, err := attributevalue.MarshalMap(ConnItemAttrs{
		Game: gameno,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 4, %v", err))
	}

	gameKey, err := attributevalue.MarshalMap(Key{
		Pk: "GAME",
		Sk: gameno,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}

	connKey, err := attributevalue.MarshalMap(Key{
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

					Key:       gameKey,
					TableName: aws.String(tableName),
					// ConditionExpression: aws.String("contains(Color, :v_sub)"),
					ExpressionAttributeNames: map[string]string{
						"#PL": "players",
					},
					ExpressionAttributeValues: gameAttrs,

					UpdateExpression: aws.String("ADD #PL :p"),
				},
			},
			{
				Update: &types.Update{

					Key:       connKey,
					TableName: aws.String(tableName),
					// ConditionExpression: aws.String("contains(Color, :v_sub)"),
					ExpressionAttributeNames: map[string]string{
						"#IG": "game",
					},
					ExpressionAttributeValues: connAttrs,

					UpdateExpression: aws.String("SET #IG = :g"),
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
