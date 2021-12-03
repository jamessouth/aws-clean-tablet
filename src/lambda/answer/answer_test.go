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
		act1, act2 := checkHiScore(test.input.score, test.input.hiScore, test.input.tied)

		act := chsExpected{
			HiScore: act1,
			Tied:    act2,
		}

		if !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - checkHiScore %s act: %+v\n; expected: %+v\n",
				test.description, act, test.expected)
		}

		t.Logf("PASS checkHiScore - %s", test.description)
	}
}

func TestAdjScore(t *testing.T) {
	t.Skip()
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

func TestGetAnswersMap(t *testing.T) {
	// t.Skip()
	for _, test := range gamTests {
		act := getAnswersMap(test.input)

		if !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - getAnswersMap %s act: %+v\n; expected: %+v\n",
				test.description, act, test.expected)
		}
		t.Logf("PASS getAnswersMap - %s", test.description)
	}
}
