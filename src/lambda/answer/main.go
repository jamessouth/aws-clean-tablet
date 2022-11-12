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
	newline   byte = 10
	byteRange int  = 25      //get 26 bytes, enough for a 12-letter word no matter what comes before or after
	highByte  int  = 1519766 //file size - 26
)

var (
	answerRE = regexp.MustCompile(`(?i)^[a-z]{1}[a-z ]{0,10}[a-z]{1}$`)
	gamenoRE = regexp.MustCompile(`^\d{19}$`)
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

func checkInput(s string, re *regexp.Regexp) string {
	if re.MatchString(s) {
		return s
	}

	return ""
}

func checkKeys(s string) error {
	if strings.Count(s, "gameno") != 1 || strings.Count(s, "aW5mb3Jt") != 1 {
		return errors.New("improper json input - duplicate or missing key")
	}

	return nil
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	bod := req.Body

	fmt.Println("answer", bod)

	err := checkKeys(bod)
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
		// tableName = os.Getenv("tableName")
		bucket    = os.Getenv("bucket")
		words     = os.Getenv("words")
		wordsETag = os.Getenv("wordsETag")
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		s3svc     = s3.NewFromConfig(cfg)
		body      struct {
			Gameno, AW5mb3Jt string
		}
		ans           string
		auth          = req.RequestContext.Authorizer.(map[string]interface{})
		id, tableName = auth["principalId"].(string), auth["tableName"].(string)
	)

	err = json.Unmarshal([]byte(bod), &body)
	if err != nil {
		return callErr(err)
	}

	checkedGameno := checkInput(body.Gameno, gamenoRE)
	if checkedGameno == "" {
		return callErr(errors.New("improper json input - wrong gameno"))
	}

	checkedAnswer := checkInput(body.AW5mb3Jt, answerRE)
	if checkedAnswer == "" {
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
		ans = checkedAnswer
	}

	//condition on player name??

	_, err = ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LIVEGAME"},
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
