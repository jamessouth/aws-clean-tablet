package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	ebtypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/smithy-go"
)

// export CGO_ENABLED=0 | go build -o main main.go | zip main.zip main | aws lambda update-function-code --function-name ct-lobby --zip-file fileb://main.zip

const maxPlayersPerGame string = "8"

type listPlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Ready  bool   `dynamodbav:"ready"`
}

func getReturnValue(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}
}

func checkInput(s string) (string, string, error) {
	var (
		maxLength  = 99
		gamenoRE   = regexp.MustCompile(`^\d{19}$|^discon$|^newgame$`)
		aW5mb3JtRE = regexp.MustCompile(`^disconnect$|^join$|^leave$|^ready$|^unready$`)
		body       struct{ Gameno, AW5mb3Jt string }
	)

	if len(s) > maxLength {
		return "", "", errors.New("improper json input - too long")
	}

	if strings.Count(s, "gameno") != 1 || strings.Count(s, "aW5mb3Jt") != 1 {
		return "", "", errors.New("improper json input - duplicate/missing key")
	}

	err := json.Unmarshal([]byte(s), &body)
	if err != nil {
		return "", "", err
	}

	if !gamenoRE.MatchString(body.Gameno) {
		return "", "", errors.New("improper json input - bad gameno")
	}

	if !aW5mb3JtRE.MatchString(body.AW5mb3Jt) {
		return "", "", errors.New("improper json input - bad aW5mb3Jt")
	}

	return body.Gameno, body.AW5mb3Jt, nil
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	bod := req.Body

	fmt.Println("lobby", bod, len(bod))

	checkedGameno, checkedAW5mb3Jt, err := checkInput(bod)
	if err != nil {
		return callErr(err)
	}

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return callErr(err)
	}

	var (
		apiid               = os.Getenv("CT_APIID")
		stage               = os.Getenv("CT_STAGE")
		endpoint            = "https://" + apiid + ".execute-api." + reg + ".amazonaws.com/" + stage
		ddbsvc              = dynamodb.NewFromConfig(cfg)
		auth                = req.RequestContext.Authorizer.(map[string]interface{})
		id, name, tableName = auth["principalId"].(string), auth["username"].(string), auth["tableName"].(string)
		connKey             = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "CONNECT"},
			"sk": &types.AttributeValueMemberS{Value: id},
		}
		ebsvc = eventbridge.NewFromConfig(cfg)
	)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == apigatewaymanagementapi.ServiceID && region == reg {
			ep := aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpoint,
				SigningRegion: reg,
			}

			return ep, nil
		}

		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	apigwcfg, err := config.LoadDefaultConfig(ctx,
		// config.WithLogger(logger),
		config.WithRegion(reg),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return callErr(err)
	}

	var apigwsvc = apigatewaymanagementapi.NewFromConfig(apigwcfg)

	if checkedGameno == "newgame" {
		checkedGameno = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	gameItemKey, err := attributevalue.MarshalMap(struct {
		Pk string `dynamodbav:"pk"`
		Sk string `dynamodbav:"sk"`
	}{
		Pk: "LISTGAME",
		Sk: checkedGameno,
	})
	if err != nil {
		return callErr(err)
	}

	removePlayerInput := dynamodb.UpdateItemInput{
		Key:       gameItemKey,
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "players",
			"#I": id,
			"#T": "timerCxld",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":t": &types.AttributeValueMemberBOOL{Value: true},
		},
		UpdateExpression: aws.String("REMOVE #P.#I SET #T = :t"),
		ReturnValues:     types.ReturnValueAllNew,
	}

	if checkedAW5mb3Jt == "join" {

		player := listPlayer{
			Name:   name,
			ConnID: req.RequestContext.ConnectionID,
			Ready:  false,
		}

		marshalledPlayersMap, err := attributevalue.Marshal(map[string]listPlayer{
			id: player,
		})
		if err != nil {
			return callErr(err)
		}

		marshalledPlayer, err := attributevalue.Marshal(player)
		if err != nil {
			return callErr(err)
		}

		updateConnInput := types.Update{
			Key:                 connKey,
			TableName:           aws.String(tableName),
			ConditionExpression: aws.String("size (#G) = :z"),
			ExpressionAttributeNames: map[string]string{
				"#G": "game",
				"#R": "returning",
				"#E": "endtoken",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: checkedGameno},
				":e": &types.AttributeValueMemberS{Value: ""},
				":z": &types.AttributeValueMemberN{Value: "0"},
				":f": &types.AttributeValueMemberBOOL{Value: false},
			},
			UpdateExpression: aws.String("SET #G = :g, #R = :f, #E = :e"),
		}

		_, err = ddbsvc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_exists(#P) AND size (#P) < :m"),
						ExpressionAttributeNames: map[string]string{
							"#P": "players",
							"#I": id,
							"#T": "timerCxld",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":t": &types.AttributeValueMemberBOOL{Value: true},
							":m": &types.AttributeValueMemberN{Value: maxPlayersPerGame},
							":p": marshalledPlayer,
						},
						UpdateExpression: aws.String("SET #P.#I = :p, #T = :t"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		if err != nil {
			return callErr(err)
		}

		_, err = ddbsvc.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                 gameItemKey,
						TableName:           aws.String(tableName),
						ConditionExpression: aws.String("attribute_not_exists(#P)"),
						ExpressionAttributeNames: map[string]string{
							"#P": "players",
							"#T": "timerCxld",
						},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":p": marshalledPlayersMap,
							":t": &types.AttributeValueMemberBOOL{Value: true},
						},
						UpdateExpression: aws.String("SET #P = :p, #T = :t"),
					},
				},
				{
					Update: &updateConnInput,
				},
			},
		})
		if err != nil {
			return callErr(err)
		}

	} else if checkedAW5mb3Jt == "leave" {

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       connKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#G": "game",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: ""},
			},
			UpdateExpression: aws.String("SET #G = :g"),
		})
		if err != nil {
			return callErr(err)
		}

		ui2, err := ddbsvc.UpdateItem(ctx, &removePlayerInput)
		if err != nil {
			return callErr(err)
		}

		err = getReadyStartGame(ui2.Attributes, gameItemKey, tableName, ctx, ddbsvc, apigwsvc, req.RequestContext.RequestTimeEpoch, ebsvc)
		if err != nil {
			return callErr(err)
		}

	} else if checkedAW5mb3Jt == "ready" {

		ui2, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#P": "players",
				"#I": id,
				"#R": "ready",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":t": &types.AttributeValueMemberBOOL{Value: true},
			},
			UpdateExpression: aws.String("SET #P.#I.#R = :t"),
			ReturnValues:     types.ReturnValueAllNew,
		})
		if err != nil {
			return callErr(err)
		}

		err = getReadyStartGame(ui2.Attributes, gameItemKey, tableName, ctx, ddbsvc, apigwsvc, req.RequestContext.RequestTimeEpoch, ebsvc)
		if err != nil {
			return callErr(err)
		}

	} else if checkedAW5mb3Jt == "unready" {

		_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key:       gameItemKey,
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#P": "players",
				"#I": id,
				"#R": "ready",
				"#T": "timerCxld",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":f": &types.AttributeValueMemberBOOL{Value: false},
				":t": &types.AttributeValueMemberBOOL{Value: true},
			},
			UpdateExpression: aws.String("SET #P.#I.#R = :f, #T = :t"),
		})
		if err != nil {
			return callErr(err)
		}

	} else if checkedAW5mb3Jt == "disconnect" {
		if checkedGameno != "discon" {

			ui2, err := ddbsvc.UpdateItem(ctx, &removePlayerInput)

			if err != nil {
				return callErr(err)
			}

			err = getReadyStartGame(ui2.Attributes, gameItemKey, tableName, ctx, ddbsvc, apigwsvc, req.RequestContext.RequestTimeEpoch, ebsvc)
			if err != nil {
				return callErr(err)
			}

		}

		_, err = ddbsvc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key:       connKey,
			TableName: aws.String(tableName),
		})
		if err != nil {
			return callErr(err)
		}

	} else {
		fmt.Println("other lobby")
	}

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
}

