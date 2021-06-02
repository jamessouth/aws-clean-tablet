package main

import (
	"context"
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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

type player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
	Color  string `dynamodbav:"color"`
}

type game struct {
	No      string   `json:"no"`
	Leader  string   `json:"leader,omitempty"`
	Players []player `json:"players"`
}

type payload struct {
	Games []game `json:"games"`
	Type  string `json:"type"`
}

// var (
// 	connResults, gamesResults *dynamodb.QueryOutput
// )

func getPlayersSlice(m map[string]player) (res []player) {
	for _, v := range m {
		res = append(res, v)
	}

	return
}

type games struct {
	*dynamodb.QueryOutput
}

func (m games) mapGames() []game {
	fmt.Print(m)
	return []game{}
}

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	// fmt.Println("reqqqq", req)
	for _, rec := range req.Records {
		tableName := strings.Split(rec.EventSourceArn, "/")[1]
		item := rec.Change.NewImage
		fmt.Printf("%s: %+v\n", "new db item", item)

		apiid, ok := os.LookupEnv("CT_APIID")
		if !ok {
			panic(fmt.Sprintf("%v", "can't find api id"))
		}

		stage, ok := os.LookupEnv("CT_STAGE")
		if !ok {
			panic(fmt.Sprintf("%v", "can't find stage"))
		}
		endpoint := "https://" + apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage

		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == apigatewaymanagementapi.ServiceID && region == rec.AWSRegion {
				ep := aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: rec.AWSRegion,
				}
				fmt.Println("eppppppppppppppppp", ep)
				return ep, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
			config.WithEndpointResolver(customResolver),
		)
		if err != nil {
			fmt.Println("cfg err")
		}

		apigwsvc := apigatewaymanagementapi.NewFromConfig(cfg)
		ddbsvc := dynamodb.NewFromConfig(cfg)

		if rec.EventName == "INSERT" {
			if strings.HasPrefix(item["pk"].String(), "CONN") {

				gamesParams := dynamodb.QueryInput{
					TableName:              aws.String(tableName),
					ScanIndexForward:       aws.Bool(false),
					KeyConditionExpression: aws.String("pk = :gm"),
					FilterExpression:       aws.String("#ST = :st"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":gm": &types.AttributeValueMemberS{Value: "GAME"},
						":st": &types.AttributeValueMemberBOOL{Value: false},
					},
					ExpressionAttributeNames: map[string]string{
						"#ST": "starting",
					},
				}
				// if err != nil {
				// 	panic(fmt.Sprintf("failed to marshal query input, %v", err))
				// }
				gamesResults, err := ddbsvc.Query(ctx, &gamesParams)
				if err != nil {

					var intServErr *types.InternalServerError
					if errors.As(err, &intServErr) {
						fmt.Printf("get item error, %v",
							intServErr.ErrorMessage())
					}

					// To get any API error
					var apiErr smithy.APIError
					if errors.As(err, &apiErr) {
						fmt.Printf("db error, Code: %v, Message: %v",
							apiErr.ErrorCode(), apiErr.ErrorMessage())
					}

				}

				payload := payload{
					Games: gamesResults.Items.mapGames(),
					Type:  "",
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
