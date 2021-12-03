package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestUpdateScores(t *testing.T) {
// 	for _, test := range updateScoresTests {
// 		if act := updateScores(test.input); !cmp.Equal(act, test.expected) {
// 			t.Errorf("FAIL - updateScores(%s: %+v\n), expected: %+v\n.",
// 				test.description.Description, act, test.expected)
// 		}
// 		t.Logf("PASS updateScores - %s", test.description.Description)
// 	}
// }

func TestCheckHiScore(t *testing.T) {
	t.Skip()
	for _, test := range checkHiScoreTests {
		if act1, _ := checkHiScore(test.input.score, test.input.hiScore, test.input.tied); !cmp.Equal(act1, test.expected.HiScore) {
			t.Errorf("FAIL - checkHiScore %s act: %d; expected: %d",
				test.description, act1, test.expected.HiScore)
		}
		if _, act2 := checkHiScore(test.input.score, test.input.hiScore, test.input.tied); !cmp.Equal(act2, test.expected.Tied) {
			t.Errorf("FAIL - checkHiScore %s act: %t; expected: %t",
				test.description, act2, test.expected.Tied)
		}
		t.Logf("PASS checkHiScore - %s", test.description)
	}
}

func TestAdjScore(t *testing.T) {
	for _, test := range adjScoreTests {
		act1, act2, act3 := adjScore(test.input.old, test.input.incr, test.input.hiScore, test.input.tied)

		act := adjScExpected{
			Lp:      act1,
			HiScore: act2,
			Tied:    act3,
		}

		if !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - adjScore %s act: %+v\n; expected: %+v\n",
				test.description, act, test.expected)
		}
		t.Logf("PASS adjScore - %s", test.description)
	}
}

// func TestRunLengthDecode(t *testing.T) {
// 	for _, test := range decodeTests {
// 		if act := RunLengthDecode(test.input); act != test.expected {
// 			t.Errorf("FAIL %s - RunLengthDecode(%s) = %q, expected %q.",
// 				test.description, test.input, act, test.expected)
// 		}
// 		t.Logf("PASS RunLengthDecode - %s", test.description)
// 	}
// }
// func TestRunLengthEncodeDecode(t *testing.T) {
// 	for _, test := range encodeDecodeTests {
// 		if act := RunLengthDecode(RunLengthEncode(test.input)); act != test.expected {
// 			t.Errorf("FAIL %s - RunLengthDecode(RunLengthEncode(%s)) = %q, expected %q.",
// 				test.description, test.input, act, test.expected)
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
