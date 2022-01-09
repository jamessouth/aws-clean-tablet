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

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"

	"github.com/aws/smithy-go"
)

type answer struct {
	PlayerID, Answer string
}

type livePlayer struct {
	Name   string `json:"name"`
	ConnID string `json:"connid"`
	Color  string `json:"color"`
	Score  int    `json:"score"`
	Answer answer `json:"answer"`
}

type livePlayerMap map[string]livePlayer

type liveGame struct {
	No       string        `json:"no"`
	Players  livePlayerMap `json:"players"`
	HiScore  int           `json:"hiScore"`
	GameTied bool          `json:"gameTied"`
}

type body struct {
	Game liveGame `json:"game"`
}

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

	sfnarn, ok := os.LookupEnv("SFNARN")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find sfn arn"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)
	sfnsvc := sfn.NewFromConfig(cfg)

	// id := req.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	winner := false
	const winThreshold int = 24

	updatedGame := updateScores(body.Game)

	if !updatedGame.GameTied && updatedGame.HiScore > winThreshold {
		winner = true
	}

	marshalledPlayersMap, err := attributevalue.Marshal(updatedGame.Players)
	if err != nil {
		return callErr(err)
	}

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGME"},
			"sk": &types.AttributeValueMemberS{Value: body.Game.No},
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

	if !winner {

		sfnInput := "{\"gameno\":\"" + body.Game.No + "\"}"

		ssei := sfn.StartSyncExecutionInput{
			StateMachineArn: aws.String(sfnarn),
			Input:           aws.String(sfnInput),
		}

		sse, err := sfnsvc.StartSyncExecution(ctx, &ssei)
		if err != nil {
			return callErr(err)
		}

		sseo := *sse
		fmt.Printf("\n%s, %+v\n", "sse op", sseo)

		if sseo.Status == sfntypes.SyncExecutionStatusFailed || sseo.Status == sfntypes.SyncExecutionStatusTimedOut {
			err := fmt.Errorf("step function %s, execution %s, failed with status %s. error code: %s. cause: %s. ", *sseo.StateMachineArn, *sseo.ExecutionArn, sseo.Status, *sseo.Error, *sseo.Cause)
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

func getAnswersMap(game liveGame) map[string][]string {
	res := make(map[string][]string)

	for _, v := range game.Players {
		// fmt.Printf("%s, %v, %+v\n", "anssss", k, v)
		res[v.Answer.Answer] = append(res[v.Answer.Answer], v.Answer.PlayerID)
	}

	return res
}

func updateScores(game liveGame) liveGame {
	answers := getAnswersMap(game)

	for k, v := range answers {
		// fmt.Printf("%s, %v, %+v\n", "anssssmapppp", k, v)

		switch {
		case len(k) < 2: // c.updateEachScore(v, 0)

		case len(v) > 2: // c.updateEachScore(v, 1)
			for _, id := range v {
				// fmt.Println("1st", id)
				game.Players[id], game.HiScore, game.GameTied = adjScore(game.Players[id], 1, game.HiScore, game.GameTied)
			}

		case len(v) == 2: // c.updateEachScore(v, 3)
			for _, id := range v {
				// fmt.Println("2nd", id)
				game.Players[id], game.HiScore, game.GameTied = adjScore(game.Players[id], 3, game.HiScore, game.GameTied)
			}

		default: // c.updateEachScore(v, 0)
		}
	}

	return game
}

func checkHiScore(score, hiScore int, tied bool) (int, bool) {
	if score == hiScore {
		tied = true
	}
	if score > hiScore {
		hiScore = score
		tied = false
	}

	return hiScore, tied
}

func adjScore(old livePlayer, incr, hiScore int, tied bool) (livePlayer, int, bool) {
	score := old.Score + incr

	hs, t := checkHiScore(score, hiScore, tied)

	return livePlayer{
		Name:   old.Name,
		ConnID: old.ConnID,
		Color:  old.Color,
		Score:  score,
		Answer: old.Answer,
	}, hs, t
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
