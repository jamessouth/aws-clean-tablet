package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/smithy-go"
)

type sfnEvent struct {
	Region, Endpoint, Word, Token string
	Conns                         []int
	Index                         int
}

func handler(ctx context.Context, req sfnEvent) (int, error) {

	fmt.Println("plaaaaaaay", req)

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == apigatewaymanagementapi.ServiceID && region == req.Region {
			ep := aws.Endpoint{
				PartitionID:   "aws",
				URL:           req.Endpoint,
				SigningRegion: req.Region,
			}

			return ep, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Region),
		// config.WithLogger(logger),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		return callErr(err)
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)
	apigwsvc := apigatewaymanagementapi.NewFromConfig(cfg)

	gameItemKey, err := attributevalue.MarshalMap(Key{
		Pk: "GAME",
		Sk: body.Gameno,
	})
	if err != nil {
		return callErr(err)
	}

	var game game
	err = attributevalue.UnmarshalMap(gi.Item, &game)
	if err != nil {
		return callErr(err)
	}

	fmt.Printf("%s%+v\n", "gammmmme ", game)

	ans, err := attributevalue.MarshalList(answer{
		PlayerID: id,
		Answer:   body.Answer,
	})
	if err != nil {
		return callErr(err)
	}

	ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                 gameItemKey,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("size (#AN) < :c"),
		ExpressionAttributeNames: map[string]string{
			// "#PL": "players",
			// "#ID": id,
			"#AN": "answers",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberL{Value: ans},
			":c": &types.AttributeValueMemberN{Value: body.PlayersCount},
		},
		UpdateExpression: aws.String("SET #AN = list_append(#AN, :a)"),
		ReturnValues:     types.ReturnValueAllNew,
	})

	if err != nil {

		return callErr(err)

	}

	for _, v := range gp.Game.Players {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

		_, err = apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return callErr(err)
		}

	}

	return 3, nil

}

func main() {
	lambda.Start(handler)
}

func callErr(err error) (int, error) {

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

	return 5, err

}
