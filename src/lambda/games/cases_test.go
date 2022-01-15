package main

// type liveGame struct {
// 	No           string         `json:"no"`
// 	CurrentWord  string         `json:"currentWord"`
// 	PreviousWord string         `json:"previousWord"`
// 	Players      livePlayerList `json:"players"`
// }

type liveGameWrap struct {
	MdLveGm liveGame `json:"mdLveGm"`
}

var modifyLiveGamePayload_MarshalJSON_Tests = []struct {
	input       modifyLiveGamePayload
	expected    liveGameWrap
	description string
}{
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}, AnswersCount: 0}},
		expected:    liveGameWrap{MdLveGm: liveGame{Sk: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}}},
		description: "answer count 0/3, straight marshal",
	},
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "12389", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}, AnswersCount: 1}},
		expected:    liveGameWrap{MdLveGm: liveGame{Sk: "12389", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}, AnswersCount: 1}},
		description: "answer count 1/3",
	},
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "16789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}}, AnswersCount: 2}},
		expected:    liveGameWrap{MdLveGm: liveGame{Sk: "16789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: true}}, AnswersCount: 2}},
		description: "answer count 2/3",
	},
}

var modifyLiveGamePayload_MarshalJSON_Tests2 = []struct {
	input       modifyLiveGamePayload
	description string
}{
	{
		input:       modifyLiveGamePayload{ModLiveGame: liveGame{Sk: "1244444489", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}}, AnswersCount: 3}},
		description: "answer count 3/3, return null",
	},
}

// ---------------------------------------------------------------------
var sortListPlayers = []struct {
	input, expected listPlayerList
	description     string
}{
	{
		input: listPlayerList{{Name: "bill", ConnID: "111", Ready: false}, {Name: "artie", ConnID: "222", Ready: false}, {Name: "wendel", ConnID: "333", Ready: false}, {Name: "will", ConnID: "111", Ready: false}, {Name: "mike", ConnID: "111", Ready: false}},

		expected:    listPlayerList{{Name: "artie", ConnID: "222", Ready: false}, {Name: "bill", ConnID: "111", Ready: false}, {Name: "mike", ConnID: "111", Ready: false}, {Name: "wendel", ConnID: "333", Ready: false}, {Name: "will", ConnID: "111", Ready: false}},
		description: "by name",
	},
}

var sortLivePlayers = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input:       livePlayerList{{Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "carl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "will", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}},
		expected:    livePlayerList{{Name: "carl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "will", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}},
		description: "by name",
	},
}

var sortByAnswerThenName = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input: livePlayerList{{Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}},

		expected: livePlayerList{{Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}},

		description: "by answer then name",
	},
}

var sortByScoreThenName = []struct {
	input, expected livePlayerList
	description     string
}{
	{
		input: livePlayerList{{Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}},

		expected: livePlayerList{{Name: "carl", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}, {Name: "earl", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "beulah", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}, {Name: "darlene", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "dean", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}, {Name: "will", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "heart"}, HasAnswered: false}, {Name: "william", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}},

		description: "by score then name",
	},
}
