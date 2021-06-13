package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type input struct {
	Iterator data `json:"iterator"`
}

type data struct {
	Index int `json:"index"`
	Count int `json:"count"`
}

type output struct {
	Index    int  `json:"index"`
	Count    int  `json:"count"`
	Continue bool `json:"continue"`
}

func handler(ctx context.Context, req input) (output, error) {
	index := req.Iterator.Index + 1

	return output{
		Index:    index,
		Count:    req.Iterator.Count,
		Continue: index < req.Iterator.Count,
	}, nil
}

func main() {
	lambda.Start(handler)
}
