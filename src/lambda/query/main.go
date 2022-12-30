package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
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

const (
	leaders   string = "leaders"
	listGames string = "listGames"
	listGame  string = "LISTGAME"
	stat_     string = "STAT"
)

type stat struct {
	Name   string  `json:"name"`
	Wins   int     `json:"wins"`
	Points int     `json:"points"`
	Games  int     `json:"games"`
	WinPct float64 `json:"winPct"`
	PPG    float64 `json:"ppg"`
}

type stats []stat

type listPlayer struct {
	Name   string `json:"name" dynamodbav:"name"`
	ConnID string `json:"connid" dynamodbav:"connid"`
}

type frontListGame struct {
	No        string       `json:"no"`
	TimerCxld bool         `json:"timerCxld"`
	Players   []listPlayer `json:"players"`
}

type backListGame struct {
	Pk        string                `dynamodbav:"pk"` //'LISTGAME'
	Sk        string                `dynamodbav:"sk"` //no
	TimerCxld bool                  `dynamodbav:"timerCxld"`
	Players   map[string]listPlayer `dynamodbav:"players"`
}

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

func getSlice(m map[string]listPlayer) (res []listPlayer) {
	for _, v := range m {
		res = append(res, v)
	}

	return
}

func sortByName(players []listPlayer) []listPlayer {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})

	return players
}

func getFrontListGames(gl []backListGame) (res []frontListGame) {
	for _, g := range gl {
		pls := sortByName(getSlice(g.Players))
		res = append(res, frontListGame{
			No:        g.Sk,
			TimerCxld: g.TimerCxld,
			Players:   pls,
		})
	}

	return
}

func checkInput(s string) (string, error) {
	var (
		maxLength = 99
		commandRE = regexp.MustCompile(`^leaders$|^listGames$`)
		body      struct{ Command string }
	)

	if len(s) > maxLength {
		return "", fmt.Errorf("improper json input - too long: %d", len(s))
	}

	if strings.Count(s, "command") != 1 {
		return "", errors.New("improper json input - duplicate/missing key")
	}

	err := json.Unmarshal([]byte(s), &body)
	if err != nil {
		return "", err
	}

	var command = body.Command

	if !commandRE.MatchString(command) {
		return "", errors.New("improper json input - bad command: " + command)
	}

	return command, nil
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

	checkedCommand, err := checkInput(bod)
	if err != nil {
		return callErr(err)
	}

	fmt.Println("query", bod, len(bod))

	var (
		region          = strings.Split(req.RequestContext.DomainName, ".")[2]
		connID          = req.RequestContext.ConnectionID
		apiid           = os.Getenv("CT_APIID")
		stage           = os.Getenv("CT_STAGE")
		endpoint        = "https://" + apiid + ".execute-api." + region + ".amazonaws.com/" + stage
		auth            = req.RequestContext.Authorizer.(map[string]interface{})
		name, tableName = auth["username"].(string), auth["tableName"].(string)
		customResolver  = aws.EndpointResolverWithOptionsFunc(func(service, awsRegion string, options ...interface{}) (aws.Endpoint, error) {
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
		payload  []byte
		qi       = dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			ScanIndexForward:       aws.Bool(false),
			KeyConditionExpression: aws.String("pk = :e"),
			Limit:                  aws.Int32(50),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":e": &types.AttributeValueMemberS{Value: listGame},
			},
		}
	)

	if checkedCommand == leaders {
		qi.ScanIndexForward = aws.Bool(true)
		qi.Limit = aws.Int32(100)
		qi.ExpressionAttributeValues[":e"] = &types.AttributeValueMemberS{Value: stat_}

		leadersResults, err := ddbsvc.Query(ctx, &qi)
		if err != nil {
			return callErr(err)
		}

		var leaders stats
		err = attributevalue.UnmarshalListOfMaps(leadersResults.Items, &leaders)
		if err != nil {
			return callErr(err)
		}

		payload, err = json.Marshal(struct {
			Leaders []stat `json:"leaders"`
		}{
			Leaders: leaders.sortByWinsThenName().calcStats(),
		})
		if err != nil {
			return callErr(err)
		}
	} else if checkedCommand == listGames {
		listGamesResults, err := ddbsvc.Query(ctx, &qi)
		if err != nil {
			return callErr(err)
		}

		var listGames []backListGame
		err = attributevalue.UnmarshalListOfMaps(listGamesResults.Items, &listGames)
		if err != nil {
			return callErr(err)
		}

		payload, err = json.Marshal(struct {
			ListGames []frontListGame `json:"listGms"`
			Name      string          `json:"name"`
			ConnID    string          `json:"connid"`
		}{
			ListGames: getFrontListGames(listGames),
			Name:      name,
			ConnID:    connID,
		})
		if err != nil {
			return callErr(err)
		}
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
