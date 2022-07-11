package main

var (
	five       int = 5
	six        int = 6
	eight      int = 8
	nine       int = 9
	ten        int = 10
	thirteen   int = 13
	twenty     int = 20
	twentyone  int = 21
	twentyfour int = 24
	twentyfive int = 25
	twentysix  int = 26
)

var getStatsTest = struct {
	players     map[string]livePlayer
	playersList []livePlayer
	expected    []stat
	description string
}{
	players: map[string]livePlayer{
		"p1": {
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &ten,
			Answer: "ans1",
		},
		"p2": {
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twenty,
			Answer: "ans1",
		},
		"p3": {
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyone,
			Answer: "ans1",
		},
		"p4": {
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &eight,
			Answer: "ans1",
		},
		"p5": {
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &five,
			Answer: "ans1",
		},
	},
	playersList: []livePlayer{
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyone,
			Answer: "ans1",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twenty,
			Answer: "ans1",
		},
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &ten,
			Answer: "ans1",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &eight,
			Answer: "ans1",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &five,
			Answer: "ans1",
		},
	},
	expected: []stat{
		{
			PlayerID: "p3",
			Name:     "ccc",
			Wins:     "1",
			Points:   "21",
		},
		{
			PlayerID: "p2",
			Name:     "bbb",
			Wins:     "0",
			Points:   "20",
		},
		{
			PlayerID: "p1",
			Name:     "aaa",
			Wins:     "0",
			Points:   "10",
		},
		{
			PlayerID: "p4",
			Name:     "ddd",
			Wins:     "0",
			Points:   "8",
		},
		{
			PlayerID: "p5",
			Name:     "eee",
			Wins:     "0",
			Points:   "5",
		},
	}, description: "get stats",
}

var updateScoresTest = struct {
	players     map[string]livePlayer
	scores      map[string]int
	expected    []livePlayer
	description string
}{
	players: map[string]livePlayer{
		"p1": {
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &ten,
			Answer: "ans1",
		},
		"p2": {
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twenty,
			Answer: "ans1",
		},
		"p3": {
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyone,
			Answer: "ans1",
		},
		"p4": {
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &eight,
			Answer: "ans1",
		},
		"p5": {
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &five,
			Answer: "ans1",
		},
	},
	scores: map[string]int{
		"111": 3,
		"222": 1,
		"333": 3,
		"444": 1,
		"555": 1,
	},
	expected: []livePlayer{
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &thirteen,
			Answer: "",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twentyone,
			Answer: "",
		},
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyfour,
			Answer: "",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &nine,
			Answer: "",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &six,
			Answer: "",
		},
	}, description: "update scores",
}

var sortTest = struct {
	players, expected []livePlayer
	description       string
}{
	players: []livePlayer{
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &thirteen,
			Answer: "",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twentyone,
			Answer: "",
		},
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyfour,
			Answer: "",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &nine,
			Answer: "",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &six,
			Answer: "",
		},
	},
	expected: []livePlayer{
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyfour,
			Answer: "",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twentyone,
			Answer: "",
		},
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &thirteen,
			Answer: "",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &nine,
			Answer: "",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &six,
			Answer: "",
		},
	}, description: "sort players",
}

var winnerTests = []struct {
	players               []livePlayer
	expected, description string
}{
	{players: []livePlayer{
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyfour,
			Answer: "",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twentyone,
			Answer: "",
		},
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &thirteen,
			Answer: "",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &nine,
			Answer: "",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &six,
			Answer: "",
		},
	}, expected: "", description: "no winner",
	},
	{players: []livePlayer{
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyfive,
			Answer: "",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twentyone,
			Answer: "",
		},
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &thirteen,
			Answer: "",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &nine,
			Answer: "",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &six,
			Answer: "",
		},
	}, expected: "ccc", description: "winner",
	},
	{players: []livePlayer{
		{
			Name:   "ccc",
			ConnID: "333",
			Color:  "green",
			Score:  &twentyfive,
			Answer: "",
		},
		{
			Name:   "bbb",
			ConnID: "222",
			Color:  "blue",
			Score:  &twentyfive,
			Answer: "",
		},
		{
			Name:   "aaa",
			ConnID: "111",
			Color:  "red",
			Score:  &thirteen,
			Answer: "",
		},
		{
			Name:   "ddd",
			ConnID: "444",
			Color:  "yellow",
			Score:  &nine,
			Answer: "",
		},
		{
			Name:   "eee",
			ConnID: "555",
			Color:  "indigo",
			Score:  &six,
			Answer: "",
		},
	}, expected: "", description: "tied, no winner",
	},
}
