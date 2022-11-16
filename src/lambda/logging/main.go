package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func checkLength(s string) error {
	if len(s) > 3200 {
		return errors.New("improper json input - too long")
	}

	return nil
}

func getReturnValue(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	bod := req.Body

	err := checkLength(bod)
	if err != nil {
		fmt.Println(err)
		return getReturnValue(http.StatusBadRequest), err
	}

	fmt.Printf("%s: %s\n", "Front end error", bod)

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
}
