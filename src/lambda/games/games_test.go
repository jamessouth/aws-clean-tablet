package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestModifyLiveGamePayload_MarshalJSON(t *testing.T) {
	t.Skip()
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

func TestModifyLiveGamePayload_MarshalJSON2(t *testing.T) {
	t.Skip()
	for _, ref := range modifyLiveGamePayload_MarshalJSON_Tests2 {
		j, err := json.Marshal(ref.input)
		// t.Log(string(j))
		if err != nil {
			t.Fatalf("MarshalJSON() returned %q, want nil.", err)
		}

		if !cmp.Equal(j, []byte("null")) {
			t.Fatalf("MarshalJSON() result: %v,\n  want:%v.", j, nil)
		}
	}
}

// ----------------------------------------------------------------------
func TestSortListPlayersByName(t *testing.T) {
	t.Skip()
	for _, ref := range sortListPlayers {
		ref.input.sortByName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortListPlayersByName: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

func TestSortLivePlayersByName(t *testing.T) {
	// t.Skip()
	for _, ref := range sortLivePlayersByName {
		ref.input.sortByName()

		if !cmp.Equal(ref.input, ref.expected) {
			t.Fatalf("SortLivePlayersByName: %s result: %v,\n  want:%v.", ref.description, ref.input, ref.expected)
		}
	}
}

// func TestSortLivePlayersByScore(t *testing.T) {
// 	t.Skip()
// 	for _, ref := range sortLivePlayersByScore {
// 		j := ref.input.sort(scores)

// 		if !cmp.Equal(j, ref.expected) {
// 			t.Fatalf("SortLivePlayersByScore: %s result: %v,\n  want:%v.", ref.description, j, ref.expected)
// 		}
// 	}
// }

// func TestSortLivePlayersByAnswer(t *testing.T) {
// 	t.Skip()
// 	for _, ref := range sortLivePlayersByAnswer {
// 		j := ref.input.sort(answers)

// 		if !cmp.Equal(j, ref.expected) {
// 			t.Fatalf("SortLivePlayersByAnswer: %s result: %v,\n  want:%v.", ref.description, j, ref.expected)
// 		}
// 	}
// }

// func TestSortLivePlayersByAnswerThenName(t *testing.T) {
// 	t.Skip()
// 	for _, ref := range sortLivePlayersByAnswerThenName {
// 		j := ref.input.sort(answers, namesLive)

// 		if !cmp.Equal(j, ref.expected) {
// 			t.Fatalf("SortLivePlayersByAnswerThenName: %s result: %v,\n\n  want:%v.", ref.description, j, ref.expected)
// 		}
// 	}
// }
