package main

import (
	"regexp"
	"testing"
)

var (
	re = regexp.MustCompile(`^____ [a-z]{2,9}$|^[a-z]{2,9} ____$`)
	m  = map[string]bool{}
)

func TestWords(t *testing.T) {
	loop := func(w []string, t string) func() (bool, string, int) {
		if t == "word" {
			return func() (bool, string, int) {
				for i, j := range w {
					if !re.MatchString(j) {
						return false, j, i + 5
					}
				}
				return true, "", 0
			}
		}
		if t == "duplicate" {
			return func() (bool, string, int) {
				for i, j := range w {
					if m[j] {
						return false, j, i + 5
					}
					m[j] = true
				}
				return true, "", 0
			}
		}
		return func() (bool, string, int) {
			return false, "", 0
		}
	}

	tests := map[string]struct {
		words []string
		test  string
	}{
		"each word has: 4 _, 1 space, 2-9 lower-case letters, OR 2-9 lower-case letters, 1 space, 4 _": {words: words, test: "word"},
		"there are no duplicate words in the list":                                                     {words: words, test: "duplicate"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, word, ind := loop(tc.words, tc.test)()
			if !got {
				t.Fatalf("the %s test failed on word: %s, line: %d", name, word, ind+269)
			}
		})
	}
}
