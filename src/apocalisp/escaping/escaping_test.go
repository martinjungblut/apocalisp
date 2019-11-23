package escaping

import (
	"testing"
)

func Test_EscapeString_UnescapeString(t *testing.T) {
	mapping := map[string]string{
		// spaces
		"":    "",
		" ":   " ",
		"   ": "   ",
		// doublequotes
		"\\\"":      "\"",
		"\\\" \\\"": "\" \"",
		// newlines
		"\\n":     "\n",
		"\\n \\n": "\n \n",
		// backslashes
		"\\\\":        "\\",
		" \\\\ \\\\ ": " \\ \\ ",
	}

	for a, b := range mapping {
		var output string

		output = UnescapeString(a)
		if output != b {
			t.Errorf("unescapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", a, b, output)
		}

		output = EscapeString(b)
		if output != a {
			t.Errorf("escapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", b, a, output)
		}
	}
}
