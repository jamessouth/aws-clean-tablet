package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

var (
	sess = session.Must(session.NewSession())
	svc  = apigatewaymanagementapi.New(sess)
)

func handler(req events.DynamoDBEvent) {
	apiid, ok := os.LookupEnv("CT_APIID")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find api id"))
	}

	stage, ok := os.LookupEnv("CT_STAGE")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find stage"))
	}

	for _, rec := range req.Records {
		// || rec.EventName == "MODIFY"
		if rec.EventName == "INSERT" {

			// for k, v := range rec.Change.NewImage {
			item := rec.Change.NewImage

			fmt.Printf("%s: %+v\n", "new db item", item)
			// if k == "pk" {
			if strings.HasPrefix(item["pk"].String(), "GAME") {

				fmt.Println("game")

				var conn *apigatewaymanagementapi.PostToConnectionInput

				conn.SetConnectionId("https://" + apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage + "/@connections/" + item["sk"].String())

				conn.SetData([]byte("yoyoyoyo"))

				svc.PostToConnection(conn)
			} else {
				fmt.Println("other")

			}
			// }
			// }
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
