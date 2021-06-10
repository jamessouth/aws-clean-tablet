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
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/aws/smithy-go"
)

var colors = []string{
	"#dc2626", //red 600
	"#0c4a6e", //light blue 900
	"#16a34a", //green 600
	"#7c2d12", //orange 900
	"#c026d3", //fuchsia 600
	"#365314", //lime 900
	"#0891b2", //cyan 600
	"#581c87", //purple 900
}

// const maxPlayersPerGame = 8

// Key holds values to be put in db
type Key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// Player holds values to be put in db
type player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
	Color  string `dynamodbav:"color,omitempty"`
	Score  int    `dynamodbav:"score"`
}

type answer struct {
	PlayerID, Answer string
}

type hiScore struct {
	Score int  `json:"score"`
	Tie   bool `json:"tie"`
}

type game struct {
	Pk       string            `dynamodbav:"pk"`
	Sk       string            `dynamodbav:"sk"`
	Starting bool              `dynamodbav:"starting"`
	Leader   string            `dynamodbav:"leader"`
	Loading  bool              `dynamodbav:"loading"`
	Players  map[string]player `dynamodbav:"players"`
	Answers  []answer          `dynamodbav:"answers"`
}

type body struct {
	Gameno, Type, Answer, PlayersCount string
}

// ConnItemAttrs holds vals for db
type ConnItemAttrs struct {
	Game string `dynamodbav:":g"`
	Zero *int   `dynamodbav:":zero,omitempty"`
}

