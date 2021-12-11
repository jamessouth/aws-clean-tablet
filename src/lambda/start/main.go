package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
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
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"

	"github.com/aws/smithy-go"
)

var colors = []string{
	"#dc2626", //red 600
	"#0c4a6e", //light blue 900
	"#16a34a", //green 600
	"#7c2d12", //orange 900
	"#c026d3", //fuchsia 600
	"#365314", //lime 900
	"#0891b2", //cyan 600
	"#581c87", //purple 900
}

// type key struct {
// 	Pk string `dynamodbav:"pk"`
// 	Sk string `dynamodbav:"sk"`
// }

// type player struct {
// 	Name   string `dynamodbav:"name"`
// 	ConnID string `dynamodbav:"connid"`
// 	Ready  bool   `dynamodbav:"ready"`
// 	Color  string `dynamodbav:"color,omitempty"`
// 	Score  int    `dynamodbav:"score"`
// 	Answer answer `dynamodbav:"answer"`
// }

type answer struct {
	PlayerID, Answer string
}

// type hiScore struct {
// 	Score int  `json:"score"`
// 	Tie   bool `json:"tie"`
// }

type livePlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
	Color  string `dynamodbav:"color"`
	Score  int    `dynamodbav:"score"`
	Answer answer `dynamodbav:"answer"`
}

type livePlayerMap map[string]livePlayer

type liveGame struct {
	// Pk           string        `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
	// CurrentWord  string        `dynamodbav:"currentWord"`
	Players livePlayerMap `dynamodbav:"players"`
	// AnswersCount int           `dynamodbav:"answersCount"`
	// SendToFront  bool          `dynamodbav:"sendToFront"`
}

// type body struct {
// 	Gameno string
// }

type sfnArrInput struct {
	Id   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

type sfnInput struct {
	Gameno  string        `dynamodbav:"gameno"`
	Players []sfnArrInput `dynamodbav:"players"`
}

func (pm livePlayerMap) assignColors() livePlayerMap {
	count := 0
	for k, v := range pm {
		v.Color = colors[count]
		pm[k] = v
		count++
	}

	return pm
}

func (pm livePlayerMap) mapToSlice() (res []sfnArrInput) {
	for k, v := range pm {
		res = append(res, sfnArrInput{
			Id:   k,
			Name: v.Name,
		})
	}

	return
}

func handler(ctx context.Context, req events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("start", req.Body)

	reg := strings.Split(req.RequestContext.DomainName, ".")[2]

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(reg),
	)
	if err != nil {
		return callErr(err)
	}

	tableName, ok := os.LookupEnv("tableName")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find table name"))
	}

	sfnarn, ok := os.LookupEnv("SFNARN")
	if !ok {
		panic(fmt.Sprintf("%v", "can't find sfn arn"))
	}

	ddbsvc := dynamodb.NewFromConfig(cfg)
	sfnsvc := sfn.NewFromConfig(cfg)

	var gameno string

	err = json.Unmarshal([]byte(req.Body), &gameno)
	if err != nil {
		return callErr(err)
	}

	di, err := ddbsvc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "LISTGME"},
			"sk": &types.AttributeValueMemberS{Value: gameno},
		},
		TableName:    aws.String(tableName),
		ReturnValues: types.ReturnValueAllOld,
	})
	callErr(err)

	var game liveGame
	err = attributevalue.UnmarshalMap(di.Attributes, &game)
	if err != nil {
		return callErr(err)
	}

	fmt.Printf("%s%+v\n", "livegame ", game)

	const numberOfWords int = 40

	marshalledWordsList, err := attributevalue.Marshal(shuffleList(words, numberOfWords))
	if err != nil {
		return callErr(err)
	}

	playersMap := game.Players.assignColors()

	marshalledPlayersMap, err := attributevalue.Marshal(playersMap)
	if err != nil {
		return callErr(err)
	}

	_, err = ddbsvc.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"pk":       &types.AttributeValueMemberS{Value: "LIVEGME"},
			"sk":       &types.AttributeValueMemberS{Value: game.Sk},
			"players":  marshalledPlayersMap,
			"wordList": marshalledWordsList,
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		return callErr(err)
	}

	sfnInput, err := json.Marshal(sfnInput{
		Gameno:  gameno,
		Players: playersMap.mapToSlice(),
	})
	if err != nil {
		return callErr(err)
	}

	ssei := sfn.StartSyncExecutionInput{
		StateMachineArn: aws.String(sfnarn),
		Input:           aws.String(string(sfnInput)),
	}

	sse, err := sfnsvc.StartSyncExecution(ctx, &ssei)
	if err != nil {
		return callErr(err)
	}

	sseo := *sse
	fmt.Printf("\n%s, %+v\n", "sse op", sseo)

	if sseo.Status == sfntypes.SyncExecutionStatusFailed || sseo.Status == sfntypes.SyncExecutionStatusTimedOut {
		err := fmt.Errorf("step function %s, execution %s, failed with status %s. error code: %s. cause: %s. ", *sseo.StateMachineArn, *sseo.ExecutionArn, sseo.Status, *sseo.Error, *sseo.Cause)
		return callErr(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, nil
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

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusBadRequest,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{},
		Body:              "",
		IsBase64Encoded:   false,
	}, err

}

