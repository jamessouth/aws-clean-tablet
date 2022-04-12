package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
	"github.com/aws/smithy-go"
)

type listPlayer struct {
	Name string `json:"name" dynamodbav:"name"`
	// ConnID string `json:"connid" dynamodbav:"connid"`
	Ready bool `json:"ready" dynamodbav:"ready"`
}

type frontListGame struct {
	No      string       `json:"no"`
	Ready   bool         `json:"ready"`
	Players []listPlayer `json:"players"`
}

type listGamePayload struct {
	Game frontListGame
	Tag  string
}

type backListGame struct {
	Pk      string                `dynamodbav:"pk"` //'LISTGME'
	Sk      string                `dynamodbav:"sk"` //no
	Ready   bool                  `dynamodbav:"ready"`
	Players map[string]listPlayer `dynamodbav:"players"`
}

type livePlayerList []struct {
	PlayerID        string `json:"playerid"`
	Name            string `json:"name"`
	ConnID          string `json:"connid"`
	Color           string `json:"color"`
	Score           int    `json:"score"`
	Index           string `json:"index"`
	Answer          string `json:"answer"`
	HasAnswered     bool   `json:"hasAnswered"`
	PointsThisRound string `json:"pointsThisRound"`
}

type liveGame struct {
	Sk           string         `json:"sk"`
	Players      livePlayerList `json:"players"`
	CurrentWord  string         `json:"currentWord"`
	PreviousWord string         `json:"previousWord"`
	AnswersCount int            `json:"answersCount"`
	ShowAnswers  bool           `json:"showAnswers"`
	Winner       string         `json:"winner"`
}

type modifyLiveGamePayload struct {
	ModLiveGame liveGame
}

func getListPlayersSlice(pm map[string]listPlayer) (res []listPlayer) {
	res = []listPlayer{}

	for _, v := range pm {
		res = append(res, v)
	}

	return
}

func getFrontListGames(gl []backListGame) (res []frontListGame) {
	res = []frontListGame{}

	for _, g := range gl {
		pls := sortByName(getListPlayersSlice(g.Players))
		res = append(res, frontListGame{
			No:      g.Sk,
			Ready:   g.Ready,
			Players: pls,
		})
	}

	return
}

// https://go.dev/play/p/P_Z4JabiTvH
func (p listGamePayload) MarshalJSON() ([]byte, error) {
	m, err := json.Marshal(p.Game)
	if err != nil {
		return m, err
	}

	return []byte(fmt.Sprintf("{%q:%s}", p.Tag, m)), nil
}

// https://go.dev/play/p/CvniMWPoLKG
func (p modifyLiveGamePayload) MarshalJSON() ([]byte, error) {
	if p.ModLiveGame.AnswersCount == len(p.ModLiveGame.Players) {
		return []byte(`null`), nil
	}

	if p.ModLiveGame.AnswersCount > 0 {
		for i, pl := range p.ModLiveGame.Players {
			if pl.HasAnswered {
				pl.Answer = ""
				p.ModLiveGame.Players[i] = pl
			}
		}
	}

	m, err := json.Marshal(p.ModLiveGame)
	if err != nil {
		return m, err
	}

	return []byte(fmt.Sprintf("{%q:%s}", "mdLveGm", m)), nil
}

func sortByName(players []listPlayer) []listPlayer {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})

	return players
}

func (players livePlayerList) sortByAnswerThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Answer != players[j].Answer:
			return players[i].Answer < players[j].Answer
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) sortByScoreThenName() {
	sort.Slice(players, func(i, j int) bool {
		switch {
		case players[i].Score != players[j].Score:
			return players[i].Score > players[j].Score
		default:
			return players[i].Name < players[j].Name
		}
	})
}

func (players livePlayerList) getPoints() livePlayerList {
	dist := map[string]int{}

	for _, v := range players {
		dist[v.Answer]++
	}

	for i, p := range players {
		if len(p.Answer) > 1 {
			freq := dist[p.Answer]
			if freq == 2 {
				p.PointsThisRound = strconv.Itoa(3)
			} else if freq > 2 {
				p.PointsThisRound = strconv.Itoa(1)
			} else {
				p.PointsThisRound = strconv.Itoa(0)
			}
		} else {
			p.PointsThisRound = strconv.Itoa(0)

		}
		players[i] = p
	}

	return players
}

