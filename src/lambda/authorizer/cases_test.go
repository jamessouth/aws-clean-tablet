package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

type input struct {
	ctx context.Context
	inp events.APIGatewayCustomAuthorizerRequestTypeRequest
}

var bunchOfTests = []struct {
	input        input
	expected_rv  events.APIGatewayCustomAuthorizerResponse
	expected_err error
	description  string
}{
	{
		input: input{
			ctx: context.TODO(),
			inp: events.APIGatewayCustomAuthorizerRequestTypeRequest{
				Type:                            "",
				MethodArn:                       "arn:aws:execute-api:catan:999:888/777/666",
				Resource:                        "",
				Path:                            "",
				HTTPMethod:                      "",
				Headers:                         map[string]string{"Origin": "http://localhost:4200"},
				MultiValueHeaders:               map[string][]string{},
				QueryStringParameters:           map[string]string{"auth": ""},
				MultiValueQueryStringParameters: map[string][]string{},
				PathParameters:                  map[string]string{},
				StageVariables:                  map[string]string{},
				RequestContext:                  events.APIGatewayCustomAuthorizerRequestTypeRequestContext{},
			},
		},
		expected_rv:  events.APIGatewayCustomAuthorizerResponse{},
		expected_err: errors.New("header error - request from wrong domain"),
		description:  "wrong origin",
	},
	{
		input: input{
			ctx: context.TODO(),
			inp: events.APIGatewayCustomAuthorizerRequestTypeRequest{
				Type:                            "",
				MethodArn:                       "arn:aws:execute-api:catan:999:888/777/666",
				Resource:                        "",
				Path:                            "",
				HTTPMethod:                      "",
				Headers:                         map[string]string{"User-Agent": "lkmlkmlkm"},
				MultiValueHeaders:               map[string][]string{},
				QueryStringParameters:           map[string]string{"auth": ""},
				MultiValueQueryStringParameters: map[string][]string{},
				PathParameters:                  map[string]string{},
				StageVariables:                  map[string]string{},
				RequestContext:                  events.APIGatewayCustomAuthorizerRequestTypeRequestContext{},
			},
		},
		expected_rv:  events.APIGatewayCustomAuthorizerResponse{},
		expected_err: errors.New("header error - request from wrong client"),
		description:  "wrong client",
	},
}
