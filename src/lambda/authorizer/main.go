package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// $env:GOOS = "linux" / $env:CGO_ENABLED = "0" / $env:GOARCH = "amd64" / go build -o main main.go | build-lambda-zip.exe -o main.zip main / sam local invoke AuthorizerFunction -e ../event.json

func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	// fmt.Printf("%s: %+v\n", "request", req.QueryStringParameters["auth"])
	if req.Headers["Origin"] != "http://localhost:8000" {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong domain")
	}

	if len(req.Headers["User-Agent"]) < 10 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong client")
	}

	userPoolID, ok := os.LookupEnv("userPoolID")
	if !ok {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("internal error - cannot find user pool id")
	}

	appClientID, ok := os.LookupEnv("appClientID")
	if !ok {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("internal error - cannot find app client id")
	}

	region := strings.Split(req.MethodArn, ":")[3]

	keyset, err := jwk.Fetch(ctx, "https://cognito-idp."+region+".amazonaws.com/"+userPoolID+"/.well-known/jwks.json")
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("fetch error - cannot find keyset")
	}

	token := []byte(req.QueryStringParameters["auth"])

	parsedToken, err := jwt.Parse(
		token,
		jwt.WithKeySet(keyset),
		jwt.WithValidate(true),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithClaimValue("client_id", appClientID),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), nil
	}

	// fmt.Println(parsedToken)

	return createPolicy(
		req.MethodArn,
		"Allow",
		parsedToken.Subject(),
		map[string]interface{}{
			"username": parsedToken.PrivateClaims()["username"].(string),
		},
	), nil
}

func main() {
	lambda.Start(handler)
}

func getErrorMsg(e error) string {
	clause := " not satisfied"
	switch e.Error() {
	case "exp" + clause:
		return "Token expired"
	case "iss" + clause:
		return "Wrong issuer"
	case "client_id" + clause:
		return "Wrong app client ID"
	case "token_use" + clause:
		return "Wrong token type"
	default:
		return e.Error()
	}
}

func createPolicy(arn, effect, pID string, context map[string]interface{}) (p events.APIGatewayCustomAuthorizerResponse) {
	p.PrincipalID = pID
	p.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Effect:   effect,
				Action:   []string{"execute-api:Invoke"},
				Resource: []string{arn},
			},
		},
	}
	p.Context = context
	p.UsageIdentifierKey = ""

	return
}
