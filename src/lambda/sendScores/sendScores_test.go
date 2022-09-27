package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetStats(t *testing.T) {
	// t.Skip()
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

func TestSort(t *testing.T) {
	// t.Skip()
	act := sortTest.players
	sortByScoreThenName(act)
	if !cmp.Equal(act, sortTest.expected) {
		t.Errorf("FAIL - sortScores - %s\n act: %+v\n exp: %+v\n",
			sortTest.description, act, sortTest.expected)
	}
}

func TestGetWinner(t *testing.T) {
	// t.Skip()
	for _, test := range winnerTests {
		if act := getWinner(test.players, test.lastword); !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - getWinner - %s\n act: %+v\n exp: %+v\n",
				test.description, act, test.expected)
		}
	}
}
