package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	slope     int    = 0
	intercept int    = 2
	connect   string = "CONNECT"
	listGame  string = "LISTGAME"
	liveGame  string = "LIVEGAME"
	sentinel  string = "game over"
)

// uncomment for test
// type ctxKey string

type DdbBatchWriteItemAPI interface {
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}

type DdbDeleteItemAPI interface {
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type listPlayer struct {
	Name   string `dynamodbav:"name"`
	ConnID string `dynamodbav:"connid"`
}

type livePlayer struct {
	Name   string `json:"name"`
	ConnID string `json:"connid"`
	Color  string `json:"color"`
	Answer string `json:"answer"`
	Score  *int   `json:"score"`
}

type output struct {
	Gameno string `json:"gameno"`
}

type stringSlice []string

func (list stringSlice) shuffleList(length int) stringSlice {
	t := time.Now().UnixNano()
	rand.Seed(t)

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	return list[:length]
}

func getLivePlayerMap(pm map[string]listPlayer, colors stringSlice) (plrs map[string]livePlayer) {
	plrs = map[string]livePlayer{}
	count := 0

	for k, v := range pm {
		plrs[k] = livePlayer{
			Name:   v.Name,
			ConnID: v.ConnID,
			Color:  colors[count],
			Answer: "",
			Score:  aws.Int(0),
		}
		count++
	}

	return
}

func deleteItem(ctx context.Context, api DdbDeleteItemAPI, pk, sk, tableName string) (dynamodb.DeleteItemOutput, error) {
	di, err := api.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		TableName:    aws.String(tableName),
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return dynamodb.DeleteItemOutput{}, err
	}

	if res := *di; len(res.Attributes) == 0 {
		return dynamodb.DeleteItemOutput{}, fmt.Errorf("error: item with pk %s and sk %s not found", pk, sk)
	} else {
		return res, nil
	}
}

func batchWriteItem(ctx context.Context, api DdbBatchWriteItemAPI, items []types.WriteRequest, tableName string) (dynamodb.BatchWriteItemOutput, error) {
	bwo, err := api.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems:                map[string][]types.WriteRequest{tableName: items},
		ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
		ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
	})
	if err != nil {
		return dynamodb.BatchWriteItemOutput{}, err
	}

	return *bwo, nil
}

func handleUnprocessedItems(ctx context.Context, api DdbBatchWriteItemAPI, batchWriteOutput dynamodb.BatchWriteItemOutput, tableName string) error {
	var (
		ret, maxDelay, factor time.Duration = 500, 5000, 2
		unprocessedItems                    = batchWriteOutput.UnprocessedItems
	)

	for len(unprocessedItems[tableName]) > 0 && ret < maxDelay {
		// uncomment to test
		// var mk ctxKey = "cKey"
		// ctx = context.WithValue(ctx, mk, ret.Nanoseconds())
		time.Sleep(ret * time.Millisecond)
		batchWriteOutput2, err := batchWriteItem(ctx, api, unprocessedItems[tableName], tableName)
		if err != nil {
			return err
		}

		unprocessedItems = batchWriteOutput2.UnprocessedItems

		ret *= factor
	}
	if len(unprocessedItems[tableName]) > 0 {
		return fmt.Errorf("error: unable to write %d items", len(unprocessedItems[tableName]))
	}

	return nil
}

func handler(ctx context.Context, req struct {
	Payload struct {
		Gameno, TableName, Region string
	}
}) (output, error) {

	fmt.Printf("%s%+v\n", "prep req ", req)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(req.Payload.Region),
	)
	if err != nil {
		return output{}, err
	}

	var (
		ddbsvc    = dynamodb.NewFromConfig(cfg)
		gameno    = req.Payload.Gameno
		tableName = req.Payload.TableName
	)

	di, err := deleteItem(ctx, ddbsvc, listGame, gameno, tableName)
	if err != nil {
		return output{}, err
	}

	var game struct {
		Players map[string]listPlayer
	}
	err = attributevalue.UnmarshalMap(di.Attributes, &game)
	if err != nil {
		return output{}, err
	}

	players := getLivePlayerMap(game.Players, colors.shuffleList(len(game.Players)))

	marshalledPlayers, err := attributevalue.Marshal(players)
	if err != nil {
		return output{}, err
	}

	wordList := append(words.shuffleList(slope*len(players)+intercept), sentinel)

	marshalledWordList, err := attributevalue.Marshal(wordList)
	if err != nil {
		return output{}, err
	}

	for k := range players {
		_, err := ddbsvc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: connect},
				"sk": &types.AttributeValueMemberS{Value: k},
			},
			TableName: aws.String(tableName),
			ExpressionAttributeNames: map[string]string{
				"#P": "playing",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":p": &types.AttributeValueMemberBOOL{Value: true},
			},
			UpdateExpression: aws.String("set #P = :p"),
		})
		if err != nil {
			return output{}, err
		}
	}

	emptyPlayersMap, err := attributevalue.Marshal(map[string]listPlayer{})
	if err != nil {
		return output{}, err
	}

	items := []types.WriteRequest{
		{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"pk":        &types.AttributeValueMemberS{Value: listGame},
					"sk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", time.Now().UnixNano())},
					"players":   emptyPlayersMap,
					"timerCxld": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
		},
		{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"pk":           &types.AttributeValueMemberS{Value: liveGame},
					"sk":           &types.AttributeValueMemberS{Value: gameno},
					"answersCount": &types.AttributeValueMemberN{Value: "0"},
					"players":      marshalledPlayers,
					"wordList":     marshalledWordList,
				},
			},
		},
	}

	unwrittenItems, err := batchWriteItem(ctx, ddbsvc, items, tableName)
	if err != nil {
		return output{}, err
	}

	err = handleUnprocessedItems(ctx, ddbsvc, unwrittenItems, tableName)
	if err != nil {
		return output{}, err
	}

	return output{Gameno: gameno}, nil
}

func main() {
	lambda.Start(handler)
}

var colors = stringSlice{
	"#007d5e",
	"#00591c",
	"#004528",
	"#064e3b",
	"#003200",
	"#365314",
	"#2b704b",
	"#00596d",
	"#2f4858",
	"#344a71",
	"#000058",
	"#1e3a8a",
	"#192a84",
	"#3d5988",
	"#423040",
	"#4b4737",
	"#5a447e",
	"#650001",
	"#701a75",
	"#893273",
	"#564516",
	"#8d814d",
	"#78350f",
	"#872c1a",
	"#b91c1c",
	"#ae652c",
	"#894f5e",
	"#cb4956",
	"#990046",
	"#ad1351",
	"#77005e",
	"#b2286e",
}

var words = stringSlice{
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
