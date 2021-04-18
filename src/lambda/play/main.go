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
	lamb "github.com/aws/aws-sdk-go-v2/service/lambda"
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
type Player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
}

type body struct {
	Game, Type string
	Value      bool
}

// ConnItemAttrs holds vals for db
type ConnItemAttrs struct {
	Game string `dynamodbav:":g"`
	Zero *int   `dynamodbav:":zero,omitempty"`
}

type lambdaInput struct {
	Game   map[string]Player `json:"game"`
	ApiId  string            `json:"apiid"`
	Stage  string            `json:"stage"`
	Region string            `json:"region"`
}

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
	svc2 := lamb.NewFromConfig(cfg)

	// auth := req.RequestContext.Authorizer.(map[string]interface{})

	// id := auth["principalId"].(string)
	// name := auth["username"].(string)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		fmt.Println("unmarshal err")
	}

	if body.Type == "start" {

		gameItemKey, err := attributevalue.MarshalMap(Key{
			Pk: "GAME",
			Sk: body.Game,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record, %v", err))
		}

		att1, err := attributevalue.Marshal(true)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal Record 22, %v", err))
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{

			// ----------------------------------------------------
			Key:                 gameItemKey,
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("#ST <> :s"),
			ExpressionAttributeNames: map[string]string{
				// "#PL": "players",
				// "#ID": id,
				"#ST": "starting",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":s": att1,
				// ":maxsize": att2,
				// ":player": att3,
			},

			UpdateExpression: aws.String("SET #ST = :s"),
			// ReturnValues:     types.ReturnValueAllNew,
		})
		// fmt.Println("op", op)
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

		} else {

			mj, err := json.Marshal(lambdaInput{
				// Game:   game,
				// ApiId:  req.RequestContext.APIID,
				// Stage:  req.RequestContext.Stage,
				Region: reg + "junkjunk",
			})
			if err != nil {
				fmt.Println("game item marshal err", err)
			}

			ii := lamb.InvokeInput{
				FunctionName: aws.String("ct-playJS"),
				// ClientContext:  new(string),
				// InvocationType: "",
				// LogType:        "",
				Payload: mj,
				// Qualifier:      new(string),
			}

			li, err := svc2.Invoke(ctx, &ii)

			fmt.Printf("\n%s, %+v\n", "liii", *li)
			// fmt.Println(*li.FunctionError, li.Payload)
			q := *li
			z := q.FunctionError
			x := q.Payload
			// fmt.Println(*z, x)

			if z != nil {
				fmt.Println(*z, x)
			}

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
		}

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