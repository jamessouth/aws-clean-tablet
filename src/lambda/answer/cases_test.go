package main

// type desc struct {
// 	ZeroScores, OneScores, ThreeScores int
// 	Description                        string
// }

// var updateScoresTests = []struct {
// 	input       liveGame
// 	expected    liveGame
// 	description desc
// }{
// 	{input: liveGame{
// 		Pk:          "",
// 		Sk:          "",
// 		CurrentWord: "",
// 		Players: map[string]livePlayer{
// 			"p1": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  10,
// 				Answer: answer{
// 					PlayerID: "111",
// 					Answer:   "ans1",
// 				},
// 			},
// 			"p2": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  20,
// 				Answer: answer{
// 					PlayerID: "222",
// 					Answer:   "ans1",
// 				},
// 			},
// 			"p3": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  7,
// 				Answer: answer{
// 					PlayerID: "333",
// 					Answer:   "ans1",
// 				},
// 			},
// 			"p4": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  11,
// 				Answer: answer{
// 					PlayerID: "444",
// 					Answer:   "ans2",
// 				},
// 			},
// 			"p5": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  15,
// 				Answer: answer{
// 					PlayerID: "555",
// 					Answer:   "ans2",
// 				},
// 			},
// 		},
// 		AnswersCount: 0,
// 		SendToFront:  false,
// 		HiScore:      20,
// 		GameTied:     false,
// 	}, expected: liveGame{
// 		Pk:          "",
// 		Sk:          "",
// 		CurrentWord: "",
// 		Players: map[string]livePlayer{
// 			"p1": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  11,
// 				Answer: answer{
// 					PlayerID: "111",
// 					Answer:   "ans1",
// 				},
// 			},
// 			"p2": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  21,
// 				Answer: answer{
// 					PlayerID: "222",
// 					Answer:   "ans1",
// 				},
// 			},
// 			"p3": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  8,
// 				Answer: answer{
// 					PlayerID: "333",
// 					Answer:   "ans1",
// 				},
// 			},
// 			"p4": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  14,
// 				Answer: answer{
// 					PlayerID: "444",
// 					Answer:   "ans2",
// 				},
// 			},
// 			"p5": {
// 				Name:   "",
// 				ConnID: "",
// 				Color:  "",
// 				Score:  18,
// 				Answer: answer{
// 					PlayerID: "555",
// 					Answer:   "ans2",
// 				},
// 			},
// 		},
// 		AnswersCount: 0,
// 		SendToFront:  false,
// 		HiScore:      21,
// 		GameTied:     false,
// 	}, description: desc{
// 		ZeroScores:  0,
// 		OneScores:   3,
// 		ThreeScores: 2,
// 		Description: "no hiscore ties",
// 	}},
// }

type chsInput struct {
	score, hiScore int
	tied           bool
}

type chsExpected struct {
	HiScore int
	Tied    bool
}

var checkHiScoreTests = []struct {
	input       chsInput
	expected    chsExpected
	description string
}{
	{
		input: chsInput{
			score:   20,
			hiScore: 20,
			tied:    false,
		}, expected: chsExpected{
			HiScore: 20,
			Tied:    true,
		}, description: "ties with high score"},
	{
		input: chsInput{
			score:   20,
			hiScore: 21,
			tied:    false,
		}, expected: chsExpected{
			HiScore: 21,
			Tied:    false,
		}, description: "no ties, no high score"},
	{
		input: chsInput{
			score:   21,
			hiScore: 20,
			tied:    true,
		}, expected: chsExpected{
			HiScore: 21,
			Tied:    false,
		}, description: "break tie with high score"},
}

type adjScInput struct {
	old           livePlayer
	incr, hiScore int
	tied          bool
}

type adjScExpected struct {
	Lp      livePlayer
	HiScore int
	Tied    bool
}

