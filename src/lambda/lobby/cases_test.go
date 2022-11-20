package main

var noErrorTests = []struct {
	input, exp1, exp2, description string
}{

	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "join"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "join",
		description: "ok",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "join"
		 }`,
		exp1:        "9156849584651978018",
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
			"aW5mb3Jt": "iuiuuhiu"
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
		msg:         "improper json input - too long",
		description: "input too long",
	},
	{
		input: `{
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - duplicate/missing key",
		description: "missing gameno key",
	},
	{
		input: `{
			"gameno": "9156849584651978018"
		 }`,
		msg:         "improper json input - duplicate/missing key",
		description: "missing aW5mb3Jt key",
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
		"aW5mb3Jt": "iuiuuhiu"}`,
		msg:         "improper json input - duplicate/missing key",
		description: "duplicate gameno key",
	},
	{
		input: `{ 
			"gameno": "9156849584651018",
		"aW5mb3Jt": "iuiuuhiu",
		"aW5mb3Jt": "iuiuuhiu"
		}`,
		msg:         "improper json input - duplicate/missing key",
		description: "duplicate aW5mb3Jt key",
	},
	{
		input: `{
			"gameno": "9156849584651018",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno too short",
	},
	{
		input: `{
			"gameno": "91568495846519p8018",
			"aW5mb3Jt": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno",
		description: "gameno has letters",
	},
	{
		input: `{
			"gameno": "91568495846519>8018",
			"aW5mb3Jt": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno",
		description: "gameno has symbol",
	},
	{
		input: `{
			"gameno": "91568495846519 8018",
			"aW5mb3Jt": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno",
		description: "gameno has space",
	},
	{
		input: `{
			"gameno": "915684958465199878018",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno too long",
	},
	{
		input: `{
			"gameno": "disco",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'discon' exactly",
	},
	{
		input: `{
			"gameno": "discoN",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'discon' exactly",
	},
	{
		input: `{
			"gameno": "discon ",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'discon' exactly",
	},
	{
		input: `{
			"gameno": "newgamee",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": "Newgame",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": " newgame",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'newgame' exactly",
	},
}
