package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, ev events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {

	fmt.Printf("%s: %+v\n", "ctx", ctx)

	fmt.Printf("%s: %+v\n", "ev", ev)

	ev.Response.AutoConfirmUser = false
	ev.Response.AutoVerifyEmail = false
	ev.Response.AutoVerifyPhone = false

	return ev, nil
}

func main() {
	lambda.Start(handler)
}