func getTimer(gik map[string]types.AttributeValue, tn string, ctx context.Context, ddbsvc *dynamodb.Client, reqTime int64) (bool, error) {
	gi, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
		Key:                  gik,
		TableName:            aws.String(tn),
		ProjectionExpression: aws.String("timerCxld, timerID"),
	})
	if err != nil {
		return false, err
	}

	fmt.Printf("%s: %+v\n", "gi", gi)

	if len(gi.Item) == 0 {
		return false, nil
	}

	var timerData struct {
		TimerID   int64
		TimerCxld bool
	}

	err = attributevalue.UnmarshalMap(gi.Item, &timerData)
	if err != nil {
		return false, err
	}

	return reqTime == timerData.TimerID && !timerData.TimerCxld, nil
}

func getReadyStartGame(rv, gik map[string]types.AttributeValue, tn string, ctx context.Context, ddbsvc *dynamodb.Client, apigwsvc *apigatewaymanagementapi.Client, reqTime int64, ebsvc *eventbridge.Client) error {
	var (
		minPlayers = 3
		gm         struct {
			Sk      string
			Players map[string]listPlayer
		}
	)
	err := attributevalue.UnmarshalMap(rv, &gm)
	if err != nil {
		return err
	}

	if len(gm.Players) < minPlayers {
		return nil
	}

	readyCount := 0
	for _, v := range gm.Players {
		if v.Ready {
			readyCount++
			if readyCount == len(gm.Players) {
				// time.Sleep(1000 * time.Millisecond)
				_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
					Key:       gik,
					TableName: aws.String(tn),
					ExpressionAttributeNames: map[string]string{
						"#I": "timerID",
						"#T": "timerCxld",
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":r": &types.AttributeValueMemberN{Value: strconv.FormatInt(reqTime, 10)},
						":f": &types.AttributeValueMemberBOOL{Value: false},
					},
					UpdateExpression: aws.String("SET #T = :f, #I = :r"),
				})
				if err != nil {
					return err
				}

				fmt.Println("cxld kick off of ticker", reqTime)
				var count byte = 54 // 6
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()
				done := make(chan bool)
				go func() {
					time.Sleep(6 * time.Second)
					done <- true
				}()
				for {
					select {
					case <-done:
						timeToGo, err := getTimer(gik, tn, ctx, ddbsvc, reqTime)
						if err != nil {
							return err
						}
						if timeToGo {
							//kick off game
							fmt.Println("starting game...", reqTime)

							for _, p := range gm.Players {
								conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(p.ConnID), Data: []byte{123, 34, 99, 110, 116, 100, 111, 119, 110, 34, 58, 34, 115, 116, 97, 114, 116, 34, 125}} //{"cntdown": "start"}

								_, err := apigwsvc.PostToConnection(ctx, &conn)

								if err != nil {
									return err
								}

							}

							po, err := ebsvc.PutEvents(ctx, &eventbridge.PutEventsInput{
								Entries: []ebtypes.PutEventsRequestEntry{
									{
										Detail:     aws.String("{\"gameno\":" + "\"" + gm.Sk + "\"" + "}"),
										DetailType: aws.String("initialize game start"),
										Source:     aws.String("lambda.ct-lobby"),
									},
								},
							})
							if err != nil {
								return err
							}

							putResults := *po

							ev := putResults.Entries[0]

							if putResults.FailedEntryCount > 0 {
								return fmt.Errorf("put event failed with msg %s, error code: %s", *ev.ErrorMessage, *ev.ErrorCode)
							}

						}
						return nil
					case <-ticker.C:

						timeToGo, err := getTimer(gik, tn, ctx, ddbsvc, reqTime)
						if err != nil {
							return err
						}
						if timeToGo {

							count -= 1

							for _, p := range gm.Players {

								conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(p.ConnID), Data: []byte{123, 34, 99, 110, 116, 100, 111, 119, 110, 34, 58, count, 125}} //{"cntdown": 4}

								_, err := apigwsvc.PostToConnection(ctx, &conn)
								if err != nil {
									return err
								}

							}
						} else {
							ticker.Stop()
							for _, p := range gm.Players {

								conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(p.ConnID), Data: []byte{123, 34, 99, 110, 116, 100, 111, 119, 110, 34, 58, 34, 34, 125}} //{"cntdown": ""}

								_, err := apigwsvc.PostToConnection(ctx, &conn)
								if err != nil {
									return err
								}

							}
						}

					}
				}

			}
		} else {
			return nil
		}
	}

	return nil
}

func callErr(err error) (events.APIGatewayProxyResponse, error) {
	var transCxldErr *types.TransactionCanceledException
	if errors.As(err, &transCxldErr) {
		fmt.Printf("put item error777, %v\n",
			transCxldErr.CancellationReasons)
	}

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
