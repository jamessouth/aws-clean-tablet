package main

import (
	// "bufio"

	"context"
	"encoding/json"
	"errors"
	"fmt"

	// "io"
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

// var buffer bytes.Buffer

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	fmt.Println("reqqqq", req)
	for _, rec := range req.Records {
		// || rec.EventName == "MODIFY"
		if rec.EventName == "INSERT" {

			// for k, v := range rec.Change.NewImage {
			item := rec.Change.NewImage

			fmt.Printf("%s: %+v\n", "new db item", item)
			fmt.Println("nnnn", item["pk"].String())
			// if k == "pk" {
			if strings.HasPrefix(item["pk"].String(), "CONN") {
				apiid, ok := os.LookupEnv("CT_APIID")
				if !ok {
					panic(fmt.Sprintf("%v", "can't find api id"))
				}

				stage, ok := os.LookupEnv("CT_STAGE")
				if !ok {
					panic(fmt.Sprintf("%v", "can't find stage"))
				}
				str := "https://" + apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage

				// fmt.Println(str)

				customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					if service == apigatewaymanagementapi.ServiceID && region == rec.AWSRegion {
						ep := aws.Endpoint{
							PartitionID:   "aws",
							URL:           str,
							SigningRegion: rec.AWSRegion,
						}
						fmt.Println("eppppppppppppppppp", ep)
						return ep, nil
					}
					return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
				})

				// logger := logging.NewStandardLogger(&buffer)
				// logger.Logf(logging.Debug, "time to %s", "log")

				cfg, err := config.LoadDefaultConfig(ctx,
					config.WithRegion(rec.AWSRegion),
					// config.WithLogger(logger),
					config.WithEndpointResolver(customResolver),
				)
				if err != nil {
					fmt.Println("cfg err")
				}
				// fmt.Println("cfggggggggggggg", cfg)
				// , &aws.Config{
				// 	Region:   aws.String(),
				// 	Endpoint: aws.String(apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage + "/@connections/"),
				// }

				svc := apigatewaymanagementapi.NewFromConfig(cfg)

				// , func(o *apigatewaymanagementapi.Options) {
				// 	o.ClientLogMode = aws.LogSigning | aws.LogRequest | aws.LogResponseWithBody
				// }
				// fmt.Println("game")

				b, err := json.Marshal("{a: 19894, b: 74156}")
				if err != nil {
					fmt.Println("error marshalling", err)
				}
				fmt.Println("xxxxxxxxxxxxx", item["GSI1SK"].String(), string(b))
				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(item["GSI1SK"].String()), Data: b}

				// conn.SetConnectionId()

				// conn.SetData(b)

				// er := conn.Validate()
				// if er != nil {
				// 	fmt.Println("val err", er)
				// }

				// fmt.Println("defff", defaults.Get())
				// fmt.Println("defff2222", defaults.Config())

				_, e := svc.PostToConnection(ctx, &conn)
				// fmt.Println("opopopo", o)
				if e != nil {
					fmt.Println("errrr", e)

					// To get any API error
					var apiErr smithy.APIError
					if errors.As(err, &apiErr) {
						fmt.Printf("db error, Code: %v, Message: %v",
							apiErr.ErrorCode(), apiErr.ErrorMessage())
					}
				}
			}
			// }
			// }
		}

	}
	// r := bufio.NewReader(&buffer)
	// log := make([]byte, 0, 1024)
	// for {
	// 	n, err := io.ReadFull(r, log[:cap(log)])
	// 	log = log[:n]
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != io.ErrUnexpectedEOF {
	// 			fmt.Fprintln(os.Stderr, err)
	// 			break
	// 		}
	// 	}

	// 	fmt.Printf("read %d bytes: ", n)

	// 	fmt.Println(string(log))
	// }

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
