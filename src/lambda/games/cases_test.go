package main

type liveGame struct {
	No           string         `json:"no"`
	CurrentWord  string         `json:"currentWord"`
	PreviousWord string         `json:"previousWord"`
	Players      livePlayerList `json:"players"`
}

type liveGameWrap struct {
	MdLveGm liveGame `json:"mdLveGm"`
}

var modifyLiveGamePayload_MarshalJSON_Tests = []struct {
	input       modifyLiveGamePayload
	expected    liveGameWrap
	description string
}{
	{
		input:       modifyLiveGamePayload{ModLiveGame: toFELiveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}, AnswersCount: 0}},
		expected:    liveGameWrap{MdLveGm: liveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}}},
		description: "answer count 0/3, straight marshal",
	},
	{
		input:       modifyLiveGamePayload{ModLiveGame: toFELiveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}, AnswersCount: 1}},
		expected:    liveGameWrap{MdLveGm: liveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: false}}}},
		description: "answer count 1/3",
	},
	{
		input:       modifyLiveGamePayload{ModLiveGame: toFELiveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}}, AnswersCount: 2}},
		expected:    liveGameWrap{MdLveGm: liveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: ""}, HasAnswered: true}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: ""}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: ""}, HasAnswered: true}}}},
		description: "answer count 2/3",
	},
}

var modifyLiveGamePayload_MarshalJSON_Tests2 = []struct {
	input       modifyLiveGamePayload
	description string
}{
	{
		input:       modifyLiveGamePayload{ModLiveGame: toFELiveGame{No: "123456789", CurrentWord: "bark", PreviousWord: "moon", Players: livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "heart"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "verb"}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "meal"}, HasAnswered: false}}, AnswersCount: 3}},
		description: "answer count 3/3, return null",
	},
}
