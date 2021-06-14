package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// type input struct {
// 	Iterator data `json:"iterator"`
// }

type data2 struct {
	// Index int `json:"index"`
	Word string `json:"word"`
}
type data struct {
	// Index int `json:"index"`
	Word   data2
	Gameno string `json:"gameno"`
	C      []string
	S      int
}

type output struct {
	Gameno string `json:"gameno"`
	// Index    int  `json:"index"`
	S int `json:"s"`
	C []string
}

func handler(ctx context.Context, req data) (output, error) {
	s := (req.S + 1) * 2

	fmt.Printf("%s%+v\n", "req: ", req)

	return output{
		Gameno: req.Gameno,
		S:      s,
		C:      req.C,
	}, nil
}

func main() {
	lambda.Start(handler)
}
