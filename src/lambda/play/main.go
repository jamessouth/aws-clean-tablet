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
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"

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

type key struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

type player struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
	Color  string `dynamodbav:"color,omitempty"`
	Score  int    `dynamodbav:"score"`
	Answer answer `dynamodbav:"answer"`
}

type answer struct {
	PlayerID, Answer string
}

type hiScore struct {
	Score int  `json:"score"`
	Tie   bool `json:"tie"`
}

type livePlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Color  string `dynamodbav:"color"`
	Score  int    `dynamodbav:"score"`
	Answer answer `dynamodbav:"answer"`
}

type livePlayerMap map[string]livePlayer

type liveGame struct {
	Pk           string        `dynamodbav:"pk"`
	Sk           string        `dynamodbav:"sk"`
	Starting     bool          `dynamodbav:"starting"`
	Loading      bool          `dynamodbav:"loading"`
	Players      livePlayerMap `dynamodbav:"players"`
	AnswersCount int           `dynamodbav:"answersCount"`
	SendToFront  bool          `dynamodbav:"sendToFront"`
}

type game struct {
	Pk           string            `dynamodbav:"pk"`
	Sk           string            `dynamodbav:"sk"`
	Starting     bool              `dynamodbav:"starting"`
	Ready        bool              `dynamodbav:"ready"`
	Token        string            `dynamodbav:"token"`
	Loading      bool              `dynamodbav:"loading"`
	Players      map[string]player `dynamodbav:"players"`
	AnswersCount int               `dynamodbav:"answersCount"`
	Wordlist     []string          `dynamodbav:"wordList"`
	SendToFront  bool              `dynamodbav:"sendToFront"`
}

type body struct {
	Gameno, Tipe, Answer, PlayersCount string
}

type sfnArrInput struct {
	Id     string `json:"id"`
	Color  string `json:"color"`
	Name   string `json:"name"`
	ConnID string `json:"connid"`
}

type sfnInput struct {
	Gameno  string        `json:"gameno"`
	Players []sfnArrInput `json:"players"`
}

func (pm livePlayerMap) getLivePlayersSlice() (res []sfnArrInput) {
	count := 0
	for k, v := range pm {
		res = append(res, sfnArrInput{
			Id:     k,
			Color:  colors[count],
			Name:   v.Name,
			ConnID: v.ConnID,
		})
		count++
	}

	return
}

