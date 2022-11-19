package main

import (
	"bytes"
	"io"
)

var getWordTests = []struct {
	input                 io.ReadCloser
	expected, description string
}{
	{
		input:       io.NopCloser(bytes.NewReader([]byte{97, 98, 98, 114, 101, 118, 105, 97, 116, 105, 110, 103, 10, 97, 98, 98, 114, 101, 118, 105, 97, 116, 105, 111, 110, 10})),
		expected:    "abbreviation",
		description: "consecutive 12-letter words, start at beginning of word",
	},
	{
		input:       io.NopCloser(bytes.NewReader([]byte{101, 118, 105, 97, 116, 105, 110, 103, 10, 97, 98, 98, 114, 101, 118, 105, 97, 116, 105, 111, 110, 10, 97, 98, 98, 114})),
		expected:    "abbreviation",
		description: "consecutive 12-letter words, start in middle of word",
	},
	{
		input:       io.NopCloser(bytes.NewReader([]byte{10, 97, 98, 98, 114, 101, 118, 105, 97, 116, 105, 111, 110, 10, 97, 98, 98, 114, 101, 118, 105, 97, 116, 111, 114, 115})),
		expected:    "abbreviation",
		description: "consecutive 12-letter words, start with newline",
	},
	{
		input:       io.NopCloser(bytes.NewReader([]byte{110, 118, 105, 10, 100, 106, 103, 111, 10, 97, 97, 97, 101, 10, 104, 97, 108, 10, 105, 117, 105, 106, 10, 104, 104, 103})),
		expected:    "djgo",
		description: "small words",
	},
}

var noErrorTests = []struct {
	input, exp1, exp2, description string
}{
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": ""
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt too short",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiuhiuhiuhiuiuiuhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt too long",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiu5uhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt has number",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiu\nuhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt has newline",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiu\tuhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt has tab",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiu<uhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt has symbols",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": " iuiuuhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt begins with spaces",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiuuhiu "
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "",
		description: "aW5mb3Jt ends with spaces",
	},
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "iuiuuhiu"
		 }`,
		exp1:        "9156849584651978018",
		exp2:        "iuiuuhiu",
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
}
