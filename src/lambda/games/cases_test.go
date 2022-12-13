package main

import "github.com/aws/aws-sdk-go-v2/aws"

var sortListPlayers = []struct {
	input, expected []listPlayer
	description     string
}{
	{
		input: []listPlayer{{Name: "bill"}, {Name: "artie"}, {Name: "wendel"}, {Name: "will"}, {Name: "mike"}},

		expected:    []listPlayer{{Name: "artie"}, {Name: "bill"}, {Name: "mike"}, {Name: "wendel"}, {Name: "will"}},
		description: "by name",
	},
}

var seven = 7
var ten = 10
var twenty = 20

var sortByAnswerThenNameTests = []struct {
	input, expected []livePlayer
	description     string
}{
	{
		input: []livePlayer{{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}},

		expected: []livePlayer{{Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}},

		description: "by answer then name",
	},
}

var sortByScoreThenNameTests = []struct {
	input, expected []livePlayer
	description     string
}{
	{
		input: []livePlayer{{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}},

		expected: []livePlayer{{Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}},

		description: "by score then name",
	},
}

var bp = []byte{123, 34, 110, 111, 34, 58, 34, 57, 57, 57, 34, 44, 34, 116, 105, 109, 101, 114, 67, 120, 108, 100, 34, 58, 102, 97, 108, 115, 101, 44, 34, 112, 108, 97, 121, 101, 114, 115, 34, 58, 91, 123, 34, 110, 97, 109, 101, 34, 58, 34, 112, 49, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 125, 44, 123, 34, 110, 97, 109, 101, 34, 58, 34, 112, 50, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 125, 44, 123, 34, 110, 97, 109, 101, 34, 58, 34, 112, 51, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 125, 93, 125, 125}

var listGamePayload_MarshalJSON_Tests = []struct {
	input       listGamePayload
	expected    []byte
	description string
}{
	{
		input:       listGamePayload{Game: frontListGame{No: "999", Players: []listPlayer{{Name: "p1"}, {Name: "p2"}, {Name: "p3"}}}, Tag: "addGame"},
		expected:    append([]byte{123, 34, 97, 100, 100, 71, 97, 109, 101, 34, 58}, bp...),
		description: "add game",
	},
	{
		input:       listGamePayload{Game: frontListGame{No: "999", Players: []listPlayer{{Name: "p1"}, {Name: "p2"}, {Name: "p3"}}}, Tag: "mdLstGm"},
		expected:    append([]byte{123, 34, 109, 100, 76, 115, 116, 71, 109, 34, 58}, bp...),
		description: "mod game",
	},
	{
		input:       listGamePayload{Game: frontListGame{No: "999", Players: []listPlayer{{Name: "p1"}, {Name: "p2"}, {Name: "p3"}}}, Tag: "rmvGame"},
		expected:    append([]byte{123, 34, 114, 109, 118, 71, 97, 109, 101, 34, 58}, bp...),
		description: "remove game",
	},
}

var prepTests = []struct {
	input, expected []livePlayer
	description     string
}{
	{
		input: []livePlayer{
			{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: true},
			{Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: true},
			{Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: true},
			{Name: "darlene", ConnID: "333", Color: "yellow", Score: &ten, Answer: "meal", HasAnswered: true},
			{Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: true},
			{Name: "william", ConnID: "333", Color: "yellow", Score: &twenty, Answer: "meal", HasAnswered: true},
			{Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: true},
		},

		expected: []livePlayer{
			{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false, PointsThisRound: aws.Int(0)},
			{Name: "darlene", ConnID: "333", Color: "yellow", Score: &ten, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "william", ConnID: "333", Color: "yellow", Score: &twenty, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
		},

		description: "prep live players",
	},
}

var showAnswersTests = []struct {
	input, expected []livePlayer
	description     string
}{
	{
		input: []livePlayer{
			{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false, PointsThisRound: aws.Int(0)},
			{Name: "darlene", ConnID: "333", Color: "yellow", Score: &ten, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "william", ConnID: "333", Color: "yellow", Score: &twenty, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
		},

		expected: []livePlayer{
			{Name: "will", ConnID: "333", Color: "yellow", Score: nil, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "earl", ConnID: "111", Color: "red", Score: nil, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "carl", ConnID: "222", Color: "green", Score: nil, Answer: "verb", HasAnswered: false, PointsThisRound: aws.Int(0)},
			{Name: "darlene", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "dean", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "william", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "beulah", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
		},

		description: "nilificate scores",
	},
}

var clearAnswersTests = []struct {
	input, expected []livePlayer
	description     string
}{
	{
		input: []livePlayer{
			{Name: "will", ConnID: "333", Color: "yellow", Score: nil, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "earl", ConnID: "111", Color: "red", Score: nil, Answer: "heart", HasAnswered: false, PointsThisRound: aws.Int(3)},
			{Name: "carl", ConnID: "222", Color: "green", Score: nil, Answer: "verb", HasAnswered: false, PointsThisRound: aws.Int(0)},
			{Name: "darlene", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "dean", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "william", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
			{Name: "beulah", ConnID: "333", Color: "yellow", Score: nil, Answer: "meal", HasAnswered: false, PointsThisRound: aws.Int(1)},
		},

		expected: []livePlayer{
			{Name: "will", ConnID: "333", Color: "yellow", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(3)},
			{Name: "earl", ConnID: "111", Color: "red", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(3)},
			{Name: "carl", ConnID: "222", Color: "green", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(0)},
			{Name: "darlene", ConnID: "333", Color: "yellow", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(1)},
			{Name: "dean", ConnID: "333", Color: "yellow", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(1)},
			{Name: "william", ConnID: "333", Color: "yellow", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(1)},
			{Name: "beulah", ConnID: "333", Color: "yellow", Score: nil, Answer: "", HasAnswered: true, PointsThisRound: aws.Int(1)},
		},

		description: "delete answers",
	},
}
