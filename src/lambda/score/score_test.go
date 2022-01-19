package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetAnswersMap(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getAnswersMap(); !cmp.Equal(act.Answers, test.expected.Answers) {
			t.Errorf("FAIL - getAnswersMap - %s\n act: %+v\n exp: %+v\n",
				test.description.Description, act.Answers, test.expected.Answers)
		}
	}
}

func TestGetScoresMap(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getAnswersMap().getScoresMap(); !cmp.Equal(act.Scores, test.expected.Scores) {
			t.Errorf("FAIL - getScoresMap - %s\n act: %+v\n exp: %+v\n",
				test.description.Description, act.Scores, test.expected.Scores)
		}
	}
}

func TestUpdateScoresAndClearAnswers(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getAnswersMap().getScoresMap().updateScoresAndClearAnswers(); !cmp.Equal(act.Players, test.expected.Players) {
			t.Errorf("FAIL - updateScoresAndClearAnswers - %s\n act: %+v\n exp: %+v\n",
				test.description.Description, act.Players, test.expected.Players)
		}
	}
}

func TestGetHiScoreAndTie(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getAnswersMap().getScoresMap().updateScoresAndClearAnswers().getHiScoreAndTie(); !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - getHiScoreAndTie - %s\n act: %+v\n exp: %+v\n",
				test.description.Description, act, test.expected)
		}
	}
}
