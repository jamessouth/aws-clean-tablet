package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomByte(t *testing.T) {
	t.Skip()
	for i := 0; i < 20; i++ {
		fmt.Println(getRandomByte())
	}
}

func TestGetWord(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests {
		if act := getWord(test.input); act != test.expected {
			t.Errorf("FAIL - getWord - %s\n act: %+v\n exp: %+v\n",
				test.description, act, test.expected)
		}
	}
}

func TestCheckInput(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests2 {
		if act := checkInput(test.input, test.re); act != test.expected {
			t.Errorf("FAIL - CheckInput - %s\n act: %+v\n exp: %+v\n",
				test.description, act, test.expected)
		}
	}
}

func TestCheckKeys(t *testing.T) {
	// t.Skip()
	for _, test := range bunchOfTests3 {
		err := checkKeys(test.input)
		assert.EqualErrorf(t, err, test.expected.Error(), "FAIL - CheckKeys - %s\n err: %+v\n exp: %s\n", test.description, err, test.expected.Error())
	}
	for _, test := range nils {
		err := checkKeys(test.input)
		assert.Nilf(t, err, "FAIL - CheckKeys - %s\n err: %+v\n exp: %v\n", test.description, err, test.expected)
	}
}

func TestCheckLength(t *testing.T) {
	// t.Skip()
	for _, test := range lens {
		err := checkLength(test.input)
		assert.EqualErrorf(t, err, test.expected.Error(), "FAIL - CheckLength - %s\n err: %+v\n exp: %s\n", test.description, err, test.expected.Error())
	}
	for _, test := range nils2 {
		err := checkLength(test.input)
		assert.Nilf(t, err, "FAIL - CheckLength - %s\n err: %+v\n exp: %v\n", test.description, err, test.expected)
	}
}
