package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go / build-lambda-zip.exe -o main.zip main / sam local invoke AuthorizerFunction -e ../event.json

func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	fmt.Printf("%s: %+v\n", "context", ctx)
	fmt.Printf("%s: %+v\n", "request", req)

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:    "koko",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{Version: "2012-10-17", Statement: []events.IAMPolicyStatement{{Effect: "Allow", Action: []string{"execute-api:Invoke"}, Resource: []string{req.MethodArn}}}},
		Context: map[string]interface{}{
			"a": 1,
			"b": 2,
		},
		UsageIdentifierKey: "",
	}, nil
}

func main() {
	lambda.Start(handler)
}
