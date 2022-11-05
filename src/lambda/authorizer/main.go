package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
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
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const apigwCustomAuthorizerPolicyVersion string = "2012-10-17"

type sink map[jwa.SignatureAlgorithm]interface{}

type keySetFetcher func(ctx context.Context, s3svc *s3.Client, s3in *s3.GetObjectInput) (jwk.Set, error)

type keyHandler struct {
	reg, upid, s3bucket, s3key string
	s3svc                      *s3.Client
	fetcher                    keySetFetcher
}

func (s sink) Key(alg jwa.SignatureAlgorithm, key interface{}) {
	s[alg] = key
}

func getKeySet(ctx context.Context, s3svc *s3.Client, s3in *s3.GetObjectInput) (jwk.Set, error) {
	obj, err := s3svc.GetObject(ctx, s3in)
	if err != nil {
		return nil, err
	}

	objOutput := *obj

	return jwk.ParseReader(objOutput.Body)
}

func (h *keyHandler) FetchKeys(ctx context.Context, sink jws.KeySink, sig *jws.Signature, msg *jws.Message) error {
	jwkSet, err := h.fetcher(ctx, h.s3svc, &s3.GetObjectInput{
		Bucket: aws.String(h.s3bucket),
		Key:    aws.String(h.s3key),
	})
	if err != nil {
		return err
	}

	var (
		kid             = sig.ProtectedHeaders().KeyID()
		alg             = sig.ProtectedHeaders().Algorithm()
		key, keyPresent = jwkSet.LookupKeyID(kid)
		finalKey        jwk.Key
	)
	if alg != jwa.RS256 {
		return errors.New("incorrect algorithm")
	}

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

	// fmt.Printf("%s: %+v\n", "request", req)

	var (
		tableName          = os.Getenv("tableName")
		appClientID        = os.Getenv("appClientID")
		userPoolID         = os.Getenv("userPoolID")
		origin             = os.Getenv("origin")
		bucket             = os.Getenv("bucket")
		jwksKey            = os.Getenv("jwksKey")
		region             = strings.Split(req.MethodArn, ":")[3]
		token              = []byte(req.QueryStringParameters["auth"])
		unsafeTokenPayload struct{ Sub string }
	)

	if req.Headers["Origin"] != origin {
		return deny(errors.New("header error - request from wrong domain"))
	}

	if len(req.Headers["User-Agent"]) < 12 {
		return deny(errors.New("header error - request from wrong client"))
	}

	msg, err := jws.Parse(token)
	if err != nil {
		return deny(err)
	}

	err = json.Unmarshal(msg.Payload(), &unsafeTokenPayload)
	if err != nil {
		return deny(err)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return deny(err)
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
		ks           = make(sink)
		sub          = unsafeTokenPayload.Sub
		specialClaim struct{ Uid string }
		validator    = jwt.ValidatorFunc(func(_ context.Context, t jwt.Token) jwt.ValidationError {
			ageLimit := 10.0
			if time.Since(t.IssuedAt()).Seconds() > ageLimit {
				return jwt.NewValidationError(errors.New("token too old"))
			}
			return nil
		})
	)

	err = kh.FetchKeys(ctx, ks, msg.Signatures()[0], msg)
	if err != nil {
		return deny(err)
	}

	di, err := ddbsvc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "TOKEN"},
			"sk": &types.AttributeValueMemberS{Value: sub},
		},
		TableName:    aws.String(tableName),
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return deny(err)
	}
	if len(di.Attributes) == 0 {
		return deny(errors.New("player token not found"))
	}

	err = attributevalue.UnmarshalMap(di.Attributes, &specialClaim)
	if err != nil {
		return deny(err)
	}

	parsedToken, err := jwt.Parse(
		token,
		jwt.WithAudience(appClientID),
		jwt.WithClaimValue("email_verified", true),
		jwt.WithClaimValue("q", specialClaim.Uid),
		jwt.WithClaimValue("token_use", "id"),
		jwt.WithContext(ctx),
		jwt.WithIssuer("https://cognito-idp."+region+".amazonaws.com/"+userPoolID),
		jwt.WithKeyProvider(kh),
		jwt.WithMaxDelta(3600*time.Second, jwt.ExpirationKey, jwt.IssuedAtKey),
		jwt.WithMinDelta(3600*time.Second, jwt.ExpirationKey, jwt.IssuedAtKey),
		jwt.WithRequiredClaim("auth_time"),
		jwt.WithRequiredClaim("cognito:username"),
		jwt.WithRequiredClaim("email"),
		jwt.WithRequiredClaim("event_id"),
		jwt.WithSubject(sub),
		jwt.WithValidator(validator),
	)
	if err != nil {
		return deny(err)
	}

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: parsedToken.Subject(),
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: apigwCustomAuthorizerPolicyVersion,
			Statement: []events.IAMPolicyStatement{{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{req.MethodArn},
			}},
		},
		Context: map[string]interface{}{
			"username":  parsedToken.PrivateClaims()["cognito:username"].(string),
			"tableName": tableName,
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}

func deny(err error) (events.APIGatewayCustomAuthorizerResponse, error) {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: apigwCustomAuthorizerPolicyVersion,
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   "Deny",
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{"*"},
				},
			},
		},
	}, err
}
