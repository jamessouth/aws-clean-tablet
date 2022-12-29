package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/stretchr/testify/assert"
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

func TestCheckInput(t *testing.T) {
	// t.Skip()
	for _, test := range noErrorTests {
		if act1, err := checkInput(test.input); assert.NoErrorf(t, err, "FAIL - checkInput - %s\n err: %+v\n", test.description, err) {
			assert.Equalf(t, act1, test.exp1, "FAIL - checkInput - %s\n act1: %s\n exp1: %s\n", test.description, act1, test.exp1)
		}
	}

	for _, test := range errorTests {
		if act1, err := checkInput(test.input); assert.EqualErrorf(t, err, test.msg, "FAIL - checkInput - %s\n err: %+v\n", test.description, err) {
			assert.Equalf(t, act1, "", "FAIL - checkInput - %s\n act1: %s\n exp1: \n", test.description, act1)
		}
	}

	for _, test := range jsonTests {
		if act1, err := checkInput(test.input); assert.Errorf(t, err, "FAIL - checkInput - %s\n err: %+v\n", test.description, err) {
			assert.Equalf(t, act1, "", "FAIL - checkInput - %s\n act1: %s\n exp1: \n", test.description, act1)
		}
	}
}
