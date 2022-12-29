package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortByWinsThenName(t *testing.T) {
	// t.Skip()
	for _, ref := range sortByWinsThenNameTests {
		ref.input.sortByWinsThenName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortByWinsThenName: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

func TestCalcStats(t *testing.T) {
	// t.Skip()
	for _, test := range calcStatsTests {
		if act := test.input.calcStats(); !cmp.Equal(act, test.expected) {
			t.Errorf("FAIL - calcStats - %s\n act: %+v\n exp: %+v\n",
				test.description, act, test.expected)
		}
	}
}
