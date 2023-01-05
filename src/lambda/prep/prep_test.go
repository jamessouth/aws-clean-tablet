package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
)

type mockBatchWriteItemAPI func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)

func (m mockBatchWriteItemAPI) BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
	return m(ctx, params, optFns...)
}

func TestBatchWriteItem(t *testing.T) {
	// t.Skip()
	ctx := context.TODO()
	tableName := "myTable"
	requestItems := []types.WriteRequest{
		{
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
					"pk": &types.AttributeValueMemberS{Value: liveGame},
					"sk": &types.AttributeValueMemberS{Value: "234"},
				},
			},
		}}

	for i, tt := range batchWriteItemNoErrorTests {
		t.Run(fmt.Sprintf("no error test %d", i), func(t *testing.T) {
			if act, err := batchWriteItem(ctx, tt.client(t), requestItems, tableName); assert.NoErrorf(t, err, "FAIL - batchWriteItem - %s\n err: %+v\n", tt.description, err) {
				assert.Equalf(t, act, tt.exp, "FAIL - batchWriteItem - %s\n act: %s\n exp: %s\n", tt.description, act, tt.exp)
			}
		})
	}

	for i, tt := range batchWriteItemErrorTests {
		t.Run(fmt.Sprintf("error test %d", i), func(t *testing.T) {
			if act, err := batchWriteItem(ctx, tt.client(t), requestItems, tableName); assert.EqualErrorf(t, err, tt.msg, "FAIL - batchWriteItem - %s\n err: %+v\n", tt.description, err) {
				assert.Equalf(t, act, dynamodb.BatchWriteItemOutput{}, "FAIL - batchWriteItem - %s\n act: %s\n exp: \n", tt.description, act)
			}
		})
	}
}

func TestHandleUnprocessedItems(t *testing.T) {
	tableName := "myTable"
	batchWriteOutput := dynamodb.BatchWriteItemOutput{
		ConsumedCapacity:      []types.ConsumedCapacity{},
		ItemCollectionMetrics: map[string][]types.ItemCollectionMetrics{},
		UnprocessedItems: map[string][]types.WriteRequest{tableName: {
			{
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
						"pk": &types.AttributeValueMemberS{Value: liveGame},
						"sk": &types.AttributeValueMemberS{Value: "234"},
					},
				},
			}}},
		ResultMetadata: middleware.Metadata{},
	}

	for i, tt := range handleUnprocessedItemsNoErrorTests {
		t.Run(fmt.Sprintf("no error test %d", i), func(t *testing.T) {
			err := handleUnprocessedItems(tt.ctx, tt.client(t), batchWriteOutput, tableName)
			assert.NoErrorf(t, err, "FAIL - handleUnprocessedItems - %s\n err: %+v\n", tt.description, err)
		})
	}

	for i, tt := range handleUnprocessedItemsErrorTests {
		t.Run(fmt.Sprintf("error test %d", i), func(t *testing.T) {
			err := handleUnprocessedItems(context.TODO(), tt.client(t), batchWriteOutput, tableName)
			assert.EqualErrorf(t, err, tt.msg, "FAIL - handleUnprocessedItems - %s\n err: %+v\n", tt.description, err)
		})
	}
}
