package escaping

import (
	"testing"
)

func Test_EscapeString_UnescapeString(t *testing.T) {
	mapping := map[string]string{
		"":            "",
		" ":           " ",
		"   ":         "   ",
		"\\\"":        "\"",
		"\\\" \\\"":   "\" \"",
		"\\\\":        "\\",
		" \\\\ \\\\ ": " \\ \\ ",
	}

	for a, b := range mapping {
		var output string

		output = EscapeString(a)
		if output != b {
			t.Errorf("escapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", a, b, output)
		}

		output = UnescapeString(b)
		if output != a {
			t.Errorf("unescapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", b, a, output)
		}
	}
}

func Test_EscapeString(t *testing.T) {
	mapping := map[string]string{
		"\n":    "\\n",
		"\n \n": "\\n \\n",
	}

	for a, b := range mapping {
		output := EscapeString(a)
		if output != b {
			t.Errorf("escapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", a, b, output)
		}
	}
}
