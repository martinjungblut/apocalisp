package manalispcore

import (
	"fmt"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"strconv"
	"strings"
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

	regex := pcre.MustCompile(`[\s,]*(~@|[\[\]{}()'`+"`"+
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`+"`"+
		`,;)]*)`, pcre.MULTILINE)
	matcher := regex.MatcherString(sexpr, 0)
	for i := 0; ; i++ {
		token := matcher.GroupString(i)
		if len(token) > 0 {
			tokens = append(tokens, token)
		} else {
			break
		}
	}

	return tokens
}

func ReadStr(sexpr string) MalType {
	reader := reader{
		position: 0,
		tokens:   tokenize(sexpr),
	}

	return readForm(reader)
}

func PrintStr(t MalType) string {
	if t._integer != nil {
		return fmt.Sprintf("%d", *t._integer)
	} else if t._symbol != nil {
		return *t._symbol
	} else {
		tokens := make([]string, 0)
		for _, _type := range t._list {
			tokens = append(tokens, PrintStr(_type))
		}
		return fmt.Sprintf("(%s)", strings.Join(tokens, " "))
	}
}

type MalType struct {
	_integer *int64
	_symbol  *string
	_list    []MalType
}

func readForm(r reader) MalType {
	token := r.peek()

	if token == "(" {
		return readList(r)
	} else {
		return readAtom(r)
	}
}

func readList(r reader) MalType {
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

func readAtom(r reader) MalType {
	token := r.next()

	i, err := strconv.ParseInt(token, 10, 64)
	if err == nil {
		return MalType{_integer: &i}
	}

	return MalType{_symbol: &token}
}
