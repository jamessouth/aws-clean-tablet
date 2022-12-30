package main

var sortByWinsThenNameTests = []struct {
	input, expected stats
	description     string
}{
	{
		input: stats{
			{
				Name:   "will",
				Wins:   3,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "bill",
				Wins:   0,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "lisa",
				Wins:   0,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "erin",
				Wins:   3,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "carl",
				Wins:   3,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "phil",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "norm",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "vera",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
		},
		expected: stats{
			{
				Name:   "norm",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "phil",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "vera",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "carl",
				Wins:   3,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "erin",
				Wins:   3,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "will",
				Wins:   3,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "bill",
				Wins:   0,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "lisa",
				Wins:   0,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
		},
		description: "sort list of stats by wins descending then name ascending",
	},
}

var calcStatsTests = []struct {
	input, expected stats
	description     string
}{
	{
		input: stats{
			{
				Name:   "norm",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "phil",
				Wins:   7,
				Points: 49,
				Games:  12,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "vera",
				Wins:   7,
				Points: 54,
				Games:  21,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "carl",
				Wins:   3,
				Points: 39,
				Games:  32,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "erin",
				Wins:   3,
				Points: 29,
				Games:  14,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "will",
				Wins:   3,
				Points: 48,
				Games:  20,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "bill",
				Wins:   0,
				Points: 63,
				Games:  33,
				WinPct: 0,
				PPG:    0,
			},
			{
				Name:   "lisa",
				Wins:   0,
				Points: 28,
				Games:  26,
				WinPct: 0,
				PPG:    0,
			},
		},
		expected: stats{
			{
				Name:   "norm",
				Wins:   7,
				Points: 59,
				Games:  22,
				WinPct: 0.32,
				PPG:    2.68,
			},
			{
				Name:   "phil",
				Wins:   7,
				Points: 49,
				Games:  12,
				WinPct: 0.58,
				PPG:    4.08,
			},
			{
				Name:   "vera",
				Wins:   7,
				Points: 54,
				Games:  21,
				WinPct: 0.33,
				PPG:    2.57,
			},
			{
				Name:   "carl",
				Wins:   3,
				Points: 39,
				Games:  32,
				WinPct: 0.09,
				PPG:    1.22,
			},
			{
				Name:   "erin",
				Wins:   3,
				Points: 29,
				Games:  14,
				WinPct: 0.21,
				PPG:    2.07,
			},
			{
				Name:   "will",
				Wins:   3,
				Points: 48,
				Games:  20,
				WinPct: 0.15,
				PPG:    2.4,
			},
			{
				Name:   "bill",
				Wins:   0,
				Points: 63,
				Games:  33,
				WinPct: 0,
				PPG:    1.91,
			},
			{
				Name:   "lisa",
				Wins:   0,
				Points: 28,
				Games:  26,
				WinPct: 0,
				PPG:    1.08,
			},
		},
		description: "get stats",
	},
}

var noErrorTests = []struct {
	input, exp1, description string
}{
	{
		input: `{
			"command": "leaders"
		 }`,
		exp1:        "leaders",
		description: "ok",
	},
	{
		input: `{
			"command": "listGames"
		 }`,
		exp1:        "listGames",
		description: "ok",
	},
}

var jsonTests = []struct {
	input, description string
}{
	{
		input: `{
			"command": "iuiuuhiu"
		 }]]]]]`,
		description: "malformed input",
	},
}

var errorTests = []struct {
	input, msg, description string
}{
	{
		input: `{
			"gameno": "9156855555555555555555555555555555549584651018",
			"aW5mhb3Jt": "iuiuu55555555555555555555555555555555555555555hiu",
			"aW5mkb3Jt": "iuiuu55555555555555555555555555555555555555555hiu",
			"aW5m3b3Jt": "iuiuu55555555555555555555555555555555555555555hiu"
		 }`,
		msg:         "improper json input - too long: 275",
		description: "input too long",
	},
	{
		input: `{
			"gameno": "9156849584651978018"
		 }`,
		msg:         "improper json input - duplicate/missing key",
		description: "missing command key",
	},
	{
		input:       `{ }`,
		msg:         "improper json input - duplicate/missing key",
		description: "containing no keys",
	},
	{
		input: `{ 
		"command": "iuiuuhiu",
		"command": "iuiu55uhiu"
		}`,
		msg:         "improper json input - duplicate/missing key",
		description: "duplicate command key",
	},
	{
		input: `{
			"command": " leaders"
		 }`,
		msg:         "improper json input - bad command:  leaders",
		description: "command not 'leaders|listGames' exactly",
	},
	{
		input: `{
			"command": "listGames "
		 }`,
		msg:         "improper json input - bad command: listGames ",
		description: "command not 'leaders|listGames' exactly",
	},
}

var getFrontListGamesTests = []struct {
	input       []backListGame
	expected    []frontListGame
	description string
}{
	{
		input: []backListGame{
			{
				Pk:        "LISTGAME",
				Sk:        "3",
				TimerCxld: true,
				Players: map[string]listPlayer{
					"x8u3": {Name: "will", ConnID: "098"},
					"n38c": {Name: "bill", ConnID: "987"},
					"a9i3": {Name: "carl", ConnID: "876"},
					"km28": {Name: "lisa", ConnID: "765"},
				},
			},
			{
				Pk:        "LISTGAME",
				Sk:        "0",
				TimerCxld: true,
				Players: map[string]listPlayer{
					"x8u3": {Name: "will", ConnID: "098"},
					"n38c": {Name: "bill", ConnID: "987"},
					"a9i3": {Name: "carl", ConnID: "876"},
					"km28": {Name: "lisa", ConnID: "765"},
				},
			},
			{
				Pk:        "LISTGAME",
				Sk:        "10",
				TimerCxld: true,
				Players: map[string]listPlayer{
					"x8u3": {Name: "will", ConnID: "098"},
					"n38c": {Name: "bill", ConnID: "987"},
					"a9i3": {Name: "carl", ConnID: "876"},
					"km28": {Name: "lisa", ConnID: "765"},
				},
			},
			{
				Pk:        "LISTGAME",
				Sk:        "13",
				TimerCxld: true,
				Players: map[string]listPlayer{
					"x8u3": {Name: "will", ConnID: "098"},
					"n38c": {Name: "bill", ConnID: "987"},
					"a9i3": {Name: "carl", ConnID: "876"},
					"km28": {Name: "lisa", ConnID: "765"},
				},
			},
			{
				Pk:        "LISTGAME",
				Sk:        "32",
				TimerCxld: true,
				Players: map[string]listPlayer{
					"x8u3": {Name: "will", ConnID: "098"},
					"n38c": {Name: "bill", ConnID: "987"},
					"a9i3": {Name: "carl", ConnID: "876"},
					"km28": {Name: "lisa", ConnID: "765"},
				},
			},
			{
				Pk:        "LISTGAME",
				Sk:        "7",
				TimerCxld: true,
				Players: map[string]listPlayer{
					"x8u3": {Name: "will", ConnID: "098"},
					"n38c": {Name: "bill", ConnID: "987"},
					"a9i3": {Name: "carl", ConnID: "876"},
					"km28": {Name: "lisa", ConnID: "765"},
				},
			},
		},
		expected: []frontListGame{
			{
				No:        "3",
				TimerCxld: true,
				Players: []listPlayer{
					{Name: "bill", ConnID: "987"},
					{Name: "carl", ConnID: "876"},
					{Name: "lisa", ConnID: "765"},
					{Name: "will", ConnID: "098"},
				},
			},
			{
				No:        "0",
				TimerCxld: true,
				Players: []listPlayer{
					{Name: "bill", ConnID: "987"},
					{Name: "carl", ConnID: "876"},
					{Name: "lisa", ConnID: "765"},
					{Name: "will", ConnID: "098"},
				},
			},
			{
				No:        "10",
				TimerCxld: true,
				Players: []listPlayer{
					{Name: "bill", ConnID: "987"},
					{Name: "carl", ConnID: "876"},
					{Name: "lisa", ConnID: "765"},
					{Name: "will", ConnID: "098"},
				},
			},
			{
				No:        "13",
				TimerCxld: true,
				Players: []listPlayer{
					{Name: "bill", ConnID: "987"},
					{Name: "carl", ConnID: "876"},
					{Name: "lisa", ConnID: "765"},
					{Name: "will", ConnID: "098"},
				},
			},
			{
				No:        "32",
				TimerCxld: true,
				Players: []listPlayer{
					{Name: "bill", ConnID: "987"},
					{Name: "carl", ConnID: "876"},
					{Name: "lisa", ConnID: "765"},
					{Name: "will", ConnID: "098"},
				},
			},
			{
				No:        "7",
				TimerCxld: true,
				Players: []listPlayer{
					{Name: "bill", ConnID: "987"},
					{Name: "carl", ConnID: "876"},
					{Name: "lisa", ConnID: "765"},
					{Name: "will", ConnID: "098"},
				},
			},
		},
		description: "convert list of back end games to front end games",
	},
}
