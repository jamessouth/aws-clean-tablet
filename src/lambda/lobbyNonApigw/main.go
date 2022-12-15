package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

const (
	maxPlayersPerGame string = "8"
	connect           string = "CONNECT"
	listGame          string = "LISTGAME"
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
		commandRE = regexp.MustCompile(`^join$`)
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
	}

	return gameno, command, nil
}

func joinEvent(ctx context.Context, connKey, gameItemKey map[string]types.AttributeValue, checkedGameno, connid, id, name, tableName string, ddbsvc *dynamodb.Client) error {
	player := listPlayer{
		Name:   name,
		ConnID: connid,
	}

	marshalledPlayersMap, err := attributevalue.Marshal(map[string]listPlayer{
		id: player,
	})
	if err != nil {
		return err
	}

	marshalledPlayer, err := attributevalue.Marshal(player)
	if err != nil {
		return err
	}

	updateConnInput := types.Update{
		Key:                 connKey,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("size (#G) = :z"),
		ExpressionAttributeNames: map[string]string{
			"#G": "game",
			"#R": "returning",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":g": &types.AttributeValueMemberS{Value: checkedGameno},
			":z": &types.AttributeValueMemberN{Value: "0"},
			":f": &types.AttributeValueMemberBOOL{Value: false},
		},
		UpdateExpression: aws.String("SET #G = :g, #R = :f"),
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
		return err
	} //TODO check error here

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
		return err
	}

	return nil
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		bod    = req.Body
		region = strings.Split(req.RequestContext.DomainName, ".")[2]
	)

	fmt.Println("lobbyNonApigw", bod, len(bod))

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
		connKey             = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: connect},
			"sk": &types.AttributeValueMemberS{Value: id},
		}
	)

	if checkedGameno == newgame {
		checkedGameno = fmt.Sprintf("%d", time.Now().UnixNano())
	}

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

	if checkedCommand == join {
		err = joinEvent(ctx, connKey, gameItemKey, checkedGameno, req.RequestContext.ConnectionID, id, name, tableName, ddbsvc)
		if err != nil {
			return callErr(err)
		}
	}

	return getReturnValue(http.StatusOK), nil
}

func main() {
	lambda.Start(handler)
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
