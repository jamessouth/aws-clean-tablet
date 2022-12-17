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

const (
	connect           string = "CONNECT"
	listGame          string = "LISTGAME"
	leave             string = "leave"
	maxPlayersPerGame string = "8"
	newgame           string = "newgame"
	join              string = "join"
)

type listPlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
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
		maxLength = 99
		gamenoRE  = regexp.MustCompile(`^\d{19}$|^newgame$`)
		commandRE = regexp.MustCompile(`^join$|^leave$`)
		body      struct{ Gameno, Command string }
	)

	if len(s) > maxLength {
		return "", "", fmt.Errorf("improper json input - too long: %d", len(s))
	}

	if strings.Count(s, "gameno") != 1 || strings.Count(s, "command") != 1 {
		return "", "", errors.New("improper json input - duplicate/missing key")
	}

	err := json.Unmarshal([]byte(s), &body)
	if err != nil {
		return "", "", err
	}

	var gameno, command = body.Gameno, body.Command

	switch {
	case !gamenoRE.MatchString(gameno):
		return "", "", errors.New("improper json input - bad gameno: " + gameno)
	case !commandRE.MatchString(command):
		return "", "", errors.New("improper json input - bad command: " + command)
	case command == leave && gameno == newgame:
		return "", "", errors.New("improper json input - leave/newgame mismatch")
	}

	return gameno, command, nil
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		bod    = req.Body
		region = strings.Split(req.RequestContext.DomainName, ".")[2]
	)

	fmt.Println("lobby", bod, len(bod))

	checkedGameno, checkedCommand, err := checkInput(bod)
	if err != nil {
		return callErr(err)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return callErr(err)
	}

	var (
		ddbsvc              = dynamodb.NewFromConfig(cfg)
		auth                = req.RequestContext.Authorizer.(map[string]interface{})
		id, name, tableName = auth["principalId"].(string), auth["username"].(string), auth["tableName"].(string)
		player              = listPlayer{
			Name:   name,
			ConnID: req.RequestContext.ConnectionID,
		}
	)

	if checkedGameno == newgame {
		marshalledPlayersMap, err := attributevalue.Marshal(map[string]listPlayer{
			id: player,
		})
		if err != nil {
			return callErr(err)
		}

		_, err = ddbsvc.PutItem(ctx, &dynamodb.PutItemInput{
			Item: map[string]types.AttributeValue{
				"pk":        &types.AttributeValueMemberS{Value: listGame},
				"sk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", time.Now().UnixNano())},
				"players":   marshalledPlayersMap,
				"timerCxld": &types.AttributeValueMemberBOOL{Value: true},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			return callErr(err)
		}
	} else {
		gameItemKey, err := attributevalue.MarshalMap(struct {
			Pk string `dynamodbav:"pk"`
			Sk string `dynamodbav:"sk"`
		}{
			Pk: listGame,
			Sk: checkedGameno,
		})
		if err != nil {
			return callErr(err)
		}

		ean := map[string]string{
			"#P": "players",
			"#I": id,
			"#T": "timerCxld",
		}

		updateParams := dynamodb.UpdateItemInput{
			Key:                      gameItemKey,
			TableName:                aws.String(tableName),
			ExpressionAttributeNames: ean,
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":t": &types.AttributeValueMemberBOOL{Value: true},
			},
			ReturnValues: types.ReturnValueUpdatedNew,
		}

		if checkedCommand == join {
			marshalledPlayer, err := attributevalue.Marshal(player)
			if err != nil {
				return callErr(err)
			}

			updateParams.ConditionExpression = aws.String("attribute_exists(sk) AND size (#P) < :m")
			updateParams.ExpressionAttributeValues[":m"] = &types.AttributeValueMemberN{Value: maxPlayersPerGame}
			updateParams.ExpressionAttributeValues[":p"] = marshalledPlayer
			updateParams.UpdateExpression = aws.String("SET #P.#I = :p, #T = :t")

		} else if checkedCommand == leave {
			updateParams.ConditionExpression = aws.String("attribute_exists(sk)")
			updateParams.UpdateExpression = aws.String("REMOVE #P.#I SET #T = :t")

		}

		ui, err := ddbsvc.UpdateItem(ctx, &updateParams)
		if err != nil {
			return callErr(err)
		}

		var (
			gm struct {
				Sk      string
				Players map[string]listPlayer
			}
			apiid    = os.Getenv("CT_APIID")
			stage    = os.Getenv("CT_STAGE")
			endpoint = "https://" + apiid + ".execute-api." + region + ".amazonaws.com/" + stage
			ebsvc    = eventbridge.NewFromConfig(cfg)
		)

		err = attributevalue.UnmarshalMap(ui.Attributes, &gm)
		if err != nil {
			return callErr(err)
		}

		apigwcfg, err := config.LoadDefaultConfig(ctx,
			// config.WithLogger(logger),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, awsRegion string, options ...interface{}) (aws.Endpoint, error) {
				if service == apigatewaymanagementapi.ServiceID && awsRegion == region {
					ep := aws.Endpoint{
						PartitionID:   "aws",
						URL:           endpoint,
						SigningRegion: awsRegion,
					}

					return ep, nil
				}

				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			})),
		)
		if err != nil {
			return callErr(err)
		}

		var apigwsvc = apigatewaymanagementapi.NewFromConfig(apigwcfg)

		err = getStartGame(ctx, gameItemKey, gm.Players, gm.Sk, tableName, req.RequestContext.RequestTimeEpoch, ddbsvc, apigwsvc, ebsvc)
		if err != nil {
			return callErr(err)
		}
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

func sendCount(ctx context.Context, players map[string]listPlayer, apigwsvc *apigatewaymanagementapi.Client, data []byte) error {
	for _, p := range players {
		conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(p.ConnID), Data: append([]byte{123, 34, 99, 110, 116, 100, 111, 119, 110, 34, 58}, data...)} //{"cntdown":

		_, err := apigwsvc.PostToConnection(ctx, &conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func getStartGame(ctx context.Context, gik map[string]types.AttributeValue, players map[string]listPlayer, sk, tn string, reqTime int64, ddbsvc *dynamodb.Client, apigwsvc *apigatewaymanagementapi.Client, ebsvc *eventbridge.Client) error {
	if len(players) < 3 { //minPlayers
		return nil
	}

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
	var count byte = 57 // 9
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	go func() {
		time.Sleep(9 * time.Second)
		done <- true
	}()

	for {
		select {
		case <-done:
			timeToGo, err := getTimer(gik, tn, ctx, ddbsvc, reqTime)
			if err != nil {
				return err
			}

			if timeToGo { //kick off game
				fmt.Println("starting game...", reqTime)

				err = sendCount(ctx, players, apigwsvc, []byte{34, 115, 116, 97, 114, 116, 34, 125}) // "start"}
				if err != nil {
					return err
				}

				po, err := ebsvc.PutEvents(ctx, &eventbridge.PutEventsInput{
					Entries: []ebtypes.PutEventsRequestEntry{
						{
							Detail:     aws.String("{\"gameno\":" + "\"" + sk + "\"" + "}"),
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

				err = sendCount(ctx, players, apigwsvc, []byte{count, 125}) // 7}
				if err != nil {
					return err
				}
			} else {
				ticker.Stop()
			}

			return nil
		}
	}
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

	return getReturnValue(http.StatusBadRequest), err
}
