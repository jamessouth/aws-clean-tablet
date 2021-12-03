package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestUpdateScores(t *testing.T) {
// 	for _, test := range updateScoresTests {
// 		if actual := updateScores(test.input); !cmp.Equal(actual, test.expected) {
// 			t.Errorf("FAIL - updateScores(%s: %+v\n), expected: %+v\n.",
// 				test.description.Description, actual, test.expected)
// 		}
// 		t.Logf("PASS updateScores - %s", test.description.Description)
// 	}
// }

func TestCheckHiScore(t *testing.T) {
	for _, test := range checkHiScoreTests {
		if actual1, _ := checkHiScore(test.input.score, test.input.hiScore, test.input.tied); !cmp.Equal(actual1, test.expected.hiScore) {
			t.Errorf("FAIL - checkHiScore %s actual: %d; expected: %d",
				test.description, actual1, test.expected.hiScore)
		}
		if _, actual2 := checkHiScore(test.input.score, test.input.hiScore, test.input.tied); !cmp.Equal(actual2, test.expected.tied) {
			t.Errorf("FAIL - checkHiScore %s actual: %t; expected: %t",
				test.description, actual2, test.expected.tied)
		}
		t.Logf("PASS checkHiScore - %s", test.description)
	}
}

// func TestRunLengthDecode(t *testing.T) {
// 	for _, test := range decodeTests {
// 		if actual := RunLengthDecode(test.input); actual != test.expected {
// 			t.Errorf("FAIL %s - RunLengthDecode(%s) = %q, expected %q.",
// 				test.description, test.input, actual, test.expected)
// 		}
// 		t.Logf("PASS RunLengthDecode - %s", test.description)
// 	}
// }
// func TestRunLengthEncodeDecode(t *testing.T) {
// 	for _, test := range encodeDecodeTests {
// 		if actual := RunLengthDecode(RunLengthEncode(test.input)); actual != test.expected {
// 			t.Errorf("FAIL %s - RunLengthDecode(RunLengthEncode(%s)) = %q, expected %q.",
// 				test.description, test.input, actual, test.expected)
// 		}
// 		t.Logf("PASS %s", test.description)
// 	}
// }

// func BenchmarkRunLengthDecode(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, test := range decodeTests {
// 			RunLengthDecode(test.input)
// 		}
// 	}
// }
