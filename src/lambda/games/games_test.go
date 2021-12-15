package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestModifyLiveGamePayload_MarshalJSON(t *testing.T) {
	// t.Skip()
	for _, ref := range modifyLiveGamePayload_MarshalJSON_Tests {
		j, err := json.Marshal(ref.input)
		// t.Log(string(j))
		if err != nil {
			t.Fatalf("MarshalJSON() returned %q, want nil.", err)
		}
		var lg liveGameWrap
		err = json.Unmarshal(j, &lg)
		if err != nil {
			t.Fatalf("json.Unmarshal() returned %q, want nil.", err)
		}
		if !cmp.Equal(lg, ref.expected) {
			t.Fatalf("MarshalJSON() result: %v,\n  want:%v.", lg, ref.expected)
		}
	}
}

// func TestModifyLiveGamePayload_MarshalJSON2(t *testing.T) {
// 	// t.Skip()
// 	for _, ref := range modifyLiveGamePayload_MarshalJSON_Tests2 {
// 		j, err := json.Marshal(ref.input)
// 		// t.Log(string(j))
// 		if err != nil {
// 			t.Fatalf("MarshalJSON() returned %q, want nil.", err)
// 		}

// 		if !cmp.Equal(j, []byte("null")) {
// 			t.Fatalf("MarshalJSON() result: %v,\n  want:%v.", j, nil)
// 		}
// 	}
// }

// ----------------------------------------------------------------------
func TestSortListPlayers(t *testing.T) {
	// t.Skip()
	for _, ref := range sortListPlayers {
		ref.input.sortByName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortListPlayers: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

func TestSortLivePlayers(t *testing.T) {
	// t.Skip()
	for _, ref := range sortLivePlayers {
		ref.input.sortByName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortLivePlayers: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
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
