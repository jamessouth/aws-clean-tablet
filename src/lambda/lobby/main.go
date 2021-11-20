package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

// export CGO_ENABLED=0 | go build -o main main.go | zip main.zip main | aws lambda update-function-code --function-name ct-lobby --zip-file fileb://main.zip

const maxPlayersPerGame string = "8"

type key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

type listPlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
}

type listGame struct {
	Pk      string                `dynamodbav:"pk"` //'GAME'
	Sk      string                `dynamodbav:"sk"` //no
	Ready   bool                  `dynamodbav:"ready"`
	Players map[string]listPlayer `dynamodbav:"players"`
}

type body struct {
	Action, Gameno, Tipe string
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
	// svc2 := lamb.NewFromConfig(cfg)

	auth := req.RequestContext.Authorizer.(map[string]interface{})

	id := auth["principalId"].(string)
	name := auth["username"].(string)

	var body body
	var gameno string

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	// var gameItemKey map[string]types.AttributeValue

	if body.Gameno == "new" {
		gameno = fmt.Sprintf("%d", time.Now().UnixNano())
	} else if body.Gameno == "dc" {
		gameno = body.Gameno
	} else if _, err = strconv.ParseInt(body.Gameno, 10, 64); err != nil {
		fmt.Println("ParseInt error: ", err)

		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusBadRequest,
			Headers:           map[string]string{"Content-Type": "application/json"},
			MultiValueHeaders: map[string][]string{},
			Body:              "",
			IsBase64Encoded:   false,
		}, err
	} else if len(body.Gameno) != 19 {
		err = errors.New("invalid game number")
		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusBadRequest,
			Headers:           map[string]string{"Content-Type": "application/json"},
			MultiValueHeaders: map[string][]string{},
			Body:              "",
			IsBase64Encoded:   false,
		}, err
	} else {
		gameno = body.Gameno
	}

	gameItemKey, err := attributevalue.MarshalMap(key{
		Pk: "GAME",
		Sk: gameno,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal gik Record, %v", err))
	}

	removePlayerInput := dynamodb.UpdateItemInput{
		Key:       gameItemKey,
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#PL": "players",
			"#ID": id,
			"#RE": "ready",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":f": &types.AttributeValueMemberBOOL{Value: false},
		},
		UpdateExpression: aws.String("REMOVE #PL.#ID SET #RE = :f"),
		ReturnValues:     types.ReturnValueAllNew,
	}

	if body.Tipe == "join" {

		player := listPlayer{
			Name:   name,
			ConnID: req.RequestContext.ConnectionID,
			Ready:  false,
		}

		marshalledPlayersMap, err := attributevalue.Marshal(map[string]listPlayer{
			id: player,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal map Record 22, %v", err))
		}

		marshalledPlayer, err := attributevalue.Marshal(player)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal indiv Record 22, %v", err))
		}

		updateConnInput := types.Update{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "CONN#" + id},
				"sk": &types.AttributeValueMemberS{Value: name},
			},
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("size (#IG) = :z"),
			ExpressionAttributeNames: map[string]string{
				"#IG": "game",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: gameno},
				":z": &types.AttributeValueMemberN{Value: "0"},
			},
			UpdateExpression: aws.String("SET #IG = :g"),
		}

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_exists(#PL) AND size (#PL) < :m"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
							"#ID": id,
							"#RE": "ready",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":f": &types.AttributeValueMemberBOOL{Value: false},
							":m": &types.AttributeValueMemberN{Value: maxPlayersPerGame},
							":p": marshalledPlayer,
						},
						UpdateExpression: aws.String("SET #PL.#ID = :p, #RE = :f"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		callErr(err)

		_, err = svc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_not_exists(#PL)"),
						ExpressionAttributeNames: map[string]string{
							"#PL": "players",
							"#RE": "ready",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":p": marshalledPlayersMap,
							":f": &types.AttributeValueMemberBOOL{Value: false},
						},
						UpdateExpression: aws.String("SET #PL = :p, #RE = :f"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		callErr(err)

	} else if body.Tipe == "leave" {

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "CONN#" + id},
				"sk": &types.AttributeValueMemberS{Value: name},
			},
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#IG": "game",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: ""},
			},
			UpdateExpression: aws.String("SET #IG = :g"),
		})
		callErr(err)

		ui2, err := svc.UpdateItem(ctx, &removePlayerInput)
		callErr(err)

		callFunction(ui2.Attributes, gameItemKey, tableName, ctx, svc)

	} else if body.Tipe == "ready" {

		ui2, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
				"#RD": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":t": &types.AttributeValueMemberBOOL{Value: true},
			},
			UpdateExpression: aws.String("SET #PL.#ID.#RD = :t"),
			ReturnValues:     types.ReturnValueAllNew,
		})

		callErr(err)

		callFunction(ui2.Attributes, gameItemKey, tableName, ctx, svc)

	} else if body.Tipe == "unready" {

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
				"#RD": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":f": &types.AttributeValueMemberBOOL{Value: false},
			},
			UpdateExpression: aws.String("SET #PL.#ID.#RD = :f, #RD = :f"),
		})
		callErr(err)

	} else if body.Tipe == "disconnect" {
		if gameno != "dc" {

			ui2, err := svc.UpdateItem(ctx, &removePlayerInput)

			callErr(err)

			callFunction(ui2.Attributes, gameItemKey, tableName, ctx, svc)

		}

		_, err = svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "CONN#" + id},
				"sk": &types.AttributeValueMemberS{Value: name},
			},
			TableName: aws.String(tableName),
		})
		callErr(err)

	} else {
		fmt.Println("other lobby")
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

func callFunction(rv, gik map[string]types.AttributeValue, tn string, ctx context.Context, svc *dynamodb.Client) {
	var gm listGame
	err := attributevalue.UnmarshalMap(rv, &gm)
	if err != nil {
		fmt.Println("unmarshal err", err)
	}

	if len(gm.Players) < 3 {
		return
	}

	readyCount := 0
	for k, v := range gm.Players {

		fmt.Printf("%s, %v, %+v", "uicf", k, v)

		if v.Ready {
			readyCount++
			if readyCount == len(gm.Players) {
				time.Sleep(1000 * time.Millisecond)
				_, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
					Key:       gik,
					TableName: aws.String(tn),
					ExpressionAttributeNames: map[string]string{
						"#RE": "ready",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":t": &types.AttributeValueMemberBOOL{Value: true},
					},
					UpdateExpression: aws.String("SET #RE = :t"),
				})
				callErr(err)
			}
		} else {
			return
		}
	}
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
