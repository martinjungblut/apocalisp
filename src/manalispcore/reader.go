package manalispcore

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type reader struct {
	position uint64
	tokens   []string
}

func (r *reader) peek() string {
	return r.tokens[r.position]
}

func (r *reader) next() string {
	token := r.tokens[r.position]
	r.position++
	return token
}

func tokenize(sexpr string) []string {
	tokens := make([]string, 0)

	regex := pcre.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	matcher := regex.MatcherString(sexpr, 0)
	for i := uint64(0); ; i++ {
		token := matcher.GroupString(i)
		if len(token) > 0 {
			tokens = append(tokens, token)
		} else {
			break
		}
	}

	return tokens
}

func readStr(sexpr string) {
	tokens := tokenize(sexpr)
	reader := reader{
		position: 0,
		tokens:   tokens,
	}
	readForm(reader)
}
