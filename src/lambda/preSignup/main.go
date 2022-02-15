package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cog "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

const (
	minEmailLength    int = 5
	minUsernameLength int = 3
	maxUsernameLength int = 10
)

func handler(ctx context.Context, ev events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {

	fmt.Println("ev", ev)

	ev.Response.AutoConfirmUser = false
	ev.Response.AutoVerifyEmail = false
	ev.Response.AutoVerifyPhone = false

	var (
		head       = ev.CognitoEventUserPoolsHeader
		req        = ev.Request
		username   = head.UserName
		nameRegex  = regexp.MustCompile(`\W`)
		emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	)

	if len(username) < minUsernameLength || len(username) > maxUsernameLength {
		return ev, errors.New("username must be 3-10 characters long")
	}

	if nameRegex.MatchString(username) {
		return ev, errors.New("username must be letters, numbers, and underscores only; no whitespace or symbols")
	}

	if req.ClientMetadata["key"] == "fp" {
		reg := head.Region

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(reg),
		)
		if err != nil {
			return ev, err
		}

		svc := cog.NewFromConfig(cfg)

		lu, err := svc.ListUsers(ctx, &cog.ListUsersInput{
			UserPoolId:      aws.String(head.UserPoolID),
			AttributesToGet: []string{},
			Filter:          aws.String(fmt.Sprintf("%s = %q", "username", username)),
			Limit:           aws.Int32(1),
		})
		if err != nil {
			return ev, err
		}

		user := *lu
		if len(user.Users) < 1 {
			return ev, errors.New("username not found")
		}
		status := user.Users[0].UserStatus
		fmt.Printf("\n%s, %+v\n", "users", status)

		if status != types.UserStatusTypeConfirmed {
			return ev, errors.New("user not confirmed - status is " + string(status))
		}

		fp, err := svc.ForgotPassword(ctx, &cog.ForgotPasswordInput{
			ClientId: aws.String(head.CallerContext.ClientID),
			Username: user.Users[0].Username,
		})
		if err != nil {
			return ev, err
		}

		fpo := *fp
		fmt.Printf("\n%s, %+v\n", "fp", fpo)

		return ev, errors.New("user found")
	}

	email, ok := req.UserAttributes["email"]
	if ok {
		if len(email) < minEmailLength || !emailRegex.MatchString(email) {
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
