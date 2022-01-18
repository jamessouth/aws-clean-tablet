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

type livePlayer struct {
	PlayerID string `json:"playerid"`
	Name     string `json:"name"`
	ConnID   string `json:"connid"`
	Color    string `json:"color"`
	Score    int    `json:"score"`
	Answer   string `json:"answer"`
}

type livePlayerList []livePlayer

type liveGame struct {
	Sk      string         `json:"sk"`
	Players livePlayerList `json:"players"`
}

type body struct {
	Game liveGame `json:"game"`
}

type master struct {
	Players livePlayerList
	Answers map[string][]string
	Scores  map[string]int
	Hi      int
	Tied    bool
}

const winThreshold int = 24

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("score", req.Body)

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

	// sfnarn, ok := os.LookupEnv("SFNARN")
	// if !ok {
	// 	panic(fmt.Sprintf("%v", "can't find sfn arn"))
	// }

	ddbsvc := dynamodb.NewFromConfig(cfg)
	// sfnsvc := sfn.NewFromConfig(cfg)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	winner := false

	scoreData := master{
		Players: body.Game.Players,
		Answers: map[string][]string{},
		Scores:  map[string]int{},
		Hi:      0,
		Tied:    false,
	}

	updatedScoreData := scoreData.getAnswersMap().getScoresMap().updateScores().checkHiScore()

	if !updatedScoreData.Tied && updatedScoreData.Hi > winThreshold {
		winner = true
	}

	marshalledPlayersMap, err := attributevalue.Marshal(updatedScoreData.Players)
	if err != nil {
		return callErr(err)
	}

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGME"},
			"sk": &types.AttributeValueMemberS{Value: body.Game.Sk},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "players",
			"#W": "winner",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p": marshalledPlayersMap,
			":w": &types.AttributeValueMemberBOOL{Value: winner},
		},
		UpdateExpression: aws.String("SET #P = :p, #W = :w"),
	})

	if err != nil {
		return callErr(err)
	}

	// if !winner {

	// 	sfnInput := "{\"gameno\":\"" + body.Game.No + "\"}"

	// 	ssei := sfn.StartSyncExecutionInput{
	// 		StateMachineArn: aws.String(sfnarn),
	// 		Input:           aws.String(sfnInput),
	// 	}

	// 	sse, err := sfnsvc.StartSyncExecution(ctx, &ssei)
	// 	if err != nil {
	// 		return callErr(err)
	// 	}

	// 	sseo := *sse
	// 	fmt.Printf("\n%s, %+v\n", "sse op", sseo)

	// 	if sseo.Status == sfntypes.SyncExecutionStatusFailed || sseo.Status == sfntypes.SyncExecutionStatusTimedOut {
	// 		err := fmt.Errorf("step function %s, execution %s, failed with status %s. error code: %s. cause: %s. ", *sseo.StateMachineArn, *sseo.ExecutionArn, sseo.Status, *sseo.Error, *sseo.Cause)
	// 		return callErr(err)
	// 	}
	// }

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

func (data master) getAnswersMap() master {

	for _, v := range data.Players {
		data.Answers[v.Answer] = append(data.Answers[v.Answer], v.PlayerID)
	}

	return data
}

func (data master) getScoresMap() master {

	for k, v := range data.Answers {
		switch {
		case len(k) < 2:
			for _, id := range v {
				data.Scores[id] = 0
			}

		case len(v) > 2:
			for _, id := range v {
				data.Scores[id] = 1

			}

		case len(v) == 2:
			for _, id := range v {
				data.Scores[id] = 3
			}

		default:
			for _, id := range v {
				data.Scores[id] = 0
			}
		}
	}

	return data
}

func (data master) updateScores() master {

	for i, p := range data.Players {
		p.Score = p.Score + data.Scores[p.PlayerID]
		p.Answer = ""
		data.Players[i] = p
	}

	return data
}

func (data master) checkHiScore() master {
	for _, p := range data.Players {
		if p.Score == data.Hi {
			data.Tied = true
		}
		if p.Score > data.Hi {
			data.Hi = p.Score
		}
	}

	return data
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
