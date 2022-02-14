package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ev events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {

	fmt.Println("ev", ev)

	if ev.Request.ClientMetadata["key"] == "fp" {

	}

	var (
		nameRegex  = regexp.MustCompile(`\W`)
		emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		username   = ev.CognitoEventUserPoolsHeader.UserName
	)

	ev.Response.AutoConfirmUser = false
	ev.Response.AutoVerifyEmail = false
	ev.Response.AutoVerifyPhone = false

	if len(username) < 3 || len(username) > 10 {
		return ev, errors.New("username must be 3-10 characters long")
	}

	if nameRegex.MatchString(username) {
		return ev, errors.New("username must be letters, numbers, and underscores only; no whitespace or symbols")
	}

	email, ok := ev.Request.UserAttributes["email"]
	if ok {
		if len(email) < 5 || !emailRegex.MatchString(email) {
			return ev, errors.New("a properly formatted email address is required")
		}
	} else {
		return ev, errors.New("email attribute not present")
	}

	return ev, nil
}

func main() {
	lambda.Start(handler)
}
