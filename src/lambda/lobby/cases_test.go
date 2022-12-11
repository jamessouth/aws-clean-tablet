package main

var noErrorTests = []struct {
	input, exp1, exp2, description string
}{

	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "disconnect"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "disconnect",
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
			"command": "leave"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "leave",
		description: "ok",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "ready"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "ready",
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
			"gameno": "9156849584651978018",
			"command": "disconnec"
		 }`,
		msg:         "improper json input - bad command",
		description: "command not 'disconnect|join|leave|ready|unready' exactly",
	},

	{
		input: `{
			"gameno": "9156849584651978018",
			"command": " leave"
		 }`,
		msg:         "improper json input - bad command",
		description: "command not 'disconnect|join|leave|ready|unready' exactly",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"command": "reaDy"
		 }`,
		msg:         "improper json input - bad command",
		description: "command not 'disconnect|join|leave|ready|unready' exactly",
	},

	{
		input: `{
			"gameno": "discon",
			"command": "leave"
			}`,
		msg:         "improper json input - leave/discon mismatch",
		description: "can only leave existing game",
	},

	{
		input: `{
			"gameno": "discon",
			"command": "ready"
			}`,
		msg:         "improper json input - ready/discon mismatch",
		description: "can only ready existing game",
	},
}
