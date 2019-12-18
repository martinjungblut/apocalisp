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
		`\"`:    "\"",
		`\" \"`: "\" \"",
		// newlines
		`\n`:    "\n",
		`\n \n`: "\n \n",
		// backslashes
		`\\`:      "\\",
		` \\ \\ `: " \\ \\ ",
		// blackslashes + newlines
		`\\n`: "\\n",
	}

	for a, b := range mapping {
		var intermediary string
		var output string
		var err error

		output, err = UnescapeString(a)
		if err != nil {
			t.Error(err)
		}
		if output != b {
			t.Errorf("unescapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", a, b, output)
		}

		output = EscapeString(b)
		if output != a {
			t.Errorf("escapeString() failed. Input: `%s`. Expected output: `%s`. Actual output: `%s`.", b, a, output)
		}

		output, err = UnescapeString(EscapeString(b))
		if err != nil {
			t.Error(err)
		}
		if output != b {
			t.Errorf("Conversion failed. Expected output: `%s`. Actual output: `%s`.", b, output)
		}

		intermediary, err = UnescapeString(a)
		if err != nil {
			t.Error(err)
		}
		output = EscapeString(intermediary)
		if output != a {
			t.Errorf("Conversion failed. Expected output: `%s`. Actual output: `%s`.", a, output)
		}
	}
}

func Test_UnescapeString(t *testing.T) {
	shouldFail := []string{`\`, `\\\`}

	for _, s := range shouldFail {
		_, err := UnescapeString(s)

		if err == nil {
			t.Error("Should have failed.")
		}
	}
}
