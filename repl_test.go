package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{{
		input:    "  hello world  ",
		expected: []string{"hello", "world"},
	}, {
		input:    "  Charmander Bulbasaur PIKACHU  ",
		expected: []string{"charmander", "bulbasaur", "pikachu"},
	},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("expected [%d] to be the same as [%d]!", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]

			if word != expected {
				t.Errorf("word [%s] does not match expected word [%s]!", word, expected)
			}
		}
	}
}
