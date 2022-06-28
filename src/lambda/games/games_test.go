package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortListPlayers(t *testing.T) {
	// t.Skip()
	for _, ref := range sortListPlayers {
		ref.input = sortByName(ref.input)

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortListPlayers: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

func TestSortByAnswerThenName(t *testing.T) {
	// t.Skip()
	for _, ref := range sortByAnswerThenName {
		ref.input.sortByAnswerThenName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortByAnswerThenName: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

func TestSortByScoreThenName(t *testing.T) {
	// t.Skip()
	for _, ref := range sortByScoreThenName {
		ref.input.sortByScoreThenName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortByScoreThenName: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

func TestListGamePayload_MarshalJSON(t *testing.T) {
	// t.Skip()
	for _, ref := range listGamePayload_MarshalJSON_Tests {
		j, err := json.Marshal(ref.input)
		// t.Log(string(j))
		if err != nil {
			t.Fatalf("MarshalJSON() returned %q, want nil.", err)
		}
		if !cmp.Equal(j, ref.expected) {
			t.Fatalf("MarshalJSON(): %s result: %v,\n  want:%v.", ref.description, j, ref.expected)
		}
	}
}
