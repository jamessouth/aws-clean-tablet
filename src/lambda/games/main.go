package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/smithy-go"
)

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {

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
				str := "https://" + apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage

				fmt.Println(str)

				customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					if service == apigatewaymanagementapi.ServiceID && region == rec.AWSRegion {
						return aws.Endpoint{
							PartitionID:   "aws",
							URL:           str,
							SigningRegion: rec.AWSRegion,
						}, nil
					}
					return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
				})

				cfg, err := config.LoadDefaultConfig(ctx,
					config.WithRegion(rec.AWSRegion),
					config.WithEndpointResolver(customResolver),
				)
				if err != nil {
					fmt.Println("cfg err")
				}

				// , &aws.Config{
				// 	Region:   aws.String(),
				// 	Endpoint: aws.String(apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage + "/@connections/"),
				// }

				svc := apigatewaymanagementapi.NewFromConfig(cfg)

				fmt.Println("game")

				b, err := json.Marshal("{a: 19894, b: 74156}")
				if err != nil {
					fmt.Println("error marshalling", err)
				}
				fmt.Println(b, item["sk"].String())
				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(item["sk"].String()), Data: b}

				// conn.SetConnectionId()

				// conn.SetData(b)

				// er := conn.Validate()
				// if er != nil {
				// 	fmt.Println("val err", er)
				// }

				// fmt.Println("defff", defaults.Get())
				// fmt.Println("defff2222", defaults.Config())

				o, e := svc.PostToConnection(ctx, &conn)
				fmt.Println("opopopo", o)
				if e != nil {
					fmt.Println("errrr", e)

					// To get any API error
					var apiErr smithy.APIError
					if errors.As(err, &apiErr) {
						fmt.Printf("db error, Code: %v, Message: %v",
							apiErr.ErrorCode(), apiErr.ErrorMessage())
					}
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
	fmt.Println("returrrrn")
	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
