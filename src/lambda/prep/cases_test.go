package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var batchWriteItemTests = []struct {
	client       func(t *testing.T) DdbBatchWriteItemAPI
	requestItems []types.WriteRequest
	tableName    string
}{
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()

				return &dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, nil
			})
		},

		requestItems: []types.WriteRequest{{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"pk": &types.AttributeValueMemberS{Value: listGame},
					"sk": &types.AttributeValueMemberS{Value: "123"},
				},
			},
		},
			{
				PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
						"sk": &types.AttributeValueMemberS{Value: "234"},
					},
				},
			}},

		tableName: "myTable",
	},
}
