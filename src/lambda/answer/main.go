package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

const (
	newline   byte   = 10
	byteRange int    = 25      //get 26 bytes, enough for a 12-letter word no matter what comes before or after
	highByte  int    = 1519766 //file size - 26
	liveGame  string = "LIVEGAME"
)

func getRandomByte() string {
	t := time.Now().UnixNano()
	rand.Seed(t)

	randobyte := rand.Intn(highByte)

	return fmt.Sprintf("bytes=%d-%d", randobyte, randobyte+byteRange)
}

func getWord(b io.ReadCloser) string {
	defer b.Close()

	rawBytes, err := io.ReadAll(b)
	if err != nil {
		return ""
	}

	words := bytes.Split(rawBytes, []byte{newline})

	return string(words[1])
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
		maxLength       = 99
		gamenoRE        = regexp.MustCompile(`^\d{19}$`)
		aW5mb3JtRE      = regexp.MustCompile(`(?i)^[a-z]{1}[a-z ]{0,10}[a-z]{1}$`)
		body            struct{ Gameno, AW5mb3Jt string }
		checkedAW5mb3Jt string
	)

	if len(s) > maxLength {
		return "", "", fmt.Errorf("improper json input - too long: %d", len(s))
	}

	if strings.Count(s, "gameno") != 1 || strings.Count(s, "aW5mb3Jt") != 1 {
		return "", "", errors.New("improper json input - duplicate/missing key")
	}

	err := json.Unmarshal([]byte(s), &body)
	if err != nil {
		return "", "", err
	}

	var gameno, aW5mb3Jt = body.Gameno, body.AW5mb3Jt

	if !gamenoRE.MatchString(gameno) {
		return "", "", errors.New("improper json input - bad gameno: " + gameno)
	}

	if aW5mb3JtRE.MatchString(aW5mb3Jt) {
		checkedAW5mb3Jt = aW5mb3Jt
	}

	return gameno, checkedAW5mb3Jt, nil
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		bod    = req.Body
		region = strings.Split(req.RequestContext.DomainName, ".")[2]
	)

	checkedGameno, checkedAW5mb3Jt, err := checkInput(bod)
	if err != nil {
		return callErr(err)
	}

	fmt.Println("answer", bod, len(bod))

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return callErr(err)
	}

	var (
		bucket        = os.Getenv("bucket")
		words         = os.Getenv("words")
		wordsETag     = os.Getenv("wordsETag")
		ddbsvc        = dynamodb.NewFromConfig(cfg)
		s3svc         = s3.NewFromConfig(cfg)
		ans           string
		auth          = req.RequestContext.Authorizer.(map[string]interface{})
		id, tableName = auth["principalId"].(string), auth["tableName"].(string)
	)

	if checkedAW5mb3Jt == "" {
		obj, err := s3svc.GetObject(ctx, &s3.GetObjectInput{
			Bucket:  aws.String(bucket),
			Key:     aws.String(words),
			IfMatch: aws.String(wordsETag),
			Range:   aws.String(getRandomByte()),
		})
		if err != nil {
			return callErr(err)
		}

		objOutput := *obj
		// fmt.Printf("\n%s, %+v\n", "getObj op", objOutput)

		ans = getWord(objOutput.Body)
	} else {
		ans = checkedAW5mb3Jt
	}

	//condition on player name??

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: liveGame},
			"sk": &types.AttributeValueMemberS{Value: checkedGameno},
		},
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "players",
			"#I": id,
			"#A": "Answer",
			"#C": "answersCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberS{Value: ans},
			":o": &types.AttributeValueMemberN{Value: "1"},
		},
		UpdateExpression: aws.String("SET #P.#I.#A = :a ADD #C :o"),
	})
	if err != nil {
		return callErr(err)
	}

	return getReturnValue(http.StatusOK), nil
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

	return getReturnValue(http.StatusBadRequest), err
}
