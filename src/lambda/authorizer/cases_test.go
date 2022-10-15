package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type input struct {
	ctx context.Context
	inp events.APIGatewayCustomAuthorizerRequestTypeRequest
}

type output struct {
	out events.APIGatewayCustomAuthorizerResponse
	err error
}

var bunchOfTests = []struct {
	input       input
	expected    output
	description string
}{}
