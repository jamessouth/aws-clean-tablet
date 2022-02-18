package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cog "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func handler(ctx context.Context, ev events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {

	fmt.Println("ev", ev)

	var (
		head     = ev.CognitoEventUserPoolsHeader
		req      = ev.Request
		username = head.UserName
	)

	if req.ClientMetadata["key"] == "forgotpassword" {
		reg := head.Region

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(reg),
		)
		if err != nil {
			return ev, err
		}

		svc := cog.NewFromConfig(cfg)
		attr := "username"

		lu, err := svc.ListUsers(ctx, &cog.ListUsersInput{
			UserPoolId:      aws.String(head.UserPoolID),
			AttributesToGet: []string{},
			Filter:          aws.String(fmt.Sprintf("%s = %q", attr, username)),
			Limit:           aws.Int32(1),
		})
		if err != nil {
			return ev, err
		}

		user := *lu
		if len(user.Users) < 1 {
			return ev, errors.New(attr + " not found")
		}
		status := user.Users[0].UserStatus
		fmt.Printf("\n%s, %+v\n", "users", status)

		if status != types.UserStatusTypeConfirmed {
			return ev, errors.New("user not confirmed - status: " + string(status))
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

	if req.ClientMetadata["key"] == "forgotusername" {
		reg := head.Region

		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(reg),
		)
		if err != nil {
			return ev, err
		}

		svc := cog.NewFromConfig(cfg)
		attr := "email"

		lu, err := svc.ListUsers(ctx, &cog.ListUsersInput{
			UserPoolId:      aws.String(head.UserPoolID),
			AttributesToGet: []string{},
			Filter:          aws.String(fmt.Sprintf("%s = %q", attr, email)),
			Limit:           aws.Int32(1),
		})
		if err != nil {
			return ev, err
		}

		res := *lu
		if len(res.Users) < 1 {
			return ev, errors.New(attr + " not found")
		}
		user := res.Users[0]
		// status := user.UserStatus
		name := user.Username
		fmt.Printf("\n%s, %+v\n", "users2", user)

		// if status != types.UserStatusTypeConfirmed {
		// 	return ev, errors.New("user not confirmed - status: " + string(status))
		// }

		fp, err := svc.ForgotPassword(ctx, &cog.ForgotPasswordInput{
			ClientId: aws.String(head.CallerContext.ClientID),
			Username: name,
		})
		if err != nil {
			return ev, err
		}

		fpo := *fp
		fmt.Printf("\n%s, %+v\n", "fp", fpo)

		return ev, errors.New("user email found")
	}

	return ev, nil
}

func main() {
	lambda.Start(handler)
}
