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

type listPlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
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

const maxPlayersPerGame string = "8"

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
			Action, Gameno, Tipe string
		}
		gameno    string
		exAttrNms = map[string]string{
			"#P": "players",
			"#I": id,
			"#R": "ready",
		}
		connKey = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "CONNECT"},
			"sk": &types.AttributeValueMemberS{Value: id},
		}
	)

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	if body.Gameno == "new" {
		gameno = fmt.Sprintf("%d", time.Now().UnixNano())
	} else if body.Gameno == "dc" {
		gameno = body.Gameno
	} else if _, err = strconv.ParseInt(body.Gameno, 10, 64); err != nil {
		return getReturnValue(http.StatusBadRequest), err
	} else if len(body.Gameno) != 19 {
		return getReturnValue(http.StatusBadRequest), errors.New("invalid game number")
	} else {
		gameno = body.Gameno
	}

	gameItemKey, err := attributevalue.MarshalMap(struct {
		Pk string `dynamodbav:"pk"`
		Sk string `dynamodbav:"sk"`
	}{
		Pk: "LISTGAME",
		Sk: gameno,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal gik Record, %v", err))
	}

	removePlayerInput := dynamodb.UpdateItemInput{
		Key:                      gameItemKey,
		TableName:                aws.String(tableName),
		ExpressionAttributeNames: exAttrNms,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":f": &types.AttributeValueMemberBOOL{Value: false},
		},
		UpdateExpression: aws.String("REMOVE #P.#I SET #R = :f"),
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
			Key:                 connKey,
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("size (#G) = :z"),
			ExpressionAttributeNames: map[string]string{
				"#G": "game",
				"#R": "returning",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: gameno},
				":z": &types.AttributeValueMemberN{Value: "0"},
				":f": &types.AttributeValueMemberBOOL{Value: false},
			},
			UpdateExpression: aws.String("SET #G = :g, #R = :f"),
		}

		_, err = ddbsvc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                      gameItemKey,
						TableName:                aws.String(tableName),
						ConditionExpression:      aws.String("attribute_exists(#P) AND size (#P) < :m"),
						ExpressionAttributeNames: exAttrNms,
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":f": &types.AttributeValueMemberBOOL{Value: false},
							":m": &types.AttributeValueMemberN{Value: maxPlayersPerGame},
							":p": marshalledPlayer,
						},
						UpdateExpression: aws.String("SET #P.#I = :p, #R = :f"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		callErr(err)

		_, err = ddbsvc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_not_exists(#P)"),
						ExpressionAttributeNames: map[string]string{
							"#P": "players",
							"#R": "ready",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":p": marshalledPlayersMap,
							":f": &types.AttributeValueMemberBOOL{Value: false},
						},
						UpdateExpression: aws.String("SET #P = :p, #R = :f"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		callErr(err)

	} else if body.Tipe == "leave" {

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       connKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#G": "game",
				"#L": "leader",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: ""},
				":f": &types.AttributeValueMemberBOOL{Value: false},
			},
			UpdateExpression: aws.String("SET #G = :g, #L = :f"),
		})
		callErr(err)

		ui2, err := ddbsvc.UpdateItem(ctx, &removePlayerInput)
		callErr(err)

		callFunction(ui2.Attributes, gameItemKey, tableName, ctx, ddbsvc)

	} else if body.Tipe == "ready" {

		ui2, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:                      gameItemKey,
			TableName:                aws.String(tableName),
			ExpressionAttributeNames: exAttrNms,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":t": &types.AttributeValueMemberBOOL{Value: true},
			},
			UpdateExpression: aws.String("SET #P.#I.#R = :t"),
			ReturnValues:     types.ReturnValueAllNew,
		})

		callErr(err)

		callFunction(ui2.Attributes, gameItemKey, tableName, ctx, ddbsvc)

	} else if body.Tipe == "unready" {

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:                      gameItemKey,
			TableName:                aws.String(tableName),
			ExpressionAttributeNames: exAttrNms,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":f": &types.AttributeValueMemberBOOL{Value: false},
			},
			UpdateExpression: aws.String("SET #P.#I.#R = :f, #R = :f"),
		})
		callErr(err)

	} else if body.Tipe == "disconnect" {
		if gameno != "dc" {

			ui2, err := ddbsvc.UpdateItem(ctx, &removePlayerInput)

			callErr(err)

			callFunction(ui2.Attributes, gameItemKey, tableName, ctx, ddbsvc)

		}

		_, err = ddbsvc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key:       connKey,
			TableName: aws.String(tableName),
		})
		callErr(err)

	} else {
		fmt.Println("other lobby")
	}

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
}

func callFunction(rv, gik map[string]types.AttributeValue, tn string, ctx context.Context, ddbsvc *dynamodb.Client) {
	var gm struct {
		Pk, Sk  string
		Ready   bool
		Players map[string]listPlayer
	}
	err := attributevalue.UnmarshalMap(rv, &gm)
	if err != nil {
		fmt.Println("unmarshal err", err)
	}

	if len(gm.Players) < 3 {
		return
	}

	readyCount := 0
	for k, v := range gm.Players {
		if v.Ready {
			readyCount++
			if readyCount == len(gm.Players) {
				time.Sleep(1000 * time.Millisecond)
				_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
					Key:       gik,
					TableName: aws.String(tn),
					ExpressionAttributeNames: map[string]string{
						"#R": "ready",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":t": &types.AttributeValueMemberBOOL{Value: true},
					},
					UpdateExpression: aws.String("SET #R = :t"),
				})
				callErr(err)

				_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
					Key: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: "CONNECT"},
						"sk": &types.AttributeValueMemberS{Value: k},
					},
					TableName: aws.String(tn),
					ExpressionAttributeNames: map[string]string{
						"#L": "leader",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":t": &types.AttributeValueMemberBOOL{Value: true},
					},
					UpdateExpression: aws.String("SET #L = :t"),
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