type lambdaInput struct {
	Game    game    `json:"game,omitempty"`
	Region  string  `json:"region"`
	HiScore hiScore `json:"hiScore,omitempty"`
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("plaaaaaaay", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return callErr(err)
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	sfnarn, ok := os.LookupEnv("SFNARN")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find sfn arn"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)
	sfnsvc := sfn.NewFromConfig(cfg)

	id := req.RequestContext.Authorizer.(map[string]interface{})["principalId"].(string)

	var body body

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	gameItemKey, err := attributevalue.MarshalMap(Key{
		Pk: "GAME",
		Sk: body.Gameno,
	})
	if err != nil {
		return callErr(err)
	}

	if body.Type == "start" {

		gi, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
		})

		if err != nil {

			return callErr(err)
		}

		var game game
		err = attributevalue.UnmarshalMap(gi.Item, &game)
		if err != nil {
			return callErr(err)
		}

		fmt.Printf("%s%+v\n", "gammmmme ", game)

		mj, err := json.Marshal(sfn.StartSyncExecutionInput{
			StateMachineArn: aws.String(sfnarn),
			Input:           new(string),
		})
		if err != nil {
			return callErr(err)
		}

		ii := lamb.InvokeInput{
			FunctionName: aws.String("ct-playJS"),
			Payload:      mj,
		}

		li, err := svc2.Invoke(ctx, &ii)

		q := *li
		fmt.Printf("\n%s, %+v\n", "liii", q)
		// fmt.Println(*li.FunctionError, li.Payload)
		z := q.FunctionError
		x := string(q.Payload)
		fmt.Println("inv pyld", x)

		if z != nil {
			fmt.Println("inv err", *z, x)
		}

		if err != nil {

			return callErr(err)

		}

	} else if body.Type == "answer" {

		ans, err := attributevalue.MarshalList(answer{
			PlayerID: id,
			Answer:   body.Answer,
		})
		if err != nil {
			return callErr(err)
		}

		ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:                 gameItemKey,
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("size (#AN) < :c"),
			ExpressionAttributeNames: map[string]string{
				// "#PL": "players",
				// "#ID": id,
				"#AN": "answers",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":a": &types.AttributeValueMemberL{Value: ans},
				":c": &types.AttributeValueMemberN{Value: body.PlayersCount},
			},
			UpdateExpression: aws.String("SET #AN = list_append(#AN, :a)"),
			ReturnValues:     types.ReturnValueAllNew,
		})

		if err != nil {

			return callErr(err)

		}

		var gm game
		err = attributevalue.UnmarshalMap(ui.Attributes, &gm)
		if err != nil {
			return callErr(err)
		}

		if len(gm.Players) == len(gm.Answers) {

			answers := map[string][]string{}
			for i, v := range gm.Answers {

				fmt.Printf("%s, %v, %+v", "anssss", i, v)

				answers[v.Answer] = append(answers[v.Answer], v.PlayerID)

			}

			for k, v := range answers {

				fmt.Printf("%s, %v, %+v", "anssssmapppp", k, v)

				switch {
				case len(k) < 2:
					// c.updateEachScore(v, 0)

				case len(v) > 2:
					// c.updateEachScore(v, 1)

					for _, id := range v {

						_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
							Key:       gameItemKey,
							TableName: aws.String(tableName),
							ExpressionAttributeNames: map[string]string{
								"#PL": "players",
								"#ID": id,
								"#SC": "score",
							},
							ExpressionAttributeValues: map[string]types.AttributeValue{
								":s": &types.AttributeValueMemberN{Value: "1"},
							},
							UpdateExpression: aws.String("ADD #PL.#ID.#SC :s"),
						})

						if err != nil {

							return callErr(err)

						}
					}

				case len(v) == 2:
					// c.updateEachScore(v, 3)

					for _, id := range v {

						_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
							Key:       gameItemKey,
							TableName: aws.String(tableName),
							ExpressionAttributeNames: map[string]string{
								"#PL": "players",
								"#ID": id,
								"#SC": "score",
							},
							ExpressionAttributeValues: map[string]types.AttributeValue{
								":s": &types.AttributeValueMemberN{Value: "3"},
							},
							UpdateExpression: aws.String("ADD #PL.#ID.#SC :s"),
						})

						if err != nil {

							return callErr(err)

						}
					}

				default:
					// c.updateEachScore(v, 0)
				}

			}

			ui2, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key:       gameItemKey,
				TableName: aws.String(tableName),
				// ConditionExpression: aws.String("size (#AN) < :c"),
				ExpressionAttributeNames: map[string]string{
					// "#PL": "players",
					// "#ID": id,
					"#AN": "answers",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":a": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
					// ":c": &types.AttributeValueMemberN{Value: body.PlayersCount},
				},
				UpdateExpression: aws.String("SET #AN = :a"),
				ReturnValues:     types.ReturnValueAllNew,
			})

			if err != nil {

				return callErr(err)

			}

			var gm2 game
			err = attributevalue.UnmarshalMap(ui2.Attributes, &gm2)
			if err != nil {
				return callErr(err)
			}

			hiScore := hiScore{
				Score: 0,
				Tie:   false,
			}
			// && hiScore.Score > 0
			for _, v := range gm2.Players {
				if v.Score == hiScore.Score {
					hiScore.Tie = true
				}
				if v.Score > hiScore.Score {
					hiScore.Score = v.Score
					hiScore.Tie = false
				}
			}

			mj2, err := json.Marshal(lambdaInput{
				HiScore: hiScore,
				Region:  reg,
			})
			if err != nil {
				return callErr(err)
			}

			ii2 := lamb.InvokeInput{
				FunctionName: aws.String("ct-wrdJS"),
				Payload:      mj2,
			}

			li2, err := svc2.Invoke(ctx, &ii2)

			q := *li2
			fmt.Printf("\n%s, %+v\n", "liii", q)
			// fmt.Println(*li.FunctionError, li.Payload)
			z := q.FunctionError
			x := string(q.Payload)
			fmt.Println("inv pyld", x)

			if z != nil {
				fmt.Println("inv err", *z, x)
			}

			if err != nil {

				return callErr(err)

			}

		}

	} else {
		fmt.Println("other play")
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

func callErr(err error) (events.APIGatewayProxyResponse, error) {

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

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusBadRequest,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, err

}
