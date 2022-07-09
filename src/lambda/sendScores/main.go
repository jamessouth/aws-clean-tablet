package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type livePlayer struct {
	Name   string `json:"name"`
	ConnID string `json:"connid"`
	Color  string `json:"color"`
	Score  *int   `json:"score,omitempty"`
	Answer string `json:"answer,omitempty"`
	// HasAnswered     bool   `json:"hasAnswered,omitempty"`
	// PointsThisRound *int `json:"pointsThisRound,omitempty"`
}

type players struct {
	Players     []livePlayer `json:"players"`
	Sk          string       `json:"sk"`
	ShowAnswers bool         `json:"showAnswers"`
	Winner      string       `json:"winner"`
}

type stat struct {
	PlayerID string `json:"playerid"`
	Name string `json:"name"`
	Wins string `json:"wins"`
	Games string `json:"games"`
	Points string `json:"points"`
}

type output struct {
	Gameno      string                        `json:"gameno"`
	Players     events.DynamoDBAttributeValue `json:"players"`
	StatsList []stat                  `json:"statsList"`
	Winner      string                        `json:"winner"`
}

func getStats(players map[string]livePlayer, playersList []livePlayer) (res []stat) {
	var s stat

	for k, v := range players {
		
	}

	for i, p := playersList {

	}



	return
}

func updateScores(players map[string]livePlayer, scores map[string]int) (res []livePlayer, plrs events.DynamoDBAttributeValue, stats []stat) {
	m := map[string]events.DynamoDBAttributeValue{}

	for k, v := range players {
		score := *v.Score + scores[v.ConnID]
		v.Score = &score
		v.Answer = ""
		res = append(res, v)

		s := stat{
			PlayerID: k,
			Name:     v.Name,
			Wins:     "0",
			Games:    "1",
			Points:   "",
		}
		stats = append(stats, s)

		p := map[string]events.DynamoDBAttributeValue{
			"name":   events.NewStringAttribute(v.Name),
			"connid": events.NewStringAttribute(v.ConnID),
			"color":  events.NewStringAttribute(v.Color),
			"answer": events.NewStringAttribute(""),
			"score":  events.NewNumberAttribute(strconv.Itoa(score)),
		}

		marshalledPlayer := events.NewMapAttribute(p)

		m[k] = marshalledPlayer
	}
	plrs = events.NewMapAttribute(m)

	return
}

func sortByScoreThenName(players []livePlayer) {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case *players[i].Score != *players[j].Score:
			return *players[i].Score > *players[j].Score
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func getWinner(players []livePlayer) string {
	if *players[0].Score != *players[1].Score && *players[0].Score > winThreshold {
		return players[0].Name
	}

	return ""
}

const (
	winThreshold int = 5
)

func handler(ctx context.Context, req struct {
	Payload struct {
		Gameno, TableName, Endpoint, Region string
		Players                             map[string]livePlayer
		Scores                              map[string]int
	}
}) (output, error) {

	fmt.Printf("%s%+v\n", "sent req ", req)

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == apigatewaymanagementapi.ServiceID && region == req.Payload.Region {
			ep := aws.Endpoint{
				PartitionID:   "aws",
				URL:           req.Payload.Endpoint,
				SigningRegion: req.Payload.Region,
			}

			return ep, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	apigwcfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Payload.Region),
		// config.WithLogger(logger),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		return output{}, err
	}

	apigwsvc := apigatewaymanagementapi.NewFromConfig(apigwcfg)

	playersList, marshalledPlayersMap := updateScores(req.Payload.Players, req.Payload.Scores)

	sortByScoreThenName(playersList)
	winner := getWinner(playersList)

	payload, err := json.Marshal(players{
		Players:     playersList,
		Sk:          req.Payload.Gameno,
		ShowAnswers: false,
		Winner:      winner,
	})
	if err != nil {
		return output{}, err
	}

	for _, v := range playersList {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

		_, err = apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return output{}, err
		}
	}

	return output{
		Gameno:      req.Payload.Gameno,
		Players:     marshalledPlayersMap,
		StatsList: playersList,
		Winner:      winner,
	}, nil

}

func main() {
	lambda.Start(handler)
}
