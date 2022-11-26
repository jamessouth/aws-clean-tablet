package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/smithy-go"
)

func checkInput(s string) (string, string, error) {
	var (
		maxLength = 750
		gamenoRE  = regexp.MustCompile(`^\d{19}$`)
		commandRE = regexp.MustCompile(`^.{648}$`)
		body      struct{ Gameno, Command string }
	)

	if len(s) > maxLength {
		return "", "", errors.New("improper json input - too long")
	}

	if strings.Count(s, "gameno") != 1 || strings.Count(s, "command") != 1 {
		return "", "", errors.New("improper json input - duplicate/missing key")
	}

	err := json.Unmarshal([]byte(s), &body)
	if err != nil {
		return "", "", err
	}

	if !gamenoRE.MatchString(body.Gameno) {
		return "", "", errors.New("improper json input - bad gameno")
	}

	if !commandRE.MatchString(body.Command) {
		return "", "", errors.New("improper json input - bad command")
	}

	return body.Gameno, body.Command, nil
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	bod := req.Body

	fmt.Println("answer", bod, len(bod))

	checkedGameno, checkedCommand, err := checkInput(bod)
	if err != nil {
		callErr(err)
	}

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		fmt.Println("cfg err")
	}

	sfnsvc := sfn.NewFromConfig(cfg)

	stsi := sfn.SendTaskSuccessInput{
		Output:    aws.String("\"\""),
		TaskToken: aws.String(checkedCommand),
	}

	_, err = sfnsvc.SendTaskSuccess(ctx, &stsi)
	if err != nil {
		callErr(err)
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

func callErr(err error) {
	if err != nil {

		// To get any API error
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("db error, Code: %v, Message: %v",
				apiErr.ErrorCode(), apiErr.ErrorMessage())
		}

	}
}
