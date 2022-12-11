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
	{
		input: `{
			"gameno": "discon",
			"command": "disconnect"
		 }`,
		exp1:        "discon",
		exp2:        "disconnect",
		description: "ok",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "unready"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "unready",
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
		msg:         "improper json input - too long",
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
		msg:         "improper json input - bad gameno",
		description: "gameno too short",
	},
	{
		input: `{
			"gameno": "91568495846519p8018",
			"command": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno",
		description: "gameno has letters",
	},
	{
		input: `{
			"gameno": "91568495846519>8018",
			"command": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno",
		description: "gameno has symbol",
	},
	{
		input: `{
			"gameno": "91568495846519 8018",
			"command": "iuiuuhiu"
			}`,
		msg:         "improper json input - bad gameno",
		description: "gameno has space",
	},
	{
		input: `{
			"gameno": "915684958465199878018",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno too long",
	},
	{
		input: `{
			"gameno": "disco",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'discon' exactly",
	},
	{
		input: `{
			"gameno": "discoN",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'discon' exactly",
	},
	{
		input: `{
			"gameno": "discon ",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'discon' exactly",
	},
	{
		input: `{
			"gameno": "newgamee",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": "Newgame",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": " newgame",
			"command": "iuiuuhiu"
		 }`,
		msg:         "improper json input - bad gameno",
		description: "gameno not 'newgame' exactly",
	},
	{
		input: `{
			"gameno": "discon",
			"command": "disconnec"
		 }`,
		msg:         "improper json input - bad command",
		description: "command not 'disconnect|join|unready' exactly",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "join "
		 }`,
		msg:         "improper json input - bad command",
		description: "command not 'disconnect|join|unready' exactly",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "Unready"
		 }`,
		msg:         "improper json input - bad command",
		description: "command not 'disconnect|join|unready' exactly",
	},
	{
		input: `{
			"gameno": "newgame",
			"command": "disconnect"
			}`,
		msg:         "improper json input - disconnect/newgame mismatch",
		description: "cannot disconnect from game while requesting that it be created",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "disconnect"
			}`,
		msg:         "improper json input - disconnect/newgame mismatch",
		description: "disconnect from numbered game handled elsewhere",
	},
	{
		input: `{
			"gameno": "discon",
			"command": "join"
			}`,
		msg:         "improper json input - join/discon mismatch",
		description: "cannot join and disconnect at same time",
	},
	{
		input: `{
			"gameno": "discon",
			"command": "unready"
			}`,
		msg:         "improper json input - unready/(discon|newgame) mismatch",
		description: "can only unready existing game",
	},
	{
		input: `{
			"gameno": "newgame",
			"command": "unready"
			}`,
		msg:         "improper json input - unready/(discon|newgame) mismatch",
		description: "can only unready existing game",
	},
}
