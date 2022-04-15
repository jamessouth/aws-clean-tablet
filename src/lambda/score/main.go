package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
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

type livePlayer struct {
	PlayerID string `json:"playerid"`
	Name     string `json:"name"`
	ConnID   string `json:"connid"`
	Color    string `json:"color"`
	Index    string `json:"index"`
	Score    int    `json:"score"`
	Answer   string `json:"answer"`
}

type game struct {
	Players []livePlayer `dynamodbav:"players"`
	Answers map[string][]string
	Scores  map[string]int
	Winner  string
}

const (
	zeroPoints int = iota
	onePoint
	twoPoints
	threePoints
	winThreshold int = 5
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

	var (
		tableName = os.Getenv("tableName")
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		body      struct {
			Game struct {
				Sk      string
				Players []livePlayer
			}
		}
	)

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	fmt.Printf("%s%+v\n", "scrsss ", body)

	updatedScoreData := game{
		Players: body.Game.Players,
		Answers: map[string][]string{},
		Scores:  map[string]int{},
		Winner:  "",
	}.getAnswersMap().getScoresMap().updateScoresAndClearAnswers().getWinner()

	marshalledPlayersMap, err := attributevalue.Marshal(updatedScoreData.Players)
	if err != nil {
		return callErr(err)
	}

	scoreInput := dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGME"},
			"sk": &types.AttributeValueMemberS{Value: body.Game.Sk},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "players",
			"#W": "winner",
			"#S": "showAnswers",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p": marshalledPlayersMap,
			":w": &types.AttributeValueMemberS{Value: updatedScoreData.Winner},
			":f": &types.AttributeValueMemberBOOL{Value: false},
		},
		UpdateExpression: aws.String("SET #P = :p, #W = :w, #S = :f"),
	}

	if updatedScoreData.Winner != "" {
		scoreInput.ReturnValues = types.ReturnValueAllOld
		ui, err := ddbsvc.UpdateItem(ctx, &scoreInput)

		if err != nil {
			return callErr(err)
		}

		var gameIDs struct {
			Ids map[string]string
		}

		err = attributevalue.UnmarshalMap(ui.Attributes, &gameIDs)
		if err != nil {
			return callErr(err)
		}

		for _, p := range updatedScoreData.Players {

			won := ":z"

			if p.Name == updatedScoreData.Winner {
				won = ":o"
			}

			ue := fmt.Sprintf("ADD #W %s, #G :o, #T :t", won)

			_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key: map[string]types.AttributeValue{
					"pk": &types.AttributeValueMemberS{Value: "STAT"},
					"sk": &types.AttributeValueMemberS{Value: gameIDs.Ids[p.ConnID+p.Color+p.Name]},
				},
				TableName: aws.String(tableName),
				ExpressionAttributeNames: map[string]string{
					"#W": "wins",
					"#G": "games",
					"#T": "totalPoints",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":z": &types.AttributeValueMemberN{Value: "0"},
					":o": &types.AttributeValueMemberN{Value: "1"},
					":t": &types.AttributeValueMemberN{Value: strconv.Itoa(p.Score)},
				},
				UpdateExpression: aws.String(ue),
			})
			if err != nil {
				return callErr(err)
			}

		}

	}

	_, err = ddbsvc.UpdateItem(ctx, &scoreInput)

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

func sortByScore(players []livePlayer) []livePlayer {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	return players
}

func sortByIndex(players []livePlayer) []livePlayer {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Index < players[j].Index
	})

	return players
}

func (data game) getWinner() game {
	sortedScores := sortByScore(data.Players)

	if sortedScores[0].Score != sortedScores[1].Score && sortedScores[0].Score > winThreshold {
		data.Winner = sortedScores[0].Name
	}

	sortByIndex(data.Players)

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
