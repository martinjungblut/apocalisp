package manalispcore

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type reader struct {
	position int
	tokens   []string

	parensCount    int
	parensPosition int

	bracesCount    int
	bracesPosition int
}

func (r *reader) check() error {
	if r.parensPosition < len(r.tokens) {
		token := r.tokens[r.parensPosition]

		if token == "(" {
			r.parensCount++
		} else if token == ")" {
			r.parensCount--
		} else if token == "[" {
			r.bracesCount++
		} else if token == "]" {
			r.bracesCount--
		}
	}

	if r.parensCount < 0 {
		return errors.New("unexpected ')'")
	}
	if (r.parensPosition == len(r.tokens)) && r.parensCount > 0 {
		return errors.New("unexpected EOF")
	}

	if r.bracesCount < 0 {
		return errors.New("unexpected ']'")
	}
	if (r.bracesPosition == len(r.tokens)) && r.bracesCount > 0 {
		return errors.New("unexpected EOF")
	}

	r.parensPosition++
	r.bracesPosition++
	return nil
}

func (r *reader) peek() (*string, error) {
	err := r.check()
	if err != nil {
		return nil, err
	}

	if r.position < len(r.tokens) {
		token := &(r.tokens[r.position])
		return token, nil
	} else {
		return nil, nil
	}
}

func (r *reader) next() (*string, error) {
	err := r.check()
	if err != nil {
		return nil, err
	}

	if r.position < len(r.tokens) {
		token := &(r.tokens[r.position])
		r.position++
		return token, nil
	} else {
		return nil, nil
	}
}

func tokenize(sexpr string) []string {
	results := make([]string, 0)

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

func readForm(r *reader) (MalType, error) {
	token, err := r.peek()
	if err != nil {
		return MalType{}, err
	}

	if token != nil && (*token == "(" || *token == "[") {
		return readList(r)
	} else {
		return readAtom(r)
	}
}

func readList(r *reader) (MalType, error) {
	list := make([]MalType, 0)

	for {
		token, err := r.next()
		if err != nil {
			return MalType{}, err
		} else if token == nil {
			break
		} else if token != nil {
			t, err := readForm(r)
			if err != nil {
				return MalType{}, err
			} else {
				list = append(list, t)
			}
		}
	}

	return MalType{_list: &list}, nil
}

func readAtom(r *reader) (MalType, error) {
	token, err := r.peek()
	if err != nil {
		return MalType{}, err
	}

	if token != nil && (*token != ")" && *token != "]") {
		i, err := strconv.ParseInt(*token, 10, 64)
		if err == nil {
			return MalType{_integer: &i}, nil
		}

		return MalType{_symbol: token}, nil
	} else {
		return MalType{}, nil
	}
}

func ReadStr(sexpr string) (MalType, error) {
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
