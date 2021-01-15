package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go / build-lambda-zip.exe -o main.zip main / sam local invoke AuthorizerFunction -e ../event.json

func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	fmt.Println("con: ", ctx)
	fmt.Println("req: ", req)

	val := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:        "koko",
		PolicyDocument:     events.APIGatewayCustomAuthorizerPolicy{Version: "2012-10-17", Statement: []events.IAMPolicyStatement{{Effect: "Allow", Action: []string{"execute-api:Invoke"}, Resource: []string{"arn:aws:execute-api:us-east-2:270222239701:m0i5yb937f/dev/$connect"}}}},
		Context:            map[string]interface{}{},
		UsageIdentifierKey: "",
	}

	return val, nil
}

func main() {
	lambda.Start(handler)
}
