package main

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandler(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act, err := handler(test.input.ctx, test.input.inp); !cmp.Equal(act, test.expected_rv) || errors.Is(err, test.expected_err) {
			t.Errorf("FAIL - handler - %s\n act: %+v\n err: %s\n exp_rv: %+v\n exp_err: %s\n",
				test.description, act, err, test.expected_rv, test.expected_err)
		}
	}
}
