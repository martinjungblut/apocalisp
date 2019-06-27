package manalispcore

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"strconv"
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
	reader := reader{
		position: 0,
		tokens:   tokenize(sexpr),
	}

	readForm(reader)
}

type MalType struct {
	_integer *int64
	_symbol  *string
	_list    []MalType
}

func readForm(r Reader) MalType {
	token := r.peek()

	if token == "(" {
		return readList(r)
	} else {
		return readAtom(r)
	}
}

func readList(r Reader) MalType {
	list := MalType{_list: make([]MalType, 0)}

	for {
		token := r.next()
		if token == ")" {
			break
		} else {
			list._list = append(list._list, readForm(r))
		}
	}

	return list
}

func readAtom(r Reader) MalType {
	token := r.next()

	i, err := strconv(token, 10, 64)
	if err == nil {
		return MalType{_integer: &i}
	}

	return MalType{_symbol: &token}
}
