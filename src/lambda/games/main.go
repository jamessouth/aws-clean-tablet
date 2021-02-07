package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.DynamoDBEvent) {
	fmt.Printf("%s: %+v\n", "cont", ctx)
	fmt.Printf("%s: %+v\n", "request", req)
}

func main() {
	lambda.Start(handler)
}
