package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func handler(ctx context.Context, req struct {
	Payload struct {
		Gameno, TableName, Region string
	}
}) error {

	fmt.Printf("%s%+v\n", "stat req ", req)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Payload.Region),
	)
	if err != nil {
		return err
	}

	var ddbsvc = dynamodb.NewFromConfig(cfg)

	gi, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
			"sk": &types.AttributeValueMemberS{Value: req.Payload.Gameno},
		},
		TableName: aws.String(req.Payload.TableName),
	})
	if err != nil {
		return err
	}

	var gameRecord players
	err = attributevalue.UnmarshalMap(gi.Item, &gameRecord)
	if err != nil {
		return err
	}

	fmt.Printf("%s%+v\n", "unmarshalledGame ", gameRecord)

	if updatedScoreData.Winner != "" {

		var gameIDs struct {
			Ids map[string]string
		}

		err = attributevalue.UnmarshalMap(ui.Attributes, &gameIDs)
		if err != nil {

		}

		for _, p := range updatedScoreData.Players {

			won := ":z"
			eav := map[string]types.AttributeValue{
				":z": &types.AttributeValueMemberN{Value: "0"},
				":o": &types.AttributeValueMemberN{Value: "1"},
				":t": &types.AttributeValueMemberN{Value: strconv.Itoa(p.Score)},
			}

			if p.Name == updatedScoreData.Winner {
				won = ":o"
				eav = map[string]types.AttributeValue{
					":o": &types.AttributeValueMemberN{Value: "1"},
					":t": &types.AttributeValueMemberN{Value: strconv.Itoa(p.Score)},
				}
			}

			ue := fmt.Sprintf("ADD #W %s, #G :o, #T :t", won)

			_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key: map[string]types.AttributeValue{
					"pk": &types.AttributeValueMemberS{Value: "STAT"},
					"sk": &types.AttributeValueMemberS{Value: gameIDs.Ids[p.PlayerID]},
				},
				TableName: aws.String(req.Payload.TableName),
				ExpressionAttributeNames: map[string]string{
					"#W": "wins",
					"#G": "games",
					"#T": "points",
				},
				ExpressionAttributeValues: eav,
				UpdateExpression:          aws.String(ue),
			})
			if err != nil {

			}

		}

	}

}

func main() {
	lambda.Start(handler)
}