// type sfev struct {
// 	Region, Endpoint, Word, Token string
// 	Conns                         []int
// 	Index                         int
// }

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

	gameItemKey, err := attributevalue.MarshalMap(key{
		Pk: "LISTGME",
		Sk: body.Gameno,
	})
	if err != nil {
		return callErr(err)
	}

	if body.Tipe == "start" {

		di, err := ddbsvc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key:          gameItemKey,
			TableName:    aws.String(tableName),
			ReturnValues: types.ReturnValueAllOld,
		})
		callErr(err)

		var game liveGame
		err = attributevalue.UnmarshalMap(di.Attributes, &game)
		if err != nil {
			return callErr(err)
		}

		fmt.Printf("%s%+v\n", "livegame ", game)

		const numberOfWords int = 40
		lambdasvc := lamb.NewFromConfig(cfg)

		wordsArg, err := json.Marshal(numberOfWords)
		if err != nil {
			fmt.Println("arg to words lambda marshal err", err)
		}

		lambdaInvInput := lamb.InvokeInput{
			FunctionName: aws.String("ct-words"),
			Payload:      wordsArg,
		}

		lambdaInv, err := lambdasvc.Invoke(ctx, &lambdaInvInput)
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

		lambdaReturn := *lambdaInv
		fmt.Printf("\n%s, %+v\n", "liii", lambdaReturn)

		if lambdaReturn.FunctionError != nil || lambdaReturn.StatusCode != 200 {
			fmt.Println("inv pyld err ", *lambdaReturn.FunctionError)
			var errPayload []string
			err = json.Unmarshal(lambdaReturn.Payload, &errPayload)
			if err != nil {
				return callErr(err)
			}
			fmt.Println("err pyld ", errPayload)
		}

		var lambdaPayload []string
		err = json.Unmarshal(lambdaReturn.Payload, &lambdaPayload)
		if err != nil {
			return callErr(err)
		}

		fmt.Println(lambdaPayload)

		words, err := attributevalue.Marshal(lambdaPayload)
		if err != nil {
			return callErr(err)
		}

		marshalledPlayersMap, err := attributevalue.Marshal(game.Players)
		if err != nil {
			return callErr(err)
		}

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "LIVEGME"},
				"sk": &types.AttributeValueMemberS{Value: game.Sk},
			},
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#P": "players",
				"#W": "wordList",
				"#S": "sendToFront",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":f": &types.AttributeValueMemberBOOL{Value: false},
				":p": marshalledPlayersMap,
				":w": words,
			},
			UpdateExpression: aws.String("SET #S = :f, #P = :p, #W = :w"),
		})

		if err != nil {

			return callErr(err)
		}

		players := game.Players.getLivePlayersSlice()

		sfnInput, err := json.Marshal(sfnInput{
			Gameno:  body.Gameno,
			Players: players,
		})
		if err != nil {
			return callErr(err)
		}

		ssei := sfn.StartSyncExecutionInput{
			StateMachineArn: aws.String(sfnarn),
			Input:           aws.String(string(sfnInput)),
		}

		sse, err := sfnsvc.StartSyncExecution(ctx, &ssei)
		if err != nil {
			return callErr(err)
		}

		sseo := *sse
		fmt.Printf("\n%s, %+v\n", "sse op", sseo)

		if sseo.Status == sfntypes.SyncExecutionStatusFailed || sseo.Status == sfntypes.SyncExecutionStatusTimedOut {

			err := fmt.Errorf("step function %s, execution %s, failed with status %s. error code: %s. cause: %s. ", *sseo.StateMachineArn, *sseo.ExecutionArn, sseo.Status, *sseo.Error, *sseo.Cause)
			return callErr(err)

		}

	} else if body.Tipe == "answer" {

		ans, err := attributevalue.Marshal(answer{
			PlayerID: id,
			Answer:   body.Answer,
		})
		if err != nil {
			return callErr(err)
		}

		ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			// ConditionExpression: aws.String("size (#AN) < :c"),
			ExpressionAttributeNames: map[string]string{
				"#PL": "players",
				"#ID": id,
				"#AN": "answer",
				"#AC": "answersCount",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":a": ans,
				":o": &types.AttributeValueMemberN{Value: "1"},
			},
			UpdateExpression: aws.String("SET #PL.#ID.#AN = :a ADD #AC :o"),
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

		if len(gm.Players) == gm.AnswersCount {

			_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				Key:       gameItemKey,
				TableName: aws.String(tableName),
				// ConditionExpression: aws.String("size (#AN) < :c"),
				ExpressionAttributeNames: map[string]string{
					"#S": "sendToFront",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":t": &types.AttributeValueMemberBOOL{Value: true},
				},
				UpdateExpression: aws.String("SET #S = :t"),
				// ReturnValues:     types.ReturnValueAllNew,
			})

			if err != nil {
				return callErr(err)
			}

			answers := map[string][]string{}
			for i, v := range gm.Players {

				fmt.Printf("%s, %v, %+v", "anssss", i, v)

				answers[v.Answer.Answer] = append(answers[v.Answer.Answer], v.Answer.PlayerID)

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
					"#AC": "answersCount",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":z": &types.AttributeValueMemberN{Value: "0"},
					// ":c": &types.AttributeValueMemberN{Value: body.PlayersCount},
				},
				UpdateExpression: aws.String("SET #AC = :z"),
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

			taskOutput, err := json.Marshal(hiScore)
			if err != nil {
				return callErr(err)
			}

			stsi := sfn.SendTaskSuccessInput{
				Output:    aws.String(string(taskOutput)),
				TaskToken: aws.String(gm2.Token),
			}

			_, err = sfnsvc.SendTaskSuccess(ctx, &stsi)
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
