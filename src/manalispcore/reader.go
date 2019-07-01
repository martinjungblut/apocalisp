package manalispcore

import (
	"fmt"
	// "github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"regexp"
	"strconv"
	"strings"
)

type reader struct {
	position int
	tokens   []string
}

func (r *reader) peek() *string {
	if r.position < len(r.tokens) {
		return &r.tokens[r.position]
	} else {
		return nil
	}
}

func (r *reader) next() *string {
	if r.position < len(r.tokens) {
		token := &(r.tokens[r.position])
		r.position++
		return token
	} else {
		return nil
	}
}

func tokenize(sexpr string) []string {
	// tokens := make([]string, 0)

	// regex := pcre.MustCompile(`[\s,]*(~@|[\[\]{}()'`+"`"+
	// 	`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`+"`"+
	// 	`,;)]*)`, pcre.MULTILINE)

	// matcher := regex.MatcherString(sexpr, 0)
	// for i := 1; i <= matcher.Groups(); i++ {
	// 	token := matcher.GroupString(i)
	// 	if len(token) > 0 {
	// 		tokens = append(tokens, token)
	// 	} else {
	// 		break
	// 	}
	// }

	// return tokens

	results := make([]string, 0, 1)
	// Work around lack of quoting in backtick
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	for _, group := range re.FindAllStringSubmatch(sexpr, -1) {
		if (group[1] == "") || (group[1][0] == ';') {
			continue
		}
		results = append(results, group[1])
	}

	return results
}

type MalType struct {
	_integer *int64
	_symbol  *string
	_list    *[]MalType
}

func readForm(r *reader) MalType {
	token := r.peek()

	if token != nil && *token == "(" {
		return readList(r)
	} else {
		return readAtom(r)
	}
}

func readList(r *reader) MalType {
	list := make([]MalType, 0)

	for token := r.next(); token != nil; token = r.next() {
		list = append(list, readForm(r))
	}

	return MalType{_list: &list}
}

func readAtom(r *reader) MalType {
	token := r.peek()

	if token != nil && *token != ")" {
		i, err := strconv.ParseInt(*token, 10, 64)
		if err == nil {
			return MalType{_integer: &i}
		}

		return MalType{_symbol: token}
	} else {
		return MalType{}
	}
}

func ReadStr(sexpr string) MalType {
	return readForm(&reader{
		position: 0,
		tokens:   tokenize(sexpr),
	})
}

func PrintStr(t MalType) string {
	if t._integer != nil {
		return fmt.Sprintf("%d", *t._integer)
	} else if t._symbol != nil {
		return *t._symbol
	} else if t._list != nil {
		tokens := make([]string, 0)
		for _, _type := range *t._list {
			token := PrintStr(_type)
			if len(token) > 0 {
				tokens = append(tokens, token)
			}
		}

		return fmt.Sprintf("(%s)", strings.Join(tokens, " "))
	} else {
		return ""
	}
}
