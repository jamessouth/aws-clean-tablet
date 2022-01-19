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

type game struct {
	Players  livePlayerList `dynamodbav:"players"`
	Answers  map[string][]string
	Scores   map[string]int
	HiScore  int
	GameTied bool
}

const (
	zeroPoints int = iota
	onePoint
	twoPoints
	threePoints
	winThreshold int = 24
)

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

	ddbsvc := dynamodb.NewFromConfig(cfg)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	winner := false

	scoreData := game{
		Players:  body.Game.Players,
		Answers:  map[string][]string{},
		Scores:   map[string]int{},
		HiScore:  zeroPoints,
		GameTied: false,
	}

	updatedScoreData := scoreData.getAnswersMap().getScoresMap().updateScoresAndClearAnswers().getHiScoreAndTie()

	if !updatedScoreData.GameTied && updatedScoreData.HiScore > winThreshold {
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

func (data game) getAnswersMap() game {
	for _, v := range data.Players {
		data.Answers[v.Answer] = append(data.Answers[v.Answer], v.PlayerID)
	}

	return data
}

func (data game) getScoresMap() game {
	for k, v := range data.Answers {
		switch {
		case len(k) < 2:
			for _, id := range v {
				data.Scores[id] = zeroPoints
			}
		case len(v) > 2:
			for _, id := range v {
				data.Scores[id] = onePoint
			}
		case len(v) == 2:
			for _, id := range v {
				data.Scores[id] = threePoints
			}
		default:
			for _, id := range v {
				data.Scores[id] = zeroPoints
			}
		}
	}

	return data
}

func (data game) updateScoresAndClearAnswers() game {
	for i, p := range data.Players {
		p.Score += data.Scores[p.PlayerID]
		p.Answer = ""
		data.Players[i] = p
	}

	return data
}

func (data game) getHiScoreAndTie() game {
	for _, p := range data.Players {
		if p.Score == data.HiScore {
			data.GameTied = true
		}
		if p.Score > data.HiScore {
			data.HiScore = p.Score
			data.GameTied = false
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
