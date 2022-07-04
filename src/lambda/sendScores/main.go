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

type connectUpdate struct {
	PlayerID string `json:"playerid"`
	Color    string `json:"color"`
	Index    string `json:"index"`
}

type livePlayer struct {
	connectUpdate
	Name   string `json:"name"`
	ConnID string `json:"connid"`
	Score  int    `json:"score"`
	Answer string `json:"answer"`
}

type players struct {
	Players     livePlayerList `json:"players"`
	Sk          string         `json:"sk"`
	ShowAnswers bool           `json:"showAnswers"`
	Winner      string         `json:"winner"`
}

type livePlayerList []struct {
	Name            string `json:"name"`
	ConnID          string `json:"connid"`
	Color           string `json:"color"`
	Score           *int   `json:"score,omitempty"`
	Answer          string `json:"answer,omitempty"`
	HasAnswered     bool   `json:"hasAnswered,omitempty"`
	PointsThisRound *int   `json:"pointsThisRound,omitempty"`
}

type output struct {
	Gameno string `json:"gameno"`
	Winner string `json:"winner"`
}

func (players livePlayerList) updateScores() livePlayerList {
	for i, p := range players {
		score := *p.Score + *p.PointsThisRound
		p.Score = &score
		p.PointsThisRound = nil
		players[i] = p
	}

	return players
}

func (players livePlayerList) sortByScoreThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case *players[i].Score != *players[j].Score:
			return *players[i].Score > *players[j].Score
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) getWinner() string {
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

	gi, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
			"sk": &types.AttributeValueMemberS{Value: req.Payload.Gameno},
		},
		TableName: aws.String(req.Payload.TableName),
	})
	if err != nil {
		return output{}, err
	}

	var gameRecord players
	err = attributevalue.UnmarshalMap(gi.Item, &gameRecord)
	if err != nil {
		return output{}, err
	}

	fmt.Printf("%s%+v\n", "unmarshalledGame ", gameRecord)

	pls := gameRecord.Players.updateScores()
	pls.sortByScoreThenName()
	winner := pls.getWinner()

	payload, err := json.Marshal(players{
		Players:     pls,
		Sk:          gameRecord.Sk,
		ShowAnswers: false,
		Winner:      winner,
	})
	if err != nil {
		return output{}, err
	}

	for _, v := range gameRecord.Players {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

		_, err = apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return output{}, err
		}
	}

	if winner == "" {

		marshalledPlayersList, err := attributevalue.Marshal(pls)
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
				":l": marshalledPlayersList,
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
