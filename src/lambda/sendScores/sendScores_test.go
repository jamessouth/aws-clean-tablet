package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetStats(t *testing.T) {
	t.Skip()
	if act := getStats(getStatsTest.players, getStatsTest.playersList); !cmp.Equal(act, getStatsTest.expected) {
		t.Errorf("FAIL - getStats - %s\n act: %+v\n exp: %+v\n",
			getStatsTest.description, act, getStatsTest.expected)
	}
}

func TestUpdateScores(t *testing.T) {
	// t.Skip()
	act, _ := updateScores(updateScoresTest.players, updateScoresTest.scores)
	sortByScoreThenName(act)
	sortByScoreThenName(updateScoresTest.expected)
	if !cmp.Equal(act, updateScoresTest.expected) {
		t.Errorf("FAIL - updateScores - %s\n act: %+v\n exp: %+v\n",
			updateScoresTest.description, act, updateScoresTest.expected)
	}
}

// func TestGetScoresMap(t *testing.T) {
// 	// t.Skip()
// 	for _, test := range bunchOfTests {
// 		if act := test.input.getScoresMap(); !cmp.Equal(act.Scores, test.expected.Scores) {
// 			t.Errorf("FAIL - getScoresMap - %s\n act: %+v\n exp: %+v\n",
// 				test.description, act.Scores, test.expected.Scores)
// 		}
// 	}
// }

// func TestUpdateScoresAndClearAnswers(t *testing.T) {
// 	// t.Skip()
// 	for _, test := range bunchOfTests {
// 		if act := test.input.updateScoresAndClearAnswers(); !cmp.Equal(act.Players, test.expected.Players) {
// 			t.Errorf("FAIL - updateScoresAndClearAnswers - %s\n act: %+v\n exp: %+v\n",
// 				test.description, act.Players, test.expected.Players)
// 		}
// 	}
// }

// func TestGetWinner(t *testing.T) {
// 	// t.Skip()
// 	for _, test := range bunchOfTests {
// 		if act := test.input.getWinner(); act.Winner != test.expected.Winner {
// 			t.Errorf("FAIL - getWinner - %s\n act: %s\n exp: %s\n",
// 				test.description, act.Winner, test.expected.Winner)
// 		}
// 	}
// }
