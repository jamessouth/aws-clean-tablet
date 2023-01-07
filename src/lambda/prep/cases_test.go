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

var (
	item1 = types.WriteRequest{
		PutRequest: &types.PutRequest{
			Item: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: listGame},
				"sk": &types.AttributeValueMemberS{Value: "123"},
			},
		},
	}
	item2 = types.WriteRequest{
		PutRequest: &types.PutRequest{
			Item: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: liveGame},
				"sk": &types.AttributeValueMemberS{Value: "234"},
			},
		},
	}
)

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
					UnprocessedItems: map[string][]types.WriteRequest{"myTable": {item1, item2}},
				}, nil
			})
		},
		exp: dynamodb.BatchWriteItemOutput{
			UnprocessedItems: map[string][]types.WriteRequest{"myTable": {item1, item2}},
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
						UnprocessedItems: map[string][]types.WriteRequest{"myTable": {item1, item2}},
					}, nil
				} else if cval < 2001 {
					return &dynamodb.BatchWriteItemOutput{
						UnprocessedItems: map[string][]types.WriteRequest{"myTable": {item1}},
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
					UnprocessedItems: map[string][]types.WriteRequest{"myTable": {item1, item2}},
				}, nil
			})
		},
		msg:         "error: unable to write 2 items",
		description: "error - never able to write items",
	},
}

var zero = 0

var getLivePlayerMapTests = []struct {
	input1      map[string]listPlayer
	input2      stringSlice
	exp         map[string]livePlayer
	description string
}{
	{
		input1: map[string]listPlayer{
			"111": {Name: "will", ConnID: "3x33"},
			"222": {Name: "earl", ConnID: "11w1"},
			"333": {Name: "carl", ConnID: "22r2"},
			"444": {Name: "darlene", ConnID: "33n3"},
		},
		input2: stringSlice{"blue", "green", "purple", "red"},
		exp: map[string]livePlayer{
			"111": {Name: "will", ConnID: "3x33", Color: "blue", Answer: "", Score: &zero},
			"222": {Name: "earl", ConnID: "11w1", Color: "green", Answer: "", Score: &zero},
			"333": {Name: "carl", ConnID: "22r2", Color: "purple", Answer: "", Score: &zero},
			"444": {Name: "darlene", ConnID: "33n3", Color: "red", Answer: "", Score: &zero},
		},
		description: "get map of livePlayers",
	},
}
