package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// type input struct {
// 	Iterator data `json:"iterator"`
// }

type data struct {
	// Index int `json:"index"`
	W, C []string
	S    int
}

// type output struct {
// 	// Index    int  `json:"index"`
// 	S int `json:"s"`
// }

func handler(ctx context.Context, req data) (data, error) {
	s := (req.S + 1) * 2

	return data{
		W: req.W[1:],
		C: req.C,
		S: s,
	}, nil
}

func main() {
	lambda.Start(handler)
}
