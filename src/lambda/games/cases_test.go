package main

var sortListPlayers = []struct {
	input, expected []listPlayer
	description     string
}{
	{
		input: []listPlayer{{Name: "bill", Ready: false}, {Name: "artie", Ready: false}, {Name: "wendel", Ready: false}, {Name: "will", Ready: false}, {Name: "mike", Ready: false}},

		expected:    []listPlayer{{Name: "artie", Ready: false}, {Name: "bill", Ready: false}, {Name: "mike", Ready: false}, {Name: "wendel", Ready: false}, {Name: "will", Ready: false}},
		description: "by name",
	},
}

var seven = 7
var ten = 10
var twenty = 20

var sortByAnswerThenName = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input: livePlayerList{{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}},

		expected: livePlayerList{{Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}},

		description: "by answer then name",
	},
}

var sortByScoreThenName = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input: livePlayerList{{Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}},

		expected: livePlayerList{{Name: "carl", ConnID: "222", Color: "green", Score: &twenty, Answer: "verb", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: &ten, Answer: "heart", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: &seven, Answer: "heart", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: &seven, Answer: "meal", HasAnswered: false}},

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
