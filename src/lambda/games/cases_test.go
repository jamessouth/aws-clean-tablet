package main

type liveGameWrap struct {
	MdLveGm liveGame `json:"mdLveGm"`
}

var modifyLiveGamePayload_MarshalJSON_Tests = []struct {
	input       modifyLiveGamePayload
	expected    liveGameWrap
	description string
}{
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "", HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "", HasAnswered: false}}, AnswersCount: 0}},
		expected:    liveGameWrap{MdLveGm: liveGame{Sk: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "", HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "", HasAnswered: false}}}},
		description: "answer count 0/3, straight marshal",
	},
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "12389", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "", HasAnswered: false}}, AnswersCount: 1}},
		expected:    liveGameWrap{MdLveGm: liveGame{Sk: "12389", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "", HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "", HasAnswered: false}}, AnswersCount: 1}},
		description: "answer count 1/3",
	},
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "16789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: true}}, AnswersCount: 2}},
		expected:    liveGameWrap{MdLveGm: liveGame{Sk: "16789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "", HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "", HasAnswered: true}}, AnswersCount: 2}},
		description: "answer count 2/3",
	},
}

var modifyLiveGamePayload_MarshalJSON_Tests2 = []struct {
	input       modifyLiveGamePayload
	description string
}{
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "1244444489", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: "verb", HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}}, AnswersCount: 3}},
		description: "answer count 3/3, return null",
	},
}

// ---------------------------------------------------------------------
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

var sortByAnswerThenName = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input: livePlayerList{{Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: "verb", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}},

		expected: livePlayerList{{Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: "verb", HasAnswered: false}},

		description: "by answer then name",
	},
}

var sortByScoreThenName = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input: livePlayerList{{Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: "verb", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}},

		expected: livePlayerList{{Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: "verb", HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: "heart", HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: "heart", HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: "meal", HasAnswered: false}},

		description: "by score then name",
	},
}

var bp = []byte{123, 34, 110, 111, 34, 58, 34, 57, 57, 57, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 44, 34, 112, 108, 97, 121, 101, 114, 115, 34, 58, 91, 123, 34, 110, 97, 109, 101, 34, 58, 34, 112, 49, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 125, 44, 123, 34, 110, 97, 109, 101, 34, 58, 34, 112, 50, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 125, 44, 123, 34, 110, 97, 109, 101, 34, 58, 34, 112, 51, 34, 44, 34, 114, 101, 97, 100, 121, 34, 58, 102, 97, 108, 115, 101, 125, 93, 125, 125}

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