var adjScoreTests = []struct {
	input       adjScInput
	expected    adjScExpected
	description string
}{
	{input: adjScInput{
		old: livePlayer{
			Name:   "",
			ConnID: "",
			Color:  "",
			Score:  18,
			Answer: answer{
				PlayerID: "555",
				Answer:   "ans2",
			},
		},
		incr:    1,
		hiScore: 18,
		tied:    false,
	}, expected: adjScExpected{
		Lp: livePlayer{
			Name:   "",
			ConnID: "",
			Color:  "",
			Score:  19,
			Answer: answer{
				PlayerID: "555",
				Answer:   "ans2",
			},
		},
		HiScore: 19,
		Tied:    false,
	}, description: "incr by 1 incr hi"},
	{input: adjScInput{
		old: livePlayer{
			Name:   "",
			ConnID: "",
			Color:  "",
			Score:  18,
			Answer: answer{
				PlayerID: "555",
				Answer:   "ans2",
			},
		},
		incr:    3,
		hiScore: 24,
		tied:    true,
	}, expected: adjScExpected{
		Lp: livePlayer{
			Name:   "",
			ConnID: "",
			Color:  "",
			Score:  21,
			Answer: answer{
				PlayerID: "555",
				Answer:   "ans2",
			},
		},
		HiScore: 24,
		Tied:    true,
	}, description: "incr by 3 tie"},
	{input: adjScInput{
		old: livePlayer{
			Name:   "",
			ConnID: "",
			Color:  "",
			Score:  18,
			Answer: answer{
				PlayerID: "555",
				Answer:   "ans2",
			},
		},
		incr:    1,
		hiScore: 19,
		tied:    false,
	}, expected: adjScExpected{
		Lp: livePlayer{
			Name:   "",
			ConnID: "",
			Color:  "",
			Score:  19,
			Answer: answer{
				PlayerID: "555",
				Answer:   "ans2",
			},
		},
		HiScore: 19,
		Tied:    true,
	}, description: "incr by 1 to tie"},
}

var gamTests = []struct {
	input       liveGame
	expected    map[string][]string
	description string
}{
	{input: liveGame{
		Pk:          "",
		Sk:          "",
		CurrentWord: "",
		Players: map[string]livePlayer{
			"p1": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  11,
				Answer: answer{
					PlayerID: "111",
					Answer:   "ans1",
				},
			},
			"p2": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  21,
				Answer: answer{
					PlayerID: "222",
					Answer:   "ans1",
				},
			},
			"p3": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  8,
				Answer: answer{
					PlayerID: "333",
					Answer:   "ans1",
				},
			},
			"p4": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  14,
				Answer: answer{
					PlayerID: "444",
					Answer:   "ans2",
				},
			},
			"p5": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  18,
				Answer: answer{
					PlayerID: "555",
					Answer:   "ans2",
				},
			},
		},
		AnswersCount: 0,
		SendToFront:  false,
		HiScore:      21,
		GameTied:     false,
	}, expected: map[string][]string{
		"ans1": {"111", "222", "333"},
		"ans2": {"444", "555"},
	}, description: "2 answers, 5 players"},
	{input: liveGame{
		Pk:          "",
		Sk:          "",
		CurrentWord: "",
		Players: map[string]livePlayer{
			"p1": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  11,
				Answer: answer{
					PlayerID: "111",
					Answer:   "ans1",
				},
			},
			"p2": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  21,
				Answer: answer{
					PlayerID: "222",
					Answer:   "ans2",
				},
			},
			"p3": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  8,
				Answer: answer{
					PlayerID: "333",
					Answer:   "ans3",
				},
			},
			"p4": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  14,
				Answer: answer{
					PlayerID: "444",
					Answer:   "ans4",
				},
			},
			"p5": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  18,
				Answer: answer{
					PlayerID: "555",
					Answer:   "ans5",
				},
			},
		},
		AnswersCount: 0,
		SendToFront:  false,
		HiScore:      21,
		GameTied:     false,
	}, expected: map[string][]string{
		"ans1": {"111"},
		"ans2": {"222"},
		"ans3": {"333"},
		"ans4": {"444"},
		"ans5": {"555"},
	}, description: "5 answers, 5 players"},
	{input: liveGame{
		Pk:          "",
		Sk:          "",
		CurrentWord: "",
		Players: map[string]livePlayer{
			"p1": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  11,
				Answer: answer{
					PlayerID: "111",
					Answer:   "ans1",
				},
			},
			"p2": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  21,
				Answer: answer{
					PlayerID: "222",
					Answer:   "ans1",
				},
			},
			"p3": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  8,
				Answer: answer{
					PlayerID: "333",
					Answer:   "ans2",
				},
			},
			"p4": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  14,
				Answer: answer{
					PlayerID: "444",
					Answer:   "ans2",
				},
			},
			"p5": {
				Name:   "",
				ConnID: "",
				Color:  "",
				Score:  18,
				Answer: answer{
					PlayerID: "555",
					Answer:   "ans3",
				},
			},
		},
		AnswersCount: 0,
		SendToFront:  false,
		HiScore:      21,
		GameTied:     false,
	}, expected: map[string][]string{
		"ans1": {"111", "222"},
		"ans2": {"333", "444"},
		"ans3": {"555"},
	}, description: "3 answers, 5 players"},
}
