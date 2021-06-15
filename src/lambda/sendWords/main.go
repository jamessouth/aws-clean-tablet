package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/smithy-go"
)

type sfnEvent struct {
	Region    string   `json:"region"`
	Endpoint  string   `json:"endpoint"`
	Word      string   `json:"word"`
	Gameno    string   `json:"gameno"`
	TableName string   `json:"tableName"`
	Conns     []string `json:"conns"`
	Token     string   `json:"token"`
}

func handler(ctx context.Context, req sfnEvent) error {

	fmt.Printf("%s%+v\n", "sndwords req ", req)

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

	ddbsvc := dynamodb.NewFromConfig(cfg)
	apigwsvc := apigatewaymanagementapi.NewFromConfig(cfg)

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "GAME"},
			"sk": &types.AttributeValueMemberS{Value: req.Gameno},
		},
		TableName: aws.String(req.TableName),
		// ConditionExpression: aws.String("size (#AN) < :c"),
		ExpressionAttributeNames: map[string]string{
			// "#PL": "players",
			// "#ID": id,
			"#T": "token",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":t": &types.AttributeValueMemberS{Value: req.Token},
		},
		UpdateExpression: aws.String("SET #T = :t)"),
		// ReturnValues:     types.ReturnValueAllNew,
	})

	if err != nil {
		return callErr(err)
	}

	for _, v := range req.Conns {

		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v), Data: []byte(req.Word)}

		_, err = apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return callErr(err)
		}

	}

	return nil

}

func main() {
	lambda.Start(handler)
}

func callErr(err error) error {

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

	return err

}
