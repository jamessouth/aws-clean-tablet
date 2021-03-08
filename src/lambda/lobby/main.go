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

// Player for player info
type Player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
}

// GameItemKey holds values to be put in db
type GameItemKey struct {
	Pk string `dynamodbav:"pk"` //GAME
	Sk string `dynamodbav:"sk"` //game no
}

// GameItemAttrs holds values to be put in db
type GameItemAttrs struct {
	Players []Player `dynamodbav:":p"`
}

type body struct {
	Game string
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("lobbbbbby", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.NewFromConfig(cfg)

	var gameno string
	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	fmt.Printf("body: %v\n", body.Game)

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	if body.Game == "new" {
		gameno = fmt.Sprintf("%d", time.Now().UnixNano())
	} else {
		gameno = body.Game
	}

	g, err := attributevalue.MarshalMap(GameItemKey{
		Pk: "GAME",
		Sk: gameno,
		// Players: []Player{
		// 	{
		// 		Name:   auth["username"].(string),
		// 		ConnID: req.RequestContext.ConnectionID,
		// 	},
		// },
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record, %v", err))
	}
	g2, err := attributevalue.MarshalMap(GameItemAttrs{
		// Pk: "GAME",
		// Sk: gameno,
		Players: []Player{
			{
				Name:   auth["username"].(string),
				ConnID: req.RequestContext.ConnectionID,
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal Record 2, %v", err))
	}

	for k, v := range g2 {

		fmt.Println("g2", k, v)
	}

	op, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:       g,
		TableName: aws.String(tableName),

		ExpressionAttributeNames: map[string]string{
			"#PL": "Players",
		},
		ExpressionAttributeValues: g2,
		ReturnConsumedCapacity:    types.ReturnConsumedCapacityTotal,

		UpdateExpression: aws.String("SET #PL = :p"),
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
		Body:              fmt.Sprintf("cap used: %v", op.ConsumedCapacity.CapacityUnits),
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
