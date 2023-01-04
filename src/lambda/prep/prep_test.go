package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockBatchWriteItemAPI func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)

func (m mockBatchWriteItemAPI) BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
	return m(ctx, params, optFns...)
}

func TestBatchWriteItem(t *testing.T) {

	for i, tt := range batchWriteItemTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := batchWriteItem(ctx, tt.client(t), tt.requestItems, tt.tableName)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

		})
	}
}
