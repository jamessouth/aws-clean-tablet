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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/smithy-go"
)

type livePlayer struct {
	PlayerID    string `dynamodbav:"playerid"`
	Name        string `dynamodbav:"name"`
	ConnID      string `dynamodbav:"connid"`
	Color       string `dynamodbav:"color"`
	Index       string `dynamodbav:"index"`
	Score       int    `dynamodbav:"score"`
	Answer      string
	HasAnswered bool `dynamodbav:"hasAnswered"`
}

const (
	newline   byte = 10
	byteRange int  = 25      //get 26 bytes, enough for a 12-letter word no matter what comes before or after
	highByte  int  = 1519766 //file size - 26
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

func clearHasAnswered(pl []livePlayer) []livePlayer {
	for i, p := range pl {
		p.HasAnswered = false
		pl[i] = p
	}

	return pl
}

func sanitize(s string) string {
	re := regexp.MustCompile(`(?i)^[a-z ]{2,12}$`)
	re2 := regexp.MustCompile(`^\s|\s$`)

	if re.MatchString(s) {
		if !re2.MatchString(s) {
			return s
		} else {
			return ""
		}
	}

	return ""
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("answer", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return callErr(err)
	}

	var (
		tableName = os.Getenv("tableName")
		bucket    = os.Getenv("bucket")
		words     = os.Getenv("words")
		wordsETag = os.Getenv("wordsETag")
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		s3svc     = s3.NewFromConfig(cfg)
		body      struct {
			Gameno, Answer, Index string
		}
		ans string
	)

	err = json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return callErr(err)
	}

	sanitizedAnswer := sanitize(body.Answer)

	if sanitizedAnswer == "" {

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

		eTag := *objOutput.ETag
		if eTag != wordsETag {
			fmt.Println("eTags do not match", eTag, wordsETag)
			ans = ""
		} else {
			ans = getWord(objOutput.Body)
		}

	} else {
		ans = sanitizedAnswer
	}

	gameItemKey, err := attributevalue.MarshalMap(struct {
		Pk string `dynamodbav:"pk"`
		Sk string `dynamodbav:"sk"`
	}{
		Pk: "LIVEGAME",
		Sk: body.Gameno,
	})
	if err != nil {
		return callErr(err)
	}

	index := body.Index
	//condition on player name??
	ui, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:       gameItemKey,
		TableName: aws.String(tableName),
		ExpressionAttributeNames: map[string]string{
			"#P": "players",
			"#A": "Answer",
			"#C": "answersCount",
			"#H": "hasAnswered",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberS{Value: ans},
			":o": &types.AttributeValueMemberN{Value: "1"},
			":t": &types.AttributeValueMemberBOOL{Value: true},
		},
		UpdateExpression: aws.String("SET #P[" + index + "].#A = :a, #P[" + index + "].#H = :t ADD #C :o"),
		ReturnValues:     types.ReturnValueAllNew,
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
