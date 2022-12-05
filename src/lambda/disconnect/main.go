package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(req.Body) > 99 {
		fmt.Printf("%s: %+v\n", "body", req.Body[:99])

		return events.APIGatewayProxyResponse{
			StatusCode:        http.StatusBadRequest,
			Headers:           map[string]string{"Content-Type": "application/json"},
			MultiValueHeaders: map[string][]string{},
			Body:              "",
			IsBase64Encoded:   false,
		}, errors.New("improper json input - too long")
	}

	fmt.Printf("%s: %+v\n", "Disconnected", req)

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
