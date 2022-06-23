package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

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

type plrs struct {
	Players livePlayerList `json:"players"`
}

type livePlayerList []struct {
	PlayerID        string `json:"playerid"`
	Name            string `json:"name"`
	ConnID          string `json:"connid"`
	Color           string `json:"color"`
	Score           int    `json:"score"`
	Index           string `json:"index"`
	Answer          string `json:"answer"`
	HasAnswered     bool   `json:"hasAnswered"`
	PointsThisRound string `json:"pointsThisRound"`
}

// type liveGame struct {
// 	Sk           string         `json:"sk"`
// 	Players      livePlayerList `json:"players"`
// 	CurrentWord  string         `json:"currentWord"`
// 	PreviousWord string         `json:"previousWord"`
// 	AnswersCount int            `json:"answersCount"`
// 	ShowAnswers  bool           `json:"showAnswers"`
// 	Winner       string         `json:"winner"`
// }

// type modifyLiveGamePayload struct {
// 	ModLiveGame liveGame
// }

// https://go.dev/play/p/CvniMWPoLKG
// func (p modifyLiveGamePayload) MarshalJSON() ([]byte, error) {
// 	if p.ModLiveGame.AnswersCount == len(p.ModLiveGame.Players) {
// 		return []byte(`null`), nil
// 	}

// 	if p.ModLiveGame.AnswersCount > 0 {
// 		for i, pl := range p.ModLiveGame.Players {
// 			if pl.HasAnswered {
// 				pl.Answer = ""
// 				p.ModLiveGame.Players[i] = pl
// 			}
// 		}
// 	}

// 	m, err := json.Marshal(p.ModLiveGame)
// 	if err != nil {
// 		return m, err
// 	}

// 	return []byte(fmt.Sprintf("{%q:%s}", "mdLveGm", m)), nil
// }

func (players livePlayerList) getPoints() livePlayerList {
	dist := map[string]int{}

	for _, v := range players {
		dist[v.Answer]++
	}

	for i, p := range players {
		if len(p.Answer) > 1 {
			freq := dist[p.Answer]
			if freq == 2 {
				p.PointsThisRound = strconv.Itoa(3)
			} else if freq > 2 {
				p.PointsThisRound = strconv.Itoa(1)
			} else {
				p.PointsThisRound = strconv.Itoa(0)
			}
		} else {
			p.PointsThisRound = strconv.Itoa(0)

		}
		players[i] = p
	}

	return players
}

func (players livePlayerList) sortByAnswerThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Answer != players[j].Answer:
			return players[i].Answer < players[j].Answer
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) sortByScoreThenName() livePlayerList {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Score != players[j].Score:
			return players[i].Score > players[j].Score
		default:
			return players[i].Name < players[j].Name
		}
	})

	return players
}

func handler(ctx context.Context, req struct {
	Payload struct {
		Region, Endpoint, Gameno, TableName string
	}
}) error {

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
		return err
	}

	apigwsvc := apigatewaymanagementapi.NewFromConfig(apigwcfg)

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

	var gameRecord plrs
	err = attributevalue.UnmarshalMap(gi.Item, &gameRecord)
	if err != nil {
		return err
	}

	fmt.Printf("%s%+v\n", "unmarshalledGame ", gameRecord)

	// if gameRecord.ShowAnswers {
	// 	pls.sortByAnswerThenName()
	// } else {
	// pls
	// }

	// gp := modifyLiveGamePayload{
	// 	ModLiveGame: liveGame{
	// 		Sk:           gameRecord.Sk,
	// 		Players:       ,
	// 		CurrentWord:  gameRecord.CurrentWord,
	// 		PreviousWord: gameRecord.PreviousWord,
	// 		// AnswersCount: gameRecord.AnswersCount,
	// 		// ShowAnswers:  gameRecord.ShowAnswers,
	// 		Winner: gameRecord.Winner,
	// 	},
	// }

	payload, err := json.Marshal(plrs{
		Players: gameRecord.Players.sortByScoreThenName().getPoints(),
	})
	if err != nil {
		return err
	}

	for _, v := range gameRecord.Players {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

		_, err = apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return err
		}
	}

	return nil

}

func main() {
	lambda.Start(handler)
}
