package main

type liveGame struct {
	No           string         `json:"no"`
	CurrentWord  string         `json:"currentWord"`
	PreviousWord string         `json:"previousWord"`
	Players      livePlayerList `json:"players"`
}

var modifyLiveGamePayload_MarshalJSON_Tests = []struct {
	input       modifyLiveGamePayload
	expected    liveGame
	description string
}{
	{input: modifyLiveGamePayload{
		ModLiveGame: toFELiveGame{
			No:           "123456789",
			CurrentWord:  "bark",
			PreviousWord: "moon",
			// Players:      livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "ans1"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "ans1"}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "ans1"}, HasAnswered: false}},
			AnswersCount: 0,
		},
	}, expected: liveGame{
		No:           "123456789",
		CurrentWord:  "bark",
		PreviousWord: "moon",
		// Players:      livePlayerList{{Name: "p1", ConnID: "111", Color: "red", Score: 10, Answer: answer{PlayerID: "p1", Answer: "ans1"}, HasAnswered: false}, {Name: "p2", ConnID: "222", Color: "green", Score: 20, Answer: answer{PlayerID: "p2", Answer: "ans1"}, HasAnswered: false}, {Name: "p3", ConnID: "333", Color: "yellow", Score: 7, Answer: answer{PlayerID: "p3", Answer: "ans1"}, HasAnswered: false}},
	}, description: "0 answers, straight marshal"},
}
