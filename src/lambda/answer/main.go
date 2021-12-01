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
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/aws/smithy-go"
)

type key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

type player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
	Color  string `dynamodbav:"color,omitempty"`
	Score  int    `dynamodbav:"score"`
	Answer answer `dynamodbav:"answer"`
}

type answer struct {
	PlayerID, Answer string
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
	SendToFront  bool          `dynamodbav:"sendToFront"`
	HiScore      int           `dynamodbav:"hiScore"`
	GameTied     bool          `dynamodbav:"gameTied"`
}

type body struct {
	Gameno, Answer, PlayersCount string
}

type sfnArrInput struct {
	Id   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

type sfnInput struct {
	Gameno  string        `dynamodbav:"gameno"`
	Players []sfnArrInput `dynamodbav:"players"`
}

// func (pm livePlayerMap) assignColors() livePlayerMap {
// 	count := 0
// 	for k, v := range pm {
// 		v.Color = colors[count]
// 		pm[k] = v
// 		count++
// 	}

// 	return pm
// }

func (pm livePlayerMap) mapToSlice() (res []sfnArrInput) {
	for k, v := range pm {
		res = append(res, sfnArrInput{
			Id:   k,
			Name: v.Name,
		})
	}

	return
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

	sfnarn, ok := os.LookupEnv("SFNARN")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find sfn arn"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)
	sfnsvc := sfn.NewFromConfig(cfg)

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
			"#PL": "players",
			"#ID": id,
			"#AN": "answer",
			"#AC": "answersCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": marshalledAnswer,
			":o": &types.AttributeValueMemberN{Value: "1"},
		},
		UpdateExpression: aws.String("SET #PL.#ID.#AN = :a ADD #AC :o"),
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
				"#P":  "previousWord",
				"#C":  "currentWord",
				"#AC": "answersCount",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":c": &types.AttributeValueMemberS{Value: gm.CurrentWord},
				":b": &types.AttributeValueMemberS{Value: ""},
				":z": &types.AttributeValueMemberN{Value: "0"},
			},
			UpdateExpression: aws.String("SET #P = :c, #C = :b, #AC = :z"),
		})

		if err != nil {
			return callErr(err)
		}

		answers := getAnswersMap(gm.Players)

		updateScores(answers, gm)

		marshalledPlayersMap, err := attributevalue.Marshal(gm.Players)
		if err != nil {
			return callErr(err)
		}

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":p": marshalledPlayersMap,
			},
			UpdateExpression: aws.String("SET #PL = :p"),
		})

		if err != nil {
			return callErr(err)
		}

		taskOutput, err := json.Marshal(hiScore)
		if err != nil {
			return callErr(err)
		}

		stsi := sfn.SendTaskSuccessInput{
			Output:    aws.String(string(taskOutput)),
			TaskToken: aws.String(gm2.Token),
		}

		_, err = sfnsvc.SendTaskSuccess(ctx, &stsi)
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

func getAnswersMap(players livePlayerMap) (res map[string][]string) {
	for k, v := range players {
		fmt.Printf("%s, %v, %+v", "anssss", k, v)
		res[v.Answer.Answer] = append(res[v.Answer.Answer], v.Answer.PlayerID)
	}

	return
}

func updateScores(answers map[string][]string, game liveGame) {
	for k, v := range answers {
		fmt.Printf("%s, %v, %+v", "anssssmapppp", k, v)

		switch {
		case len(k) < 2: // c.updateEachScore(v, 0)

		case len(v) > 2: // c.updateEachScore(v, 1)
			for _, id := range v {
				game.Players[id] = adjScore(game.Players[id], 1, game.HiScore, game.GameTied)
			}

		case len(v) == 2: // c.updateEachScore(v, 3)
			for _, id := range v {
				game.Players[id] = adjScore(game.Players[id], 3, game.HiScore, game.GameTied)
			}

		default: // c.updateEachScore(v, 0)
		}
	}
}

func adjScore(old livePlayer, incr int, hiScore int, tied bool) livePlayer {
	score := old.Score + incr

	if score == hiScore {
		tied = true
	}
	if score > hiScore {
		hiScore = score
		tied = false
	}

	return livePlayer{
		Name:   old.Name,
		ConnID: old.ConnID,
		Color:  old.Color,
		Score:  score,
		Answer: old.Answer,
	}
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
