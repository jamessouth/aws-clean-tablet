package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
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
		PlayersList               []livePlayer
	}
}) error {

	fmt.Printf("%s%+v\n", "stat req ", req)

	return nil

}

func main() {
	lambda.Start(handler)
}
