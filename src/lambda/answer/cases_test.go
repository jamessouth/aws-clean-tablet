package main

import (
	"bytes"
	"io"
)

var bunchOfTests = []struct {
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

var bunchOfTests2 = []struct {
	input, expected1, expected2, description string
	expected3                                error
}{
	{
		input: `{
			"gameno": "9156849584651978018",
			"aW5mb3Jt": "",
		 }`,
		expected1:   "9156849584651978018",
		expected2:   "",
		expected3:   nil,
		description: "aW5mb3Jt too short",
	},
	// {
	// 	input:       "jjjjjjjjjjjjj",
	// 	expected:    "",
	// 	description: "too long",
	// },
	// {
	// 	input:       "bgt5gb",
	// 	expected:    "",
	// 	description: "number",
	// },
	// {
	// 	input:       "\nbhbhvg",
	// 	expected:    "",
	// 	description: "newline",
	// },
	// {
	// 	input:       "bhbhvg\t",
	// 	expected:    "",
	// 	description: "tab",
	// },
	// {
	// 	input:       "m*.kjns",
	// 	expected:    "",
	// 	description: "symbols",
	// },
	// {
	// 	input:       "  j",
	// 	expected:    "",
	// 	description: "begins with spaces",
	// },
	// {
	// 	input:       "mkjns  ",
	// 	expected:    "",
	// 	description: "ends with spaces",
	// },
	// {
	// 	input:       "bhb hv g",
	// 	expected:    "bhb hv g",
	// 	description: "ok",
	// },
	// {
	// 	input:       "bhb hv g",
	// 	expected:    "",
	// 	description: "letters",
	// },
	// {
	// 	input:       "987987987987987987",
	// 	expected:    "",
	// 	description: "too short",
	// },
	// {
	// 	input:       "98765432198765432194",
	// 	expected:    "",
	// 	description: "too long",
	// },
	// {
	// 	input:       "1546879451598456357",
	// 	expected:    "1546879451598456357",
	// 	description: "ok",
	// },
}

// var errKey = errors.New("improper json input - duplicate or missing key")
// var errLen = errors.New("improper json input - too long")

// var bunchOfTests3 = []struct {
// 	input, description string
// 	expected           error
// }{
// 	{
// 		input: `{
// 		   "aW5mb3Jt": "ggg",
// 		}`,
// 		expected:    errKey,
// 		description: "missing gameno key",
// 	},
// 	{
// 		input: `{
// 		   "gameno": "ggg",
// 		}`,
// 		expected:    errKey,
// 		description: "missing aW5mb3Jt key",
// 	},
// 	{
// 		input:       `{}`,
// 		expected:    errKey,
// 		description: "containing no keys",
// 	},
// 	{
// 		input: `{
// 		   "gameno": "ggg",
// 		   "gameno": "gggvvv",
// 		   "aW5mb3Jt": "hhh",
// 		}`,
// 		expected:    errKey,
// 		description: "duplicate gameno key",
// 	},
// 	{
// 		input: `{
// 		   "gameno": "gggvvv",
// 		   "aW5mb3Jt": "hhh",
// 		   "aW5mb3Jt": "hhddh",
// 		}`,
// 		expected:    errKey,
// 		description: "duplicate aW5mb3Jt key",
// 	},
// }

// var nils = []struct {
// 	input, description string
// 	expected           error
// }{
// 	{
// 		input: `{
// 		   "gameno": "ggg",
// 		   "aW5mb3Jt": "hhh",
// 		}`,
// 		expected:    nil,
// 		description: "ok",
// 	},
// }

// var nils2 = []struct {
// 	input, description string
// 	expected           error
// }{
// 	{
// 		input: `{
// 		   "gameno": "ggg",
// 		   "aW5mb3Jt": "hhh",
// 		}`,
// 		expected:    nil,
// 		description: "ok",
// 	},
// }

// var lens = []struct {
// 	input, description string
// 	expected           error
// }{
// 	{
// 		input: `{
// 		   "gameno": "ggg",
// 		   "aW5mb3Jt": "hhh",
// 		   "aW5mb3Jt": "hhh",
// 		   "aW5mb3Jt": "hhh",
// 		   "aW5mb3Jt": "hhh",
// 		}`,
// 		expected:    errLen,
// 		description: "too long",
// 	},
// }
