package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type sink map[jwa.SignatureAlgorithm]interface{}

func (s sink) Key(alg jwa.SignatureAlgorithm, key interface{}) {
	s[alg] = key
}

type keySetFetcher func(ctx context.Context, s3svc *s3.Client, s3in *s3.GetObjectInput) (jwk.Set, error)

func getKeySet(ctx context.Context, s3svc *s3.Client, s3in *s3.GetObjectInput) (jwk.Set, error) {
	obj, err := s3svc.GetObject(ctx, s3in)
	if err != nil {
		return nil, err
	}

	objOutput := *obj
	// fmt.Printf("\n%s, %+v\n", "getObj op", objOutput)

	return jwk.ParseReader(objOutput.Body)
}

type keyHandler struct {
	reg, upid, s3bucket, s3key string
	s3svc                      *s3.Client
	fetcher                    keySetFetcher
}

func (h *keyHandler) FetchKeys(ctx context.Context, sink jws.KeySink, sig *jws.Signature, msg *jws.Message) error {
	alg := sig.ProtectedHeaders().Algorithm()
	kid := sig.ProtectedHeaders().KeyID()

	jwkSet, err := h.fetcher(ctx, h.s3svc, &s3.GetObjectInput{
		Bucket: aws.String(h.s3bucket),
		Key:    aws.String(h.s3key),
	})
	if err != nil {
		return err
	}

	var finalKey jwk.Key
	key, keyPresent := jwkSet.LookupKeyID(kid)

	if !keyPresent {
		keyset, err := jwk.Fetch(ctx, "https://cognito-idp."+h.reg+".amazonaws.com/"+h.upid+"/.well-known/jwks.json")
		if err != nil {
			return err
		}

		marshalledKeyset, err := json.Marshal(keyset)
		if err != nil {
			return err
		}

		dig := md5.New()
		dig.Write(marshalledKeyset)
		md5 := base64.StdEncoding.EncodeToString(dig.Sum(nil))

		_, err = h.s3svc.PutObject(ctx, &s3.PutObjectInput{
			Bucket:     aws.String(h.s3bucket),
			Key:        aws.String(h.s3key),
			Body:       bytes.NewReader(marshalledKeyset),
			ContentMD5: aws.String(md5),
		})
		if err != nil {
			return err
		}

		finalKey, _ = keyset.LookupKeyID(kid)
	} else {
		finalKey = key
	}

	sink.Key(alg, finalKey)
	return nil
}

func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	// fmt.Printf("%s: %+v\n", "request", req.QueryStringParameters["auth"])

	var (
		tableName   = os.Getenv("tableName")
		appClientID = os.Getenv("appClientID")
		userPoolID  = os.Getenv("userPoolID")
		origin      = os.Getenv("origin")
		bucket      = os.Getenv("bucket")
		jwksKey     = os.Getenv("jwksKey")
		region      = strings.Split(req.MethodArn, ":")[3]
		accessToken = []byte(req.QueryStringParameters["authToken"])
		idToken     = []byte(req.QueryStringParameters["idToken"])
	)

	if req.Headers["Origin"] != origin {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong domain")
	}

	if len(req.Headers["User-Agent"]) < 10 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong client")
	}

	authMsg, err := jws.Parse(accessToken)
	if err != nil {
		fmt.Printf(`failed to parse serialized JWT: %s`, err)
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	idMsg, err := jws.Parse(idToken)
	if err != nil {
		fmt.Printf(`failed to parse serialized JWT: %s`, err)
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	var (
		s3svc  = s3.NewFromConfig(cfg)
		ddbsvc = dynamodb.NewFromConfig(cfg)
		kh     = &keyHandler{
			reg:      region,
			upid:     userPoolID,
			s3svc:    s3svc,
			s3bucket: bucket,
			s3key:    jwksKey,
			fetcher:  getKeySet,
		}
		ks = make(sink)
	)

	err = kh.FetchKeys(ctx, ks, authMsg.Signatures()[0], authMsg)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	parsedAccessToken, err := jwt.Parse(
		accessToken,
		jwt.WithContext(ctx),
		jwt.WithKeyProvider(kh),
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithClaimValue("client_id", appClientID),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, err
	}

	fmt.Println(parsedAccessToken)

	sub := parsedAccessToken.Subject()

	err = kh.FetchKeys(ctx, ks, idMsg.Signatures()[0], idMsg)
	if err != nil {
		deny(req.MethodArn, sub, err)
	}

	gi, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "TOKEN"},
			"sk": &types.AttributeValueMemberS{Value: sub},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		deny(req.MethodArn, sub, err)
	}

	fmt.Printf("%s: %+v\n", "gi", gi)

	var specialClaim struct {
		Uid string
	}
	err = attributevalue.UnmarshalMap(gi.Item, &specialClaim)
	if err != nil {
		deny(req.MethodArn, sub, err)
	}

	parsedIdToken, err := jwt.Parse(
		idToken,
		jwt.WithContext(ctx),
		jwt.WithAudience(appClientID),
		jwt.WithKeyProvider(kh),
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithClaimValue("q", specialClaim.Uid),
		jwt.WithClaimValue("token_use", "id"),
	)
	if err != nil {
		deny(req.MethodArn, sub, err)
	}

	fmt.Println(parsedIdToken)

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: sub,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{req.MethodArn},
			}},
		},
		Context: map[string]interface{}{"username": parsedAccessToken.PrivateClaims()["username"].(string)},
	}, nil
}

func main() {
	lambda.Start(handler)
}

func deny(arn, pID string, err error) (events.APIGatewayCustomAuthorizerResponse, error) {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: pID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   "Deny",
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{arn},
				},
			},
		},
	}, err
}
