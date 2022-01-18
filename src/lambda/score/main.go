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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/smithy-go"
)

// type answer struct {
// 	PlayerID string `json,dynamodbav:"playerid"`
// 	Answer   string `json,dynamodbav:"answer"`
// }

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
	Sk       string         `json:"sk"`
	Players  livePlayerList `json:"players"`
	HiScore  int            `json:"hiScore"`
	GameTied bool           `json:"gameTied"`
}

type body struct {
	Game liveGame `json:"game"`
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

	updatedPayersList := updateScores(body.Game.Players)

	hiScore, gameTied := checkHiScore(updatedPayersList, body.Game.HiScore, body.Game.GameTied)

	if !gameTied && hiScore > winThreshold {
		winner = true
	}

	marshalledPlayersMap, err := attributevalue.Marshal(updatedPayersList)
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
			"#H": "hiScore",
			"#G": "gameTied",
			"#W": "winner",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p": marshalledPlayersMap,
			":h": &types.AttributeValueMemberN{Value: strconv.Itoa(updatedGame.HiScore)},
			":g": &types.AttributeValueMemberBOOL{Value: updatedGame.GameTied},
			":w": &types.AttributeValueMemberBOOL{Value: winner},
		},
		UpdateExpression: aws.String("SET #P = :p, #H = :h, #G = :g, #W = :w"),
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

func getAnswersMap(players livePlayerList) (res map[string][]string) {
	res = make(map[string][]string)

	for _, v := range players {
		res[v.Answer] = append(res[v.Answer], v.PlayerID)
	}

	return
}

func getScoresMap(answers map[string][]string) (res map[string]int) {
	res = make(map[string]int)

	for k, v := range answers {
		switch {
		case len(k) < 2:
			for _, id := range v {
				res[id] = 0
			}

		case len(v) > 2:
			for _, id := range v {
				res[id] = 1

			}

		case len(v) == 2:
			for _, id := range v {
				res[id] = 3
			}

		default:
			for _, id := range v {
				res[id] = 0
			}
		}
	}

	return
}

func updateScores(players livePlayerList) (res livePlayerList) {
	answers := getAnswersMap(players)
	scores := getScoresMap(answers)

	for _, p := range players {
		res = append(res, p.adjScore(scores[p.PlayerID]))
	}

	return
}

func checkHiScore(players livePlayerList, hiScore int, tied bool) (hi int, tie bool) {
	for _, p := range players {
		if p.Score == hiScore {
			tie = true
		}
		if p.Score > hiScore {
			// fmt.Println("hiscore", p.Score, hiScore, tied)
			hi = p.Score
			tie = false
		}
	}

	return
}

// , hiScore int, tied bool    , int, bool
func (old livePlayer) adjScore(incr int) livePlayer {
	// score := old.Score + incr

	// hs, t := checkHiScore(score, hiScore, tied)

	return livePlayer{
		PlayerID: old.PlayerID,
		Name:     old.Name,
		ConnID:   old.ConnID,
		Color:    old.Color,
		Score:    old.Score + incr,
		Answer:   "",
	}

	// , hs, t
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
