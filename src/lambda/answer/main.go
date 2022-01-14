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

type key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

type answer struct {
	PlayerID string `json,dynamodbav:"playerid"`
	Answer   string `json,dynamodbav:"answer"`
}

type livePlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Color  string `dynamodbav:"color"`
	Score  int    `dynamodbav:"score"`
	Answer answer `dynamodbav:"answer"`
}

type livePlayerMap map[string]livePlayer

type liveGame struct {
	Pk           string        `dynamodbav:"pk"`
	Sk           string        `dynamodbav:"sk"`
	CurrentWord  string        `dynamodbav:"currentWord"`
	Players      livePlayerMap `dynamodbav:"players"`
	AnswersCount int           `dynamodbav:"answersCount"`
	HiScore      int           `dynamodbav:"hiScore"`
	GameTied     bool          `dynamodbav:"gameTied"`
}

type body struct {
	Gameno string `json:"gameno"`
	Answer string `json:"answer"`
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("answer", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return callErr(err)
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)

	id := req.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	gameItemKey, err := attributevalue.MarshalMap(key{
		Pk: "LIVEGME",
		Sk: body.Gameno,
	})
	if err != nil {
		return callErr(err)
	}

	marshalledAnswer, err := attributevalue.Marshal(answer{
		PlayerID: id,
		Answer:   body.Answer,
	})
	if err != nil {
		return callErr(err)
	}

	ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:       gameItemKey,
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "players",
			"#I": id,
			"#A": "answer",
			"#C": "answersCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": marshalledAnswer,
			":o": &types.AttributeValueMemberN{Value: "1"},
		},
		UpdateExpression: aws.String("SET #P.#I.#A = :a ADD #C :o"),
		ReturnValues:     types.ReturnValueAllNew,
	})

	if err != nil {
		return callErr(err)
	}

	var gm liveGame
	err = attributevalue.UnmarshalMap(ui.Attributes, &gm)
	if err != nil {
		return callErr(err)
	}

	if len(gm.Players) == gm.AnswersCount {

		_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#P": "previousWord",
				"#C": "currentWord",
				"#A": "answersCount",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":c": &types.AttributeValueMemberS{Value: gm.CurrentWord},
				":b": &types.AttributeValueMemberS{Value: ""},
				":z": &types.AttributeValueMemberN{Value: "0"},
			},
			UpdateExpression: aws.String("SET #P = :c, #C = :b, #A = :z"),
		})

		if err != nil {
			return callErr(err)
		}

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
