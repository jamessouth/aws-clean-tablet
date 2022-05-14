package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cog "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	minEmailLength    int = 5
	maxEmailLength    int = 99
	minUsernameLength int = 3
	maxUsernameLength int = 10
)

var body struct {
	Swear []string
}

func getBadWords(b io.ReadCloser) ([]string, error) {
	defer b.Close()

	rawBytes, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawBytes, &body)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return nil, err
	}

	return body.Swear, nil
}

func handler(ctx context.Context, ev events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {

	fmt.Println("ev", ev)

	ev.Response.AutoConfirmUser = false
	ev.Response.AutoVerifyEmail = false
	ev.Response.AutoVerifyPhone = false

	var (
		bucket     = os.Getenv("bucket")
		swear      = os.Getenv("swear")
		swearETag  = os.Getenv("swearETag")
		head       = ev.CognitoEventUserPoolsHeader
		req        = ev.Request
		username   = head.UserName
		reg        = head.Region
		upid       = head.UserPoolID
		nameRegex  = regexp.MustCompile(`\W`)
		emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		badWords   []string
	)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return ev, err
	}

	if len(username) < minUsernameLength || len(username) > maxUsernameLength {
		return ev, errors.New("username must be 3-10 characters long")
	}

	if nameRegex.MatchString(username) {
		return ev, errors.New("username must be letters, numbers, and underscores only; no whitespace or symbols")
	}

	s3svc := s3.NewFromConfig(cfg)

	obj, err := s3svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket:  aws.String(bucket),
		Key:     aws.String(swear),
		IfMatch: aws.String(swearETag),
	})

	if err != nil {
		return ev, err
	}

	objOutput := *obj
	fmt.Printf("\n%s, %+v\n", "getObj op", objOutput)

	eTag := *objOutput.ETag
	if eTag != swearETag {
		fmt.Println("eTags do not match", eTag, swearETag)
		return ev, errors.New(fmt.Sprintf("eTags do not match: %s != %s", eTag, swearETag))
	} else {
		badWords, err = getBadWords(objOutput.Body)
		if err != nil {
			return ev, err
		}
	}

	for _, w := range badWords {
		if strings.Contains(username, w) {
			return ev, errors.New("unacceptable username; please submit another")
		}
	}

	if req.ClientMetadata["key"] == "forgotpassword" {

		svc := cog.NewFromConfig(cfg)
		attr := "username"

		lu, err := svc.ListUsers(ctx, &cog.ListUsersInput{
			UserPoolId:      aws.String(upid),
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

	email, ok := req.UserAttributes["email"]
	if ok {
		if len(email) < minEmailLength || len(email) > maxEmailLength || !emailRegex.MatchString(email) {
			return ev, errors.New("a properly formatted email address is required")
		}
	} else {
		return ev, errors.New("email attribute not present")
	}

	if req.ClientMetadata["key"] == "forgotusername" {

		svc := cog.NewFromConfig(cfg)
		attr := "email"

		lu, err := svc.ListUsers(ctx, &cog.ListUsersInput{
			UserPoolId:      aws.String(upid),
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
		name := *user.Username
		fmt.Printf("\n%s, %+v\n", "users2", user)

		return ev, errors.New("user found - " + name)
	}

	return ev, nil
}

func main() {
	lambda.Start(handler)
}
