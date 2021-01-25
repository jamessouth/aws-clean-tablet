package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go / build-lambda-zip.exe -o main.zip main / sam local invoke AuthorizerFunction -e ../event.json

func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	// fmt.Printf("%s: %+v\n", "request", req.QueryStringParameters["auth"])

	token := []byte(req.QueryStringParameters["auth"])

	userPoolID, ok := os.LookupEnv("userPoolID")
	if !ok {
		panic("cannot find user pool id")
	}
	appClientID, ok := os.LookupEnv("appClientID")
	if !ok {
		panic("cannot find app client id")
	}

	region := strings.Split(req.MethodArn, ":")[3]

	keyset, err := jwk.Fetch("https://cognito-idp." + region + ".amazonaws.com/" + userPoolID + "/.well-known/jwks.json")
	if err != nil {

	}

	parsedToken, err := jwt.Parse(
		bytes.NewReader(token),
		jwt.WithKeySet(keyset),
		jwt.WithValidate(true),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithClaimValue("client_id", appClientID),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(parsedToken)
	fmt.Println(parsedToken.Subject())

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
