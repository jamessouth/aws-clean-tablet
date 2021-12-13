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

func TestModifyLiveGamePayload_MarshalJSON2(t *testing.T) {
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
