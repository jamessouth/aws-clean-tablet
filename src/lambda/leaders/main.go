package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

const stat_ string = "STAT"

type stat struct {
	Name   string  `json:"name"`
	Wins   int     `json:"wins"`
	Points int     `json:"points"`
	Games  int     `json:"games"`
	WinPct float64 `json:"winPct"`
	PPG    float64 `json:"ppg"`
}

type stats []stat

func (stats stats) sortByWinsThenName() stats {
	sort.Slice(stats, func(i, j int) bool {
		switch {
		case stats[i].Wins != stats[j].Wins:
			return stats[i].Wins > stats[j].Wins
		default:
			return stats[i].Name < stats[j].Name
		}
	})

	return stats
}

func (stats stats) calcStats() stats {
	for i, s := range stats {
		w := float64(s.Wins)
		g := float64(s.Games)
		p := float64(s.Points)

		s.WinPct = math.Round((w/g)*100) / 100
		s.PPG = math.Round((p/g)*100) / 100

		if math.IsNaN(s.WinPct) || math.IsInf(s.WinPct, 0) {
			s.WinPct = 0
		}
		if math.IsNaN(s.PPG) || math.IsInf(s.PPG, 0) {
			s.PPG = 0
		}

		stats[i] = s
	}

	return stats
}

func getReturnValue(status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        status,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	bod := req.Body

	if len(bod) > 75 { //TODO replace with observed value
		return callErr(errors.New("improper json input - too long"))
	}

	fmt.Println("leaders", bod, len(bod))

	var (
		region    = strings.Split(req.RequestContext.DomainName, ".")[2]
		connID    = req.RequestContext.ConnectionID
		apiid     = os.Getenv("CT_APIID")
		stage     = os.Getenv("CT_STAGE")
		endpoint  = "https://" + apiid + ".execute-api." + region + ".amazonaws.com/" + stage
		tableName = req.RequestContext.Authorizer.(map[string]interface{})["tableName"].(string)
		// id, name  = auth["principalId"].(string), auth["username"].(string)
		customResolver = aws.EndpointResolverWithOptionsFunc(func(service, awsRegion string, options ...interface{}) (aws.Endpoint, error) {
			if service == apigatewaymanagementapi.ServiceID && awsRegion == region {
				ep := aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: awsRegion,
				}

				return ep, nil
			}

			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})
	)

	apigwcfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return callErr(err)
	}

	ddbcfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return callErr(err)
	}

	var (
		apigwsvc = apigatewaymanagementapi.NewFromConfig(apigwcfg)
		ddbsvc   = dynamodb.NewFromConfig(ddbcfg)
	)

	leadersResults, err := ddbsvc.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("pk = :s"),
		Limit:                  aws.Int32(100),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: stat_},
		},
	})
	if err != nil {
		return callErr(err)
	}

	var leaders stats
	err = attributevalue.UnmarshalListOfMaps(leadersResults.Items, &leaders)
	if err != nil {
		return callErr(err)
	}

	fmt.Printf("%s%+v\n", "res ", leaders)

	payload, err := json.Marshal(struct {
		Leaders []stat `json:"leaders"`
	}{
		Leaders: leaders.sortByWinsThenName().calcStats(),
	})
	if err != nil {
		return callErr(err)
	}

	conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connID), Data: payload}

	_, err = apigwsvc.PostToConnection(ctx, &conn)
	if err != nil {
		return callErr(err)
	}

	return getReturnValue(http.StatusOK), nil
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

	return getReturnValue(http.StatusBadRequest), err
}
