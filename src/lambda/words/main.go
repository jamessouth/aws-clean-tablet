package main

import (
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

const amountOfWords int = 50

func handler(length int) ([]string, error) {
	if length < 1 || length > len(words) {
		length = amountOfWords
	}
	return shuffleList(words, length), nil // error required by lambda
}

func main() {
	lambda.Start(handler)
}

func shuffleList(words []string, length int) []string {
	t := time.Now().UnixNano()
	rand.Seed(t)

	nl := append([]string(nil), words...)

	rand.Shuffle(len(nl), func(i, j int) {
		nl[i], nl[j] = nl[j], nl[i]
	})

	return nl[:length]
}
