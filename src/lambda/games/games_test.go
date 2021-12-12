package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestModifyLiveGamePayload_MarshalJSON(t *testing.T) {
	for _, ref := range modifyLiveGamePayload_MarshalJSON_Tests {
		j, err := json.Marshal(ref.input)
		// t.Log(j)
		if err != nil {
			t.Fatalf("MarshalJSON() returned %q, want nil.", err)
		}
		var lg liveGame
		err = json.Unmarshal(j, &lg)
		if err != nil {
			t.Fatalf("json.Unarshal() returned %q, want nil.", err)
		}
		if !cmp.Equal(lg, ref.expected) {
			t.Fatalf("MarshalJSON() result: %v,\n  want:%v.", lg, ref.expected)
		}
	}
}
