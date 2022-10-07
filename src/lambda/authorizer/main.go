package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
)



	
type keySetFetcher func(ctx context.Context, s3svc *s3.Client, s3in s3.GetObjectInput) (jwk.Set, error)

func getKeySet(ctx context.Context, s3svc *s3.Client, s3in s3.GetObjectInput) (jwk.Set, error) {
	obj, err := s3svc.GetObject(ctx, &s3in)

	if err != nil {
		return nil, err
	}

	objOutput := *obj
	// fmt.Printf("\n%s, %+v\n", "getObj op", objOutput)



	return jwk.ParseReader(objOutput.Body)
	
}



type keyHandler struct {
	ctx context.Context
	s3svc *s3.Client
	s3in *s3.GetObjectInput
	fetcher keySetFetcher
}

func (h *keyHandler) fetchKeys(ctx context.Context, sink jws.KeySink, sig *jws.Signature, msg *jws.Message) error {

	
	alg := sig.ProtectedHeaders().Algorithm()
	kid := sig.ProtectedHeaders().KeyID()

	
	key, err := h.fetcher(ctx, h.s3svc, h.s3in)
	if err != nil {
		return err
	}




	key, keyPresent := jwkSet.LookupKeyID(kid)

	if !keyPresent {

	} else {

	}



	sink.Key(alg, key)
	return nil
}






	
	 jwt.Parse(jwtRaw, jwt.WithKeyProvider(handler), jwt.WithContext(c))




func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	// fmt.Printf("%s: %+v\n", "request", req.QueryStringParameters["auth"])

	var (
		tableName   = os.Getenv("tableName")
		appClientID = os.Getenv("appClientID")
		userPoolID  = os.Getenv("userPoolID")
		origin      = os.Getenv("origin")

		bucket   = aws.String(os.Getenv("bucket"))
		jwksKey  = aws.String(os.Getenv("jwksKey"))
		jwksETag = os.Getenv("jwksETag")
	)

	if req.Headers["Origin"] != origin {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong domain")
	}

	if len(req.Headers["User-Agent"]) < 10 {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("header error - request from wrong client")
	}

	region := strings.Split(req.MethodArn, ":")[3]

	accessToken := []byte(req.QueryStringParameters["auth"])

	msg, err := jws.Parse(accessToken)
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


	kh := keyHandler{
		s3svc: s3svc,
		s3in: &s3.GetObjectInput{
			Bucket:  bucket,
			Key:     jwksKey,
			IfMatch: aws.String(jwksETag),
		},
		fetcher: getKeySet,

	}

	err = kh.fetchKeys(ctx, interface {
		Key(jwa.SignatureAlgorithm, interface{})
	}, msg.Signatures()[0], msg)



	





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

	// keyset = jwk.p

	// for _, k := range jwks.Keys {
	// 	if k.Kid == kid {

	// 	}
	// }

	// parsedToken, err := jwt.Parse(
	// 	token,
	// 	jwt.WithKeySet(keyset),
	// 	jwt.WithValidate(true),
	// 	jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
	// 	jwt.WithClaimValue("client_id", appClientID),
	// 	jwt.WithClaimValue("token_use", "access"),
	// )
	// if err != nil {
	// 	return createPolicy(
	// 		req.MethodArn,
	// 		"Deny",
	// 		"ID",
	// 		map[string]interface{}{
	// 			"error": getErrorMsg(err),
	// 		},
	// 	), err
	// }

	// fmt.Println(parsedToken)

	// return createPolicy(
	// 	req.MethodArn,
	// 	"Allow",
	// 	parsedToken.Subject(),
	// 	map[string]interface{}{
	// 		"username": parsedToken.PrivateClaims()["username"].(string),
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
	p.UsageIdentifierKey = ""

	return
}
