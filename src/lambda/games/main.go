package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.DynamoDBEvent) {

	for _, rec := range req.Records {

		if rec.EventName == "INSERT" || rec.EventName == "MODIFY" {

			for k, v := range rec.Change.NewImage {

				fmt.Printf("%s - %v: %v - %s\n", "k v", k, v, rec.EventName)
				if k == "pk" {
					if strings.HasPrefix(v.String(), "GAME") {
						fmt.Println("game")
					} else {
						fmt.Println("other")

					}
				}
			}
		} else {
			for k, v := range rec.Change.Keys {

				fmt.Printf("%s - %v: %v - %s\n", "k v", k, v, rec.EventName)

			}

		}

	}

}

func main() {
	lambda.Start(handler)
}
