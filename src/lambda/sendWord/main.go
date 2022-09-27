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

type livePlayer struct {
	Name   string `json:"name"`
	ConnID string `json:"connid"`
	Color  string `json:"color"`
	Score  *int   `json:"score,omitempty"`
	Answer string `json:"answer,omitempty"`
}

func handler(ctx context.Context, req struct {
	Token, Gameno, TableName, Endpoint, Region string
}) error {

	fmt.Printf("%s%+v\n", "sent req ", req)

	reg := req.Region

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == apigatewaymanagementapi.ServiceID && region == reg {
			ep := aws.Endpoint{
				PartitionID:   "aws",
				URL:           req.Endpoint,
				SigningRegion: reg,
			}

			return ep, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	apigwcfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
		// config.WithLogger(logger),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return err
	}

	var (
		apigwsvc = apigatewaymanagementapi.NewFromConfig(apigwcfg)
		ddbsvc   = dynamodb.NewFromConfig(cfg)
	)

	ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
			"sk": &types.AttributeValueMemberS{Value: req.Gameno},
		},
		TableName: aws.String(req.TableName),
		ExpressionAttributeNames: map[string]string{
			"#T": "token",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":t": &types.AttributeValueMemberS{Value: req.Token},
		},
		UpdateExpression: aws.String("set #T = :t remove wordList[0]"),
		ReturnValues:     types.ReturnValueAllOld,
	})
	if err != nil {
		return err
	}

	var words struct {
		WordList []string
		Players  map[string]livePlayer
	}

	err = attributevalue.UnmarshalMap(ui.Attributes, &words)
	if err != nil {
		return err
	}

	fmt.Printf("%s%+v\n", "words ", words)

	current, next := words.WordList[0], words.WordList[1]

	if next == "game over" {
		_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
				"sk": &types.AttributeValueMemberS{Value: req.Gameno},
			},
			TableName: aws.String(req.TableName),
			ExpressionAttributeNames: map[string]string{
				"#L": "lastword",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":l": &types.AttributeValueMemberBOOL{Value: true},
			},
			UpdateExpression: aws.String("set #L = :l"),
		})
		if err != nil {
			return err
		}
	}

	payload, err := json.Marshal(struct {
		Word string `json:"newword"`
	}{
		Word: current,
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
