package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func handler(ctx context.Context, req struct {
	Token, Gameno, TableName, Endpoint, Region string
}) error {

	fmt.Printf("%s%+v\n", "sent req ", req)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Region),
	)
	if err != nil {
		return err
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
			"sk": &types.AttributeValueMemberS{Value: req.Gameno},
		},
		TableName: aws.String(req.TableName),
		ExpressionAttributeNames: map[string]string{
			"#T": "endtoken",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":t": &types.AttributeValueMemberS{Value: req.Token},
		},
		UpdateExpression: aws.String("set #T = :t"),
	})
	if err != nil {
		return err
	}

	return nil

}

func main() {
	lambda.Start(handler)
}
