package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

var (
	sess = session.Must(session.NewSession())
)

func handler(req events.DynamoDBEvent) {

	for _, rec := range req.Records {
		// || rec.EventName == "MODIFY"
		if rec.EventName == "INSERT" {

			// for k, v := range rec.Change.NewImage {
			item := rec.Change.NewImage

			fmt.Printf("%s: %+v\n", "new db item", item)
			// if k == "pk" {
			if strings.HasPrefix(item["pk"].String(), "GAME") {
				apiid, ok := os.LookupEnv("CT_APIID")
				if !ok {
					panic(fmt.Sprintf("%v", "can't find api id"))
				}

				stage, ok := os.LookupEnv("CT_STAGE")
				if !ok {
					panic(fmt.Sprintf("%v", "can't find stage"))
				}

				svc := apigatewaymanagementapi.New(sess, &aws.Config{
					Region:   aws.String(rec.AWSRegion),
					Endpoint: aws.String(apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage + "/@connections/"),
				})
				fmt.Println("game")

				var conn apigatewaymanagementapi.PostToConnectionInput

				conn.SetConnectionId(item["sk"].String())

				b, err := json.Marshal("{a: 19894, b: 74156}")
				if err != nil {
					fmt.Println("error marshalling", err)
				}

				conn.SetData(b)

				er := conn.Validate()
				if er != nil {
					fmt.Println("val err", er)
				}

				// fmt.Println("defff", defaults.Get())
				// fmt.Println("defff2222", defaults.Config())

				o, e := svc.PostToConnection(&conn)
				fmt.Println("ooo", o.GoString())
				if e != nil {
					fmt.Println("errrr", e)
				}
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
