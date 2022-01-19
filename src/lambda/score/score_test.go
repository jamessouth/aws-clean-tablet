package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetAnswersMap(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		act := test.input.getAnswersMap()

		if !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - getAnswersMap %s act: %+v\n; expected: %+v\n",
				test.description, act, test.expected)
		}
		t.Logf("PASS getAnswersMap - %s", test.description)
	}
}

func TestGetScoresMap(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getScoresMap(); !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - getScoresMap(%s: %+v\n), expected: %+v\n.",
				test.description.Description, act, test.expected)
		}
		t.Logf("PASS getScoresMap - %s", test.description.Description)
	}
}

func TestUpdateScoresAndClearAnswers(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		act := test.input.updateScoresAndClearAnswers()

		if !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - updateScoresAndClearAnswers %s act: %+v\n; expected: %+v\n",
				test.description, act, test.expected)
		}

		t.Logf("PASS updateScoresAndClearAnswers - %s", test.description)
	}
}

func TestGetHiScoreAndTie(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		act := test.input.getHiScoreAndTie()

		if !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - getHiScoreAndTie %s act: %+v\n; expected: %+v\n",
				test.description, act, test.expected)
		}
		t.Logf("PASS getHiScoreAndTie - %s", test.description)
	}
}
