package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

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

type key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// type player struct {
// 	Name   string `dynamodbav:"name"`
// 	ConnID string `dynamodbav:"connid"`
// 	Ready  bool   `dynamodbav:"ready"`
// 	Color  string `dynamodbav:"color,omitempty"`
// 	Score  int    `dynamodbav:"score"`
// 	Answer answer `dynamodbav:"answer"`
// }

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
	// SendToFront  bool          `dynamodbav:"sendToFront"`
	HiScore  int  `dynamodbav:"hiScore"`
	GameTied bool `dynamodbav:"gameTied"`
}

// type body struct {
// 	Gameno string `json:"gameno"`
// 	Answer string `json:"answer"`
// }

// func (pm livePlayerMap) assignColors() livePlayerMap {
// 	count := 0
// 	for k, v := range pm {
// 		v.Color = colors[count]
// 		pm[k] = v
// 		count++
// 	}

// 	return pm
// }

// func (pm livePlayerMap) mapToSlice() (res []sfnArrInput) {
// 	for k, v := range pm {
// 		res = append(res, sfnArrInput{
// 			Id:   k,
// 			Name: v.Name,
// 		})
// 	}

// 	return
// }

type sfnEvent struct {
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
	TableName string `json:"tableName"`
	Gameno    string `json:"gameno"`
}

func handler(ctx context.Context, req sfnEvent) error {

	fmt.Printf("%s%+v\n", "score req ", req)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Region),
	)
	if err != nil {
		return callErr(err)
	}

	sfnarn, ok := os.LookupEnv("SFNARN")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find sfn arn"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)
	sfnsvc := sfn.NewFromConfig(cfg)

	// id := req.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	// var body body

	// err = json.Unmarshal([]byte(req.Body), &body)
	// if err != nil {
	// 	return callErr(err)
	// }

	gameItemKey, err := attributevalue.MarshalMap(key{
		Pk: "LIVEGME",
		Sk: req.Gameno,
	})
	if err != nil {
		return callErr(err)
	}

	gm, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       gameItemKey,
		TableName: aws.String(req.TableName),
	})

	if err != nil {
		return callErr(err)
	}

	var game liveGame
	err = attributevalue.UnmarshalMap(gm.Item, &game)
	if err != nil {
		return callErr(err)
	}

	winner := false
	const winThreshold int = 24

	updatedGame := updateScores(game)

	if !updatedGame.GameTied && updatedGame.HiScore > winThreshold {
		winner = true
	}

	marshalledPlayersMap, err := attributevalue.Marshal(updatedGame.Players)
	if err != nil {
		return callErr(err)
	}

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:       gameItemKey,
		TableName: aws.String(req.TableName),
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

		sfnInput := "{\"gameno\":\"" + req.Gameno + "\"}"

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

	return nil
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

func callErr(err error) error {

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

	return err

}
