package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

type output struct {
	Gameno string `json:"gameno"`
	Winner string `json:"winner"`
}

func getSlice(m map[string]livePlayer) (res []livePlayer) {
	for _, v := range m {
		res = append(res, v)
	}

	return
}

func updateScores(players map[string]livePlayer, scores map[string]int) map[string]livePlayer {
	for k, v := range players {
		score := *v.Score + scores[v.ConnID]
		v.Score = &score
		v.Answer = ""
		players[k] = v
	}

	return players
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

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Payload.Region),
	)
	if err != nil {
		return output{}, err
	}

	var ddbsvc = dynamodb.NewFromConfig(cfg)

	plsMap := updateScores(req.Payload.Players, req.Payload.Scores)
	plsSlice := getSlice(plsMap)
	sortByScoreThenName(plsSlice)
	winner := getWinner(plsSlice)

	payload, err := json.Marshal(players{
		Players:     plsSlice,
		Sk:          req.Payload.Gameno,
		ShowAnswers: false,
		Winner:      winner,
	})
	if err != nil {
		return output{}, err
	}

	for _, v := range plsSlice {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

		_, err = apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return output{}, err
		}
	}

	if winner == "" {

		marshalledPlayersMap, err := attributevalue.Marshal(plsMap)
		if err != nil {
			return output{}, err
		}

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
				"sk": &types.AttributeValueMemberS{Value: req.Payload.Gameno},
			},
			TableName: aws.String(req.Payload.TableName),
			ExpressionAttributeNames: map[string]string{
				"#P": "players",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":l": marshalledPlayersMap,
			},
			UpdateExpression: aws.String("SET #P = :l"),
		})

		if err != nil {
			return output{}, err
		}
	}

	return output{
		Gameno: req.Payload.Gameno,
		Winner: winner,
	}, nil

}

func main() {
	lambda.Start(handler)
}
