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
				test.description, act.Answers, test.expected.Answers)
		}
	}
}

func TestGetScoresMap(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getScoresMap(); !cmp.Equal(act.Scores, test.expected.Scores) {
			t.Errorf("FAIL - getScoresMap - %s\n act: %+v\n exp: %+v\n",
				test.description, act.Scores, test.expected.Scores)
		}
	}
}

func TestUpdateScoresAndClearAnswers(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.updateScoresAndClearAnswers(); !cmp.Equal(act.Players, test.expected.Players) {
			t.Errorf("FAIL - updateScoresAndClearAnswers - %s\n act: %+v\n exp: %+v\n",
				test.description, act.Players, test.expected.Players)
		}
	}
}

func TestGetWinner(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := test.input.getWinner(); act.Winner != test.expected.Winner {
			t.Errorf("FAIL - getWinner - %s\n act: %t\n exp: %t\n",
				test.description, act.Winner, test.expected.Winner)
		}
	}
}