func FromDynamoDBEventAVMap(m map[string]events.DynamoDBAttributeValue) (res map[string]types.AttributeValue, err error) {
	// fmt.Println("av map: ", m)
	res = make(map[string]types.AttributeValue, len(m))

	for k, v := range m {
		res[k], err = FromDynamoDBEventAV(v)
		if err != nil {
			return nil, err
		}
	}

	return
}

func FromDynamoDBEventAVList(l []events.DynamoDBAttributeValue) (res []types.AttributeValue, err error) {
	// fmt.Println("av list: ", l)
	res = make([]types.AttributeValue, len(l))

	for i, v := range l {
		res[i], err = FromDynamoDBEventAV(v)
		if err != nil {
			return nil, err
		}
	}

	return
}

func FromDynamoDBEventAV(av events.DynamoDBAttributeValue) (types.AttributeValue, error) {
	// fmt.Println("av type: ", av, av.DataType())
	switch av.DataType() {

	case events.DataTypeBoolean: // 1
		return &types.AttributeValueMemberBOOL{Value: av.Boolean()}, nil

	case events.DataTypeList: // 3
		values, err := FromDynamoDBEventAVList(av.List())
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberL{Value: values}, nil

	case events.DataTypeMap: // 4
		values, err := FromDynamoDBEventAVMap(av.Map())
		if err != nil {
			return nil, err
		}
		return &types.AttributeValueMemberM{Value: values}, nil

	case events.DataTypeNumber: // 5
		return &types.AttributeValueMemberN{Value: av.Number()}, nil

	case events.DataTypeNull: // 7
		return &types.AttributeValueMemberNULL{Value: av.IsNull()}, nil

	case events.DataTypeString: // 8
		return &types.AttributeValueMemberS{Value: av.String()}, nil

	default:
		return nil, fmt.Errorf("unknown AttributeValue union member, %T", av)
	}
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

func handler(ctx context.Context, req events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	// fmt.Println("reqqqq", req)
	for i, rec := range req.Records {
		// fmt.Println("rekkkk", req.Records, len(req.Records))
		fmt.Println("reccc: ", i, rec, len(req.Records))

		tableName := strings.Split(rec.EventSourceArn, "/")[1]

		var ni map[string]events.DynamoDBAttributeValue

		if rec.EventName == dynamodbstreams.OperationTypeRemove {
			ni = rec.Change.OldImage
		} else {
			ni = rec.Change.NewImage
		}
		fmt.Printf("%s: %+v\n", "new db ni", ni)

		item, err := FromDynamoDBEventAVMap(ni)
		if err != nil {
			return callErr(err)
		}

		var (
			apiid    = os.Getenv("CT_APIID")
			stage    = os.Getenv("CT_STAGE")
			endpoint = "https://" + apiid + ".execute-api." + rec.AWSRegion + ".amazonaws.com/" + stage
		)

		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == apigatewaymanagementapi.ServiceID && region == rec.AWSRegion {
				ep := aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: rec.AWSRegion,
				}

				return ep, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

		apigwcfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
			config.WithEndpointResolver(customResolver),
		)
		if err != nil {
			return callErr(err)
		}

		ddbcfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(rec.AWSRegion),
			// config.WithLogger(logger),
		)
		if err != nil {
			return callErr(err)
		}

		var (
			apigwsvc = apigatewaymanagementapi.NewFromConfig(apigwcfg)
			ddbsvc   = dynamodb.NewFromConfig(ddbcfg)
			recType  = item["pk"].(*types.AttributeValueMemberS).Value
		)

		if recType == "CONNECT" {

			var connRecord struct {
				Pk, Sk, Game, Name, Color, Index, ConnID string
				Playing, Leader, Returning               bool
			}
			err = attributevalue.UnmarshalMap(item, &connRecord)
			if err != nil {
				return callErr(err)
			}

			fmt.Printf("%s%+v\n", "connrecord ", connRecord)

			var payload []byte

			if rec.EventName == dynamodbstreams.OperationTypeInsert || (rec.EventName == dynamodbstreams.OperationTypeModify && connRecord.Returning) {

				listGamesResults, err := ddbsvc.Query(ctx, &dynamodb.QueryInput{
					TableName:              aws.String(tableName),
					ScanIndexForward:       aws.Bool(false),
					KeyConditionExpression: aws.String("pk = :g"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":g": &types.AttributeValueMemberS{Value: "LISTGME"},
					},
				})
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
					Returning bool            `json:"returning"`
				}{
					ListGames: getFrontListGames(listGames),
					Name:      connRecord.Name,
					Returning: connRecord.Returning,
				})
				if err != nil {
					return callErr(err)
				}

			} else if rec.EventName == dynamodbstreams.OperationTypeModify {

				payload, err = json.Marshal(struct {
					ModConnGm string `json:"modConn"`
					Color     string `json:"color"`
					Index     string `json:"index"`
					Leader    bool   `json:"leader"`
				}{
					ModConnGm: connRecord.Game,
					Color:     connRecord.Color,
					Index:     connRecord.Index,
					Leader:    connRecord.Leader,
				})
				if err != nil {
					return callErr(err)
				}

			} else {
				oi := rec.Change.OldImage
				fmt.Printf("%s: %+v\n", "remove conn oi", oi)
				continue
			}

			conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(connRecord.ConnID), Data: payload}

			_, err = apigwsvc.PostToConnection(ctx, &conn)
			if err != nil {
				return callErr(err)
			}

		} else if recType == "LISTGME" {

			var listGameRecord backListGame
			err = attributevalue.UnmarshalMap(item, &listGameRecord)
			if err != nil {
				return callErr(err)
			}

			fmt.Printf("%s%+v\n", "list gammmmme ", listGameRecord)

			gp := listGamePayload{
				Game: frontListGame{
					No:      listGameRecord.Sk,
					Players: sortByName(getListPlayersSlice(listGameRecord.Players)),
				},
				Tag: "mdLstGm",
			}

			switch rec.EventName {
			case dynamodbstreams.OperationTypeInsert:
				gp.Tag = "addGame"
			case dynamodbstreams.OperationTypeModify:
				gp.Game.Ready = listGameRecord.Ready
			default:
				fmt.Printf("%s: %+v\n", "remove list game oi", rec.Change.OldImage)
				gp.Game.Players = nil
				gp.Tag = "rmvGame"
			}

			payload, err := json.Marshal(gp)
			if err != nil {
				return callErr(err)
			}

			connResults, err := ddbsvc.Query(ctx, &dynamodb.QueryInput{
				TableName:              aws.String(tableName),
				ScanIndexForward:       aws.Bool(true),
				KeyConditionExpression: aws.String("pk = :c"),
				FilterExpression:       aws.String("#P = :f"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":c": &types.AttributeValueMemberS{Value: "CONNECT"},
					":f": &types.AttributeValueMemberBOOL{Value: false},
				},
				ExpressionAttributeNames: map[string]string{
					"#P": "playing",
				},
			})
			if err != nil {
				return callErr(err)
			}

			var conns []struct{ ConnID string }
			err = attributevalue.UnmarshalListOfMaps(connResults.Items, &conns)
			if err != nil {
				return callErr(err)
			}

			for _, v := range conns {

				conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

				_, err := apigwsvc.PostToConnection(ctx, &conn)
				if err != nil {
					return callErr(err)
				}

			}

		} else if recType == "LIVEGME" {

			if rec.EventName == dynamodbstreams.OperationTypeInsert || rec.EventName == dynamodbstreams.OperationTypeModify {

				var gameRecord struct {
					Sk, CurrentWord, PreviousWord, Winner string
					Players                               livePlayerList
					AnswersCount                          int
					ShowAnswers                           bool
				}
				err = attributevalue.UnmarshalMap(item, &gameRecord)
				if err != nil {
					return callErr(err)
				}

				fmt.Printf("%s%+v\n", "live gammmmme ", gameRecord)

				pls := gameRecord.Players

				if gameRecord.ShowAnswers {
					pls.sortByAnswerThenName()
				} else {
					pls.sortByScoreThenName()
				}

				gp := modifyLiveGamePayload{
					ModLiveGame: liveGame{
						Sk:           gameRecord.Sk,
						Players:      pls.getPoints(),
						CurrentWord:  gameRecord.CurrentWord,
						PreviousWord: gameRecord.PreviousWord,
						AnswersCount: gameRecord.AnswersCount,
						ShowAnswers:  gameRecord.ShowAnswers,
						Winner:       gameRecord.Winner,
					},
				}

				payload, err := json.Marshal(gp)
				if err != nil {
					return callErr(err)
				}

				for _, v := range pls {

					conn := apigatewaymanagementapi.PostToConnectionInput{ConnectionId: aws.String(v.ConnID), Data: payload}

					_, err = apigwsvc.PostToConnection(ctx, &conn)
					if err != nil {
						return callErr(err)
					}

				}

			} else {
				oi := rec.Change.OldImage
				fmt.Printf("%s: %+v\n", "remove live game oi", oi)
			}

		} else {
			fmt.Println("other record type", item)
		}

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
