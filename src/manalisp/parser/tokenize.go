package parser

import (
	"regexp"
)

func Tokenize(sexpr string) []string {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	rawTokens := []string{}
	for _, group := range re.FindAllStringSubmatch(sexpr, -1) {
		if (group[1] == "") || (group[1][0] == ';') {
			continue
		}
		rawTokens = append(rawTokens, group[1])
	}

	tokens := []string{}
	for index, rawToken := range rawTokens {
		lToken := rawToken
		rToken := rawToken
		if index+1 < len(rawTokens) {
			rToken = rawTokens[index+1]
			if lToken == "~" && rToken == "@" {
				tokens = append(tokens, "~@")
			} else {
				tokens = append(tokens, rawToken)
			}
		} else {
			tokens = append(tokens, rawToken)
		}
	}

	return tokens
}
