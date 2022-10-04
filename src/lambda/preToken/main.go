package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

func panicProtectedUID() (id string) {
	defer func() {
		if err := recover(); err != nil {
			id = fmt.Sprintf("uuid panic: %v", err)
		}
	}()

	return uuid.NewString()
}

func handler(ctx context.Context, ev events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {

	uid := panicProtectedUID()

	if uid[0] == 117 { //"u"
		return ev, errors.New(uid)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(ev.CognitoEventUserPoolsHeader.Region),
	)
	if err != nil {
		return ev, err
	}

	var (
		tableName = os.Getenv("tableName")
		ddbsvc    = dynamodb.NewFromConfig(cfg)
	)

	_, err = ddbsvc.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "TOKEN#" + ev.Request.UserAttributes["sub"]},
			"sk": &types.AttributeValueMemberS{Value: uid},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return ev, err
	}

	ev.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride = map[string]string{
		"q": uid,
	}

	return ev, nil
}

func main() {
	lambda.Start(handler)
}
