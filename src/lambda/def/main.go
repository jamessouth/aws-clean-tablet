package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func getReturnValue(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}
}

func handler(req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(req.Body) > 99 {
		fmt.Printf("%s: %+v\n", "body", req.Body[:99])

		return getReturnValue(http.StatusBadRequest), errors.New("improper json input - too long")
	}

	fmt.Printf("%s: %+v\n", "request", req)

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
}
