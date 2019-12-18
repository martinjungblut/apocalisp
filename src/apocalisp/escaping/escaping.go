package escaping

import (
	"errors"
	"strings"
)

func EscapeString(input string) string {
	if len(strings.Trim(input, " ")) == 0 {
		return input
	}

	output := input
	output = strings.ReplaceAll(output, "\\", "\u029e")
	output = strings.ReplaceAll(output, "\n", `\n`)
	output = strings.ReplaceAll(output, "\"", `\"`)
	output = strings.ReplaceAll(output, "\u029e", `\\`)
	return output
}

func UnescapeString(input string) (string, error) {
	if len(strings.Trim(input, " ")) == 0 {
		return input, nil
	}

	eof := func() (string, error) {
		return "", errors.New("Error: unexpected EOF.")
	}
	runes := []rune(strings.ReplaceAll(input, " ", ""))
	unbalanced := false
	for i, r := range runes {
		if unbalanced {
			if r == '\\' || r == '"' || r == 'n' {
				unbalanced = false
			} else {
				return eof()
			}
		} else if r == '\\' {
			if i+1 < len(runes) {
				unbalanced = true
			} else {
				return eof()
			}
		}
	}

	output := input
	output = strings.ReplaceAll(output, `\\`, "\u029e")
	output = strings.ReplaceAll(output, `\"`, "\"")
	output = strings.ReplaceAll(output, `\n`, "\n")
	output = strings.ReplaceAll(output, "\u029e", "\\")
	return output, nil
}
