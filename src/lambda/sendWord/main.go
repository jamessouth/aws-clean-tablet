package main

import (
	"context"
	"encoding/json"
	"fmt"

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

	ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
			"sk": &types.AttributeValueMemberS{Value: req.Payload.Gameno},
		},
		TableName:        aws.String(req.Payload.TableName),
		UpdateExpression: aws.String("remove wordList[0]"),
		ReturnValues:     types.ReturnValueAllOld,
	})
	if err != nil {
		return err
	}

	var words struct {
		WordList []string
		Players  livePlayerList
	}

	err = attributevalue.UnmarshalMap(ui.Attributes, &words)
	if err != nil {
		return err
	}

	fmt.Printf("%s%+v\n", "words ", words)

	payload, err := json.Marshal(struct {
		Word string `json:"word"`
	}{
		Word: words.WordList[0],
	})
	if err != nil {
		return err
	}

	for _, v := range words.Players {

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
