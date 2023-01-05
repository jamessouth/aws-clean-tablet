package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ctxKey string

var batchWriteItemNoErrorTests = []struct {
	client      func(t *testing.T) DdbBatchWriteItemAPI
	exp         dynamodb.BatchWriteItemOutput
	description string
}{
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()

				return &dynamodb.BatchWriteItemOutput{}, nil
			})
		},
		exp:         dynamodb.BatchWriteItemOutput{},
		description: "no unprocessed items, no error",
	},
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()
				return &dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{"myTable": {{
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
				}, nil
			})
		},
		exp: dynamodb.BatchWriteItemOutput{
			UnprocessedItems: map[string][]types.WriteRequest{"myTable": {{
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
		},
		description: "unprocessed items, no error",
	},
}

var batchWriteItemErrorTests = []struct {
	client           func(t *testing.T) DdbBatchWriteItemAPI
	msg, description string
}{
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()
				return &dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{},
				}, errors.New("fake error")
			})
		},
		msg:         "fake error",
		description: "api error",
	},
}

var handleUnprocessedItemsNoErrorTests = []struct {
	client      func(t *testing.T) DdbBatchWriteItemAPI
	ctx         context.Context
	description string
}{
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()

				return &dynamodb.BatchWriteItemOutput{}, nil
			})
		},
		ctx:         context.TODO(),
		description: "unprocessed items are successfully processed the first loop iteration, no error",
	},
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()
				var mk ctxKey = "cKey"
				cval, _ := strconv.Atoi(fmt.Sprintf("%v", ctx.Value(mk)))
				if cval < 1001 {
					return &dynamodb.BatchWriteItemOutput{
						UnprocessedItems: map[string][]types.WriteRequest{"myTable": {{
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
					}, nil
				} else if cval < 2001 {
					return &dynamodb.BatchWriteItemOutput{
						UnprocessedItems: map[string][]types.WriteRequest{"myTable": {{
							PutRequest: &types.PutRequest{
								Item: map[string]types.AttributeValue{
									"pk": &types.AttributeValueMemberS{Value: listGame},
									"sk": &types.AttributeValueMemberS{Value: "123"},
								},
							},
						}}},
					}, nil
				} else {
					return &dynamodb.BatchWriteItemOutput{}, nil
				}
			})
		},
		ctx:         context.Background(),
		description: "items eventually processed before loop ends, no error",
	},
}

var handleUnprocessedItemsErrorTests = []struct {
	client           func(t *testing.T) DdbBatchWriteItemAPI
	msg, description string
}{
	{
		client: func(t *testing.T) DdbBatchWriteItemAPI {
			return mockBatchWriteItemAPI(func(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
				t.Helper()
				return &dynamodb.BatchWriteItemOutput{
					UnprocessedItems: map[string][]types.WriteRequest{"myTable": {{
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
				}, nil
			})
		},
		msg:         "error: unable to write 2 items",
		description: "error - never able to write items",
	},
}
