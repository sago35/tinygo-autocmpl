package main

import (
	"strings"
	"testing"
)

func TestCompletionBash(t *testing.T) {
	tests := []struct {
		summary  string
		input    []string
		expected string
	}{
		{
			summary:  "no arguments",
			input:    []string{""},
			expected: strings.Join(validCommands, " "),
		},

		{
			summary:  "commands - incomplete",
			input:    []string{"f"},
			expected: strings.Join(validCommands, " "),
		},

		{
			summary:  "commands - incomplete",
			input:    []string{"flas"},
			expected: strings.Join(validCommands, " "),
		},

		{
			summary:  "commands",
			input:    []string{"flash"},
			expected: "",
		},

		{
			summary:  "flag",
			input:    []string{"flash", "-size"},
			expected: "none short full",
		},
	}

	for _, test := range tests {
		if g, e := completionBash(test.input), test.expected; g != e {
			t.Errorf("got %q, want %q (%s)", g, e, test.summary)
		}
	}
}

func TestCompletionBashFlags(t *testing.T) {
	tests := []struct {
		summary  string
		input    []string
		expected string
	}{
		{
			summary:  "flag - incomplete",
			input:    []string{"flash", "-prog"},
			expected: "-programmer",
		},

		{
			summary:  "flag - incomplete && double-dash",
			input:    []string{"flash", "--prog"},
			expected: "--programmer",
		},
	}

	for _, test := range tests {
		got := completionBash(test.input)
		if !strings.Contains(got, test.expected) {
			t.Errorf("%q is not contained in %q (%s)", test.expected, got, test.summary)
		}
	}

}
