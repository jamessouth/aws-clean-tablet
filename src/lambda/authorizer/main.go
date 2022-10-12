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
	"github.com/aws/aws-sdk-go/aws"

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
		},
		)

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
		// tableName   = os.Getenv("tableName")
		appClientID = os.Getenv("appClientID")
		userPoolID  = os.Getenv("userPoolID")
		origin      = os.Getenv("origin")
		bucket      = os.Getenv("bucket")
		jwksKey     = os.Getenv("jwksKey")
	)

	if req.Headers["Origin"] != origin {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong domain")
	}

	if len(req.Headers["User-Agent"]) < 10 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong client")
	}

	region := strings.Split(req.MethodArn, ":")[3]

	accessToken := []byte(req.QueryStringParameters["authToken"])
	idToken := []byte(req.QueryStringParameters["idToken"])

	authMsg, err := jws.Parse(accessToken)
	if err != nil {
		fmt.Printf(`failed to parse serialized JWT: %s`, err)
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}

	idMsg, err := jws.Parse(idToken)
	if err != nil {
		fmt.Printf(`failed to parse serialized JWT: %s`, err)
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}

	var (
		s3svc = s3.NewFromConfig(cfg)

		// ddbsvc = dynamodb.NewFromConfig(cfg)
	)

	authKH := &keyHandler{
		reg:      region,
		upid:     userPoolID,
		s3svc:    s3svc,
		s3bucket: bucket,
		s3key:    jwksKey,
		fetcher:  getKeySet,
	}

	idKH := &keyHandler{
		reg:      region,
		upid:     userPoolID,
		s3svc:    s3svc,
		s3bucket: bucket,
		s3key:    jwksKey,
		fetcher:  getKeySet,
	}

	ks := make(sink)

	err = authKH.FetchKeys(ctx, ks, authMsg.Signatures()[0], authMsg)
	if err != nil {
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}

	parsedAccessToken, err := jwt.Parse(
		accessToken,
		jwt.WithContext(ctx),
		jwt.WithKeyProvider(authKH),
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithClaimValue("client_id", appClientID),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}
	fmt.Println(parsedAccessToken)

	err = idKH.FetchKeys(ctx, ks, idMsg.Signatures()[0], idMsg)
	if err != nil {
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}

	parsedIdToken, err := jwt.Parse(
		accessToken,
		jwt.WithContext(ctx),
		jwt.WithKeyProvider(idKH),
		jwt.WithValidate(true),
		jwt.WithVerify(true),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithClaimValue("client_id", appClientID),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return createPolicy(
			req.MethodArn,
			"Deny",
			"ID",
			map[string]interface{}{
				"error": getErrorMsg(err),
			},
		), err
	}
	fmt.Println(parsedIdToken)

	// gi, err := ddbsvc.GetItem(ctx, &dynamodb.GetItemInput{
	// 	Key: map[string]types.AttributeValue{
	// 		"pk": &types.AttributeValueMemberS{Value: "JWKS"},
	// 		"sk": &types.AttributeValueMemberS{Value: "keys"},
	// 	},
	// 	TableName: aws.String(tableName),
	// 	// ProjectionExpression: aws.String("keys"),
	// })
	// if err != nil {
	// 	callErr(err)
	// }

	// fmt.Printf("%s: %+v\n", "gi", gi)

	// var jwks struct {
	// 	Keys []struct {
	// 		Kty, Alg, E, N, Use, Kid string
	// 	}
	// }
	// err = attributevalue.UnmarshalMap(gi.Item, &jwks)
	// if err != nil {
	// 	callErr(err)
	// }

	// return createPolicy(
	// 	req.MethodArn,
	// 	"Allow",
	// 	parsedAccessToken.Subject(),
	// 	map[string]interface{}{
	// 		"username": parsedAccessToken.PrivateClaims()["username"].(string),
	// 	},
	// ), nil

	return createPolicy(
		req.MethodArn,
		"Deny",
		"ID",
		map[string]interface{}{
			"error": getErrorMsg(err),
		},
	), err

}

func main() {
	lambda.Start(handler)
}

func getErrorMsg(e error) string {
	clause := " not satisfied"
	switch e.Error() {
	case "exp" + clause:
		return "Token expired"
	case "iss" + clause:
		return "Wrong issuer"
	case "client_id" + clause:
		return "Wrong app client ID"
	case "token_use" + clause:
		return "Wrong token type"
	default:
		return e.Error()
	}
}

func createPolicy(arn, effect, pID string, context map[string]interface{}) (p events.APIGatewayCustomAuthorizerResponse) {
	p.PrincipalID = pID
	p.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Effect:   effect,
				Action:   []string{"execute-api:Invoke"},
				Resource: []string{arn},
			},
		},
	}
	p.Context = context

	return
}
