package main

import (
	"fmt"
	"testing"
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
