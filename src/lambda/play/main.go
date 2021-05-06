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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/smithy-go"
)

// var colors = []string{
// 	"#dc2626", //red 600
// 	"#0c4a6e", //light blue 900
// 	"#16a34a", //green 600
// 	"#7c2d12", //orange 900
// 	"#c026d3", //fuchsia 600
// 	"#365314", //lime 900
// 	"#0891b2", //cyan 600
// 	"#581c87", //purple 900
// }

// const maxPlayersPerGame = 8

// Key holds values to be put in db
type Key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// Player holds values to be put in db
// type player struct {
// 	Name   string `json:"name"`
// 	ConnID string `json:"connid"`
// 	Ready  bool   `json:"ready"`
// 	Color  string `json:"color"`
// }

// type game struct {
// 	Pk       string   `json:"pk"`
// 	Sk       string   `json:"sk"`
// 	Starting bool     `json:"starting"`
// 	Leader   string   `json:"leader"`
// 	Loading  bool     `json:"loading"`
// 	Players  []player `json:"players"`
// }

type body struct {
	Gameno string
	Type   string
}

// ConnItemAttrs holds vals for db
type ConnItemAttrs struct {
	Game string `dynamodbav:":g"`
	Zero *int   `dynamodbav:":zero,omitempty"`
}

// type lambdaInput struct {
// 	Game   game   `json:"game"`
// 	Region string `json:"region"`
// }

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("plaaaaaaay", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	// .WithEndpoint("http://192.168.4.27:8000")

	svc := dynamodb.NewFromConfig(cfg)
	// svc2 := lamb.NewFromConfig(cfg)

	// auth := req.RequestContext.Authorizer.(map[string]interface{})

	// id := auth["principalId"].(string)
	// name := auth["username"].(string)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal body err")
	}

	if body.Type == "start" {

		gameItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: body.Gameno,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal gik Record, %v", err))
		}

		gi, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
		})
		// j := &gi

		fmt.Printf("%s%+v\n", "get item", gi)
		fmt.Printf("%s%+v\n", "get item22222", *gi)
		// fmt.Println("get item2222", j.Item)

		if err != nil {

			var intServErr *types.InternalServerError
			if errors.As(err, &intServErr) {
				fmt.Printf("get item error, %v",
					intServErr.ErrorMessage())
			}

			// To get any API error
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				fmt.Printf("play rt db error, Code: %v, Message: %v",
					apiErr.ErrorCode(), apiErr.ErrorMessage())
			}
		}
		// } else {

		// var game game
		// err = attributevalue.UnmarshalMap(ui.Attributes, &game)
		// if err != nil {
		// 	fmt.Println("play update item unmarshal err", err)
		// }

		// mj, err := json.Marshal(lambdaInput{
		// 	Game: gi.Item,
		// 	// ApiId:  req.RequestContext.APIID,
		// 	// Stage:  req.RequestContext.Stage,
		// 	Region: reg,
		// })
		// if err != nil {
		// 	fmt.Println("game item marshal err", err)
		// }

		// ii := lamb.InvokeInput{
		// 	FunctionName: aws.String("ct-playJS"),
		// 	// ClientContext:  new(string),
		// 	// InvocationType: "",
		// 	// LogType:        "",
		// 	Payload: mj,
		// 	// Qualifier:      new(string),
		// }

		// li, err := svc2.Invoke(ctx, &ii)

		// q := *li
		// fmt.Printf("\n%s, %+v\n", "liii", q)
		// // fmt.Println(*li.FunctionError, li.Payload)
		// z := q.FunctionError
		// x := string(q.Payload)
		// fmt.Println("inv pyld", x)

		// if z != nil {
		// 	fmt.Println("inv err", *z, x)
		// }

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
		// }

	}

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
