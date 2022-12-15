package main

var noErrorTests = []struct {
	input, exp1, exp2, description string
}{
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "join"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "join",
		description: "ok",
	},
	{
		input: `{
			"gameno": "newgame",
			"command": "join"
		 }`,
		exp1:        "newgame",
		exp2:        "join",
		description: "ok",
	},
}

var jsonTests = []struct {
	input, description string
}{
	{
		input: `{
			"gameno": "9156849584651978018",
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
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - duplicate/missing key",
		description: "missing gameno key",
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
			"gameno": "9156849584651018",
			"gameno": "9156849584651018",
		"command": "iuiuuhiu"}`,
		msg:         "improper json input - duplicate/missing key",
		description: "duplicate gameno key",
	},
	{
		input: `{ 
			"gameno": "9156849584651018",
		"command": "iuiuuhiu",
		"command": "iuiuuhiu"
		}`,
		msg:         "improper json input - duplicate/missing key",
		description: "duplicate command key",
	},
	{
		input: `{
			"gameno": "9156849584651018",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno: 9156849584651018",
		description: "gameno too short",
	},
	{
		input: `{
			"gameno": "91568495846519p8018",
			"command": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno: 91568495846519p8018",
		description: "gameno has letters",
	},
	{
		input: `{
			"gameno": "91568495846519>8018",
			"command": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno: 91568495846519>8018",
		description: "gameno has symbol",
	},
	{
		input: `{
			"gameno": "91568495846519 8018",
			"command": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno: 91568495846519 8018",
		description: "gameno has space",
	},
	{
		input: `{
			"gameno": "915684958465199878018",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno: 915684958465199878018",
		description: "gameno too long",
	},
	{
		input: `{
			"gameno": "newgamee",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno: newgamee",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": "Newgame",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno: Newgame",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": " newgame",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno:  newgame",
		description: "gameno not 'newgame' exactly",
	},

	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "join "
		 }`,
		msg:         "improper json input - bad command: join ",
		description: "command not 'leave|join' exactly",
	},
}