func shuffleList(words []string, length int) []string {
	t := time.Now().UnixNano()
	rand.Seed(t)

	nl := append([]string(nil), words...)

	rand.Shuffle(len(nl), func(i, j int) {
		nl[i], nl[j] = nl[j], nl[i]
	})

	return nl[:length]
}

var words = []string{
	"half ____",
	"____ child",
	"middle ____",
	"____ wash",
	"car ____",
	"tooth ____",
	"time ____",
	"running ____",
	"party ____",
	"social ____",
	"night ____",
	"____ gear",
	"____ dollar",
	"chop ____",
	"milk ____",
	"____ water",
	"south ____",
	"pillow ____",
	"head ____",
	"____ powder",
	"happy ____",
	"____ potato",
	"____ storm",
	"lime ____",
	"roller ____",
	"____ language",
	"world ____",
	"evening ____",
	"____ service",
	"shopping ____",
	"shoe ____",
	"wet ____",
	"____ cow",
	"hold ____",
	"____ finger",
	"mouth ____",
	"____ course",
	"chain ____",
	"____ doll",
	"just ____",
	"under ____",
	"sea ____",
	"tropical ____",
	"chicken ____",
	"____ land",
	"rock ____",
	"tree ____",
	"cherry ____",
	"wide ____",
	"short ____",
	"pay ____",
	"____ grown",
	"____ bench",
	"food ____",
	"training ____",
	"screw ____",
	"____ bread",
	"body ____",
	"dinner ____",
	"____ fish",
	"grand ____",
	"____ paper",
	"____ keeper",
	"fine ____",
	"____ juice",
	"field ____",
	"belly ____",
	"____ berry",
	"salad ____",
	"____ limit",
	"bed ____",
	"____ belly",
	"mass ____",
	"____ tank",
	"flat ____",
	"____ mark",
	"good ____",
	"____ pad",
	"bull ____",
	"____ class",
	"lucky ____",
	"____ story",
	"front ____",
	"____ bite",
	"center ____",
	"____ corn",
	"chocolate ____",
	"____ dog",
	"baby ____",
	"____ aid",
	"grape ____",
	"____ party",
	"hang ____",
	"____ pole",
	"off ____",
	"____ glove",
	"bus ____",
	"thank ____",
	"parking ____",
	"____ green",
	"air ____",
	"snake ____",
	"finger ____",
	"____ bean",
	"golf ____",
	"____ body",
	"stock ____",
	"pot ____",
	"super ____",
	"real ____",
	"____ name",
	"fresh ____",
	"____ weight",
	"mixed ____",
	"tennis ____",
	"black ____",
	"____ mate",
	"for ____",
	"top ____",
	"sand ____",
	"____ blue",
	"gas ____",
	"____ duty",
	"lip ____",
	"____ cup",
	"horse ____",
	"____ ever",
	"crab ____",
	"white ____",
	"drive ____",
	"____ fry",
	"nice ____",
	"____ cake",
	"box ____",
	"____ bag",
	"baked ____",
	"____ stop",
	"light ____",
	"____ star",
	"land ____",
	"____ station",
	"left ____",
	"summer ____",
	"better ____",
	"____ free",
	"double ____",
	"____ blanket",
	"fruit ____",
	"____ night",
	"game ____",
	"____ salad",
	"hot ____",
	"____ ring",
	"home ____",
	"speed ____",
	"base ____",
	"____ shop",
	"jelly ____",
	"____ gun",
	"draw ____",
	"too ____",
	"candle ____",
	"silver ____",
	"no ____",
	"____ hand",
	"pig ____",
	"sour ____",
	"perfect ____",
	"sitting ____",
	"north ____",
	"____ bell",
	"flower ____",
	"____ fire",
	"moving ____",
	"____ front",
	"never ____",
	"third ____",
	"business ____",
	"____ egg",
	"cow ____",
	"____ ticket",
	"master ____",
	"____ food",
	"door ____",
	"____ cycle",
	"hyper ____",
	"____ drop",
	"cold ____",
	"____ order",
	"golden ____",
	"____ oil",
	"go ____",
	"sub ____",
	"punch ____",
	"spit ____",
	"beach ____",
	"____ pit",
	"school ____",
	"____ bug",
	"bottom ____",
	"strip ____",
	"prime ____",
	"smooth ____",
	"out ____",
	"____ well",
	"mini ____",
	"sweat ____",
	"big ____",
	"side ____",
	"mud ____",
	"team ____",
	"bowling ____",
	"sound ____",
	"back ____",
	"____ case",
	"brass ____",
	"____ club",
	"health ____",
	"vegetable ____",
	"country ____",
	"____ frame",
	"name ____",
	"____ pie",
	"guess ____",
	"____ print",
	"heavy ____",
	"____ shower",
	"jungle ____",
	"____ house",
	"false ____",
	"____ block",
	"garbage ____",
	"____ book",
	"board ____",
	"____ fee",
	"main ____",
	"spoiled ____",
	"pine ____",
	"____ meat",
	"free ____",
	"truck ____",
	"christmas ____",
	"tax ____",
	"birth ____",
	"wild ____",
	"shot ____",
	"____ break",
	"green ____",
	"____ floor",
	"long ____",
	"sun ____",
	"raw ____",
	"sure ____",
	"rest ____",
	"____ dance",
	"jail ____",
	"____ wine",
	"monkey ____",
	"____ less",
	"safety ____",
	"tail ____",
	"root ____",
	"____ age",
	"ever ____",
	"____ pot",
	"hard ____",
	"sweet ____",
	"right ____",
	"____ father",
	"day ____",
	"____ walk",
	"mid ____",
	"soft ____",
	"american ____",
	"tough ____",
	"cash ____",
	"small ____",
	"open ____",
	"____ guard",
	"pea ____",
	"____ work",
	"motor ____",
	"____ good",
	"oil ____",
	"____ sauce",
	"human ____",
	"____ coat",
	"heart ____",
	"____ driver",
	"coffee ____",
	"____ horse",
	"elbow ____",
	"____ guess",
	"pepper ____",
	"____ beer",
	"bath ____",
	"____ time",
	"meat ____",
	"____ hour",
	"red ____",
	"welcome ____",
	"dirt ____",
	"____ spot",
	"single ____",
	"star ____",
	"pocket ____",
	"____ duck",
	"court ____",
	"____ suit",
	"magic ____",
	"____ pick",
	"growing ____",
	"____ field",
	"mother ____",
	"____ stool",
	"life ____",
	"____ skate",
	"junk ____",
	"upper ____",
	"club ____",
	"____ rest",
	"holy ____",
	"____ business",
	"hair ____",
	"____ burn",
	"guest ____",
	"____ bird",
	"fried ____",
	"____ paint",
	"same ____",
	"____ bone",
	"birthday ____",
	"____ life",
	"fire ____",
	"____ friend",
	"down ____",
	"so ____",
	"paper ____",
	"soul ____",
	"penny ____",
	"____ done",
	"key ____",
	"tea ____",
	"round ____",
	"test ____",
	"blue ____",
	"toilet ____",
	"busy ____",
	"____ bed",
	"fish ____",
	"____ hard",
	"rain ____",
	"____ key",
	"rice ____",
	"____ word",
	"moon ____",
	"____ band",
	"eye ____",
	"____ bar",
	"face ____",
	"____ fight",
	"deep ____",
	"spring ____",
	"play ____",
	"train ____",
	"cheese ____",
	"____ jam",
	"fast ____",
	"____ court",
	"cheap ____",
	"____ bee",
	"barn ____",
	"____ fly",
	"dog ____",
	"____ chip",
	"broken ____",
	"____ neck",
	"full ____",
	"silent ____",
	"neck ____",
	"straight ____",
	"bean ____",
	"slow ____",
	"oh ____",
	"____ job",
	"fat ____",
	"____ fair",
	"love ____",
	"rubber ____",
	"tight ____",
	"bubble ____",
	"____ chocolate",
	"honey ____",
	"____ room",
	"pretty ____",
	"string ____",
	"salt ____",
	"____ load",
	"pin ____",
	"spare ____",
	"second ____",
	"water ____",
	"book ____",
	"____ brush",
	"deadly ____",
	"ball ____",
	"leading ____",
	"____ drum",
	"micro ____",
	"____ town",
	"nose ____",
	"search ____",
	"jet ____",
	"flash ____",
	"best ____",
	"odd ____",
	"picnic ____",
	"french ____",
	"great ____",
	"cat ____",
	"banana ____",
	"pop ____",
	"dirty ____",
	"kick ____",
	"even ____",
	"lunch ____",
	"____ flow",
	"wine ____",
	"____ down",
	"____ flakes",
	"stiff ____",
	"____ basket",
	"traffic ____",
	"____ bowl",
	"____ mouth",
	"____ guy",
	"____ glass",
	"____ boat",
	"____ luck",
	"____ shot",
	"up ____",
	"sky ____",
	"ice ____",
	"make ____",
	"candy ____",
	"easter ____",
	"apple ____",
	"semi ____",
	"man ____",
	"pit ____",
	"bare ____",
	"jack ____",
	"ground ____",
	"wedding ____",
	"dead ____",
	"high ____",
	"keep ____",
	"security ____",
	"jump ____",
	"gift ____",
	"hand ____",
	"first ____",
	"cream ____",
	"over ____",
	"get ____",
	"house ____",
	"lap ____",
	"mountain ____",
	"egg ____",
	"check ____",
	"foot ____",
	"____ market",
	"____ cream",
	"window ____",
	"show ____",
	"____ drive",
	"____ cut",
	"____ office",
	"snow ____",
	"____ face",
	"____ light",
	"____ chance",
	"____ board",
	"____ date",
	"what's ____",
	"____ door",
	"____ clock",
	"____ feet",
	"____ ball",
	"____ pen",
	"____ shrine",
	"____ bear",
	"spot ____",
	"____ tag",
	"____ power",
	"____ ache",
	"____ hole",
	"____ control",
	"____ table",
	"____ seat",
}
