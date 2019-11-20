package escaping

import (
	"fmt"
	"strings"
)

func EscapeString(input string) string {
	if len(strings.Trim(input, " ")) == 0 {
		return input
	}

	runes := []rune(input)
	output := string(runes[1 : len(runes)-1])
	output = strings.ReplaceAll(output, "\\\"", "\"")
	output = strings.ReplaceAll(output, "\\n", "\n")
	output = strings.ReplaceAll(output, "\\\\", "\\")

	return fmt.Sprintf("\"%s\"", output)
}

func UnescapeString(input string) string {
	if len(strings.Trim(input, " ")) == 0 {
		return input
	}

	runes := []rune(input)
	output := string(runes[1 : len(runes)-1])
	output = strings.ReplaceAll(output, "\\", "\\\\")
	output = strings.ReplaceAll(output, "\n", "\\n")
	output = strings.ReplaceAll(output, "\"", "\\\"")

	return fmt.Sprintf("\"%s\"", output)
}
