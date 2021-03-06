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
	input, expected, description string
}{
	{
		input:       "j",
		expected:    "",
		description: "too short",
	},
	{
		input:       "jjjjjjjjjjjjj",
		expected:    "",
		description: "too long",
	},
	{
		input:       "bgt5gb",
		expected:    "",
		description: "number",
	},
	{
		input:       "\nbhbhvg",
		expected:    "",
		description: "newline",
	},
	{
		input:       "bhbhvg\t",
		expected:    "",
		description: "tab",
	},
	{
		input:       "m*.kjns",
		expected:    "",
		description: "symbols",
	},
	{
		input:       "  j",
		expected:    "",
		description: "begins with spaces",
	},
	{
		input:       "mkjns  ",
		expected:    "",
		description: "ends with spaces",
	},
	{
		input:       "bhb hv g",
		expected:    "bhb hv g",
		description: "ok",
	},
}
