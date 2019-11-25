package escaping

import (
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

func UnescapeString(input string) string {
	if len(strings.Trim(input, " ")) == 0 {
		return input
	}

	output := input
	output = strings.ReplaceAll(output, `\\`, "\u029e")
	output = strings.ReplaceAll(output, `\"`, "\"")
	output = strings.ReplaceAll(output, `\n`, "\n")
	output = strings.ReplaceAll(output, "\u029e", "\\")
	return output
}
