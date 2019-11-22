package escaping

import (
	"strings"
)

func EscapeString(input string) string {
	if len(strings.Trim(input, " ")) == 0 {
		return input
	}

	output := strings.ReplaceAll(input, "\\\"", "\"")
	output = strings.ReplaceAll(output, "\\n", "\n")
	output = strings.ReplaceAll(output, "\\\\", "\\")

	return output
}

func UnescapeString(input string) string {
	if len(strings.Trim(input, " ")) == 0 {
		return input
	}

	output := strings.ReplaceAll(input, "\\", "\\\\")
	output = strings.ReplaceAll(output, "\n", "\\n")
	output = strings.ReplaceAll(output, "\"", "\\\"")

	return output
}
