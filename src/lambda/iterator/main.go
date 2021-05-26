package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type input struct {
	Comment  string `json:"comment"`
	Iterator data   `json:"iterator"`
}

type data struct {
	Index int `json:"index"`
	Step  int `json:"step"`
	Count int `json:"count"`
}

type output struct {
	Index    int  `json:"index"`
	Step     int  `json:"step"`
	Count    int  `json:"count"`
	Continue bool `json:"continue"`
}

func handler(ctx context.Context, req input) (output, error) {
	index := req.Iterator.Index + req.Iterator.Step

	return output{
		Index: index,
		Step:  req.Iterator.Step,
		Count: req.Iterator.Count,
		// Continue: index < req.Iterator.Count,
	}, nil
}

func main() {
	lambda.Start(handler)
}
