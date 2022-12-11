package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckInput(t *testing.T) {
	// t.Skip()
	for _, test := range noErrorTests {
		if act1, act2, err := checkInput(test.input); assert.NoErrorf(t, err, "FAIL - checkInput - %s\n err: %+v\n", test.description, err) {
			if assert.Equalf(t, act1, test.exp1, "FAIL - checkInput - %s\n act1: %s\n exp1: %s\n", test.description, act1, test.exp1) {
				assert.Equalf(t, act2, test.exp2, "FAIL - checkInput - %s\n act2: %s\n exp2: %s\n", test.description, act2, test.exp2)
			}
		}
	}

	for _, test := range errorTests {
		if act1, act2, err := checkInput(test.input); assert.EqualErrorf(t, err, test.msg, "FAIL - checkInput - %s\n err: %+v\n", test.description, err) {
			if assert.Equalf(t, act1, "", "FAIL - checkInput - %s\n act1: %s\n exp1: \n", test.description, act1) {
				assert.Equalf(t, act2, "", "FAIL - checkInput - %s\n act2: %s\n exp2: \n", test.description, act2)
			}
		}
	}

	for _, test := range jsonTests {
		if act1, act2, err := checkInput(test.input); assert.Errorf(t, err, "FAIL - checkInput - %s\n err: %+v\n", test.description, err) {
			if assert.Equalf(t, act1, "", "FAIL - checkInput - %s\n act1: %s\n exp1: \n", test.description, act1) {
				assert.Equalf(t, act2, "", "FAIL - checkInput - %s\n act2: %s\n exp2: \n", test.description, act2)
			}
		}
	}
}
