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

	parensCount   int
	bracketsCount int
	bracesCount   int

	foundString bool

	readAheadPosition int
}

func (r *reader) readAhead() error {
	if r.readAheadPosition < len(r.tokens) {
		token := r.tokens[r.readAheadPosition]

		if token == "(" {
			r.parensCount++
		} else if token == ")" {
			r.parensCount--
		} else if token == "[" {
			r.bracketsCount++
		} else if token == "]" {
			r.bracketsCount--
		} else if token == "{" {
			r.bracesCount++
		} else if token == "}" {
			r.bracesCount--
		} else if token == "\"" {
			r.foundString = !r.foundString
		} else if strings.HasPrefix(token, "\"") && !strings.HasSuffix(token, "\"") {
			return errors.New("unexpected EOF")
		}
	} else if r.foundString {
		return errors.New("unexpected EOF")
	}

	if r.parensCount < 0 {
		return errors.New("unexpected ')'")
	}
	if (r.readAheadPosition == len(r.tokens)) && r.parensCount > 0 {
		return errors.New("unexpected EOF")
	}

	if r.bracketsCount < 0 {
		return errors.New("unexpected ']'")
	}
	if (r.readAheadPosition == len(r.tokens)) && r.bracketsCount > 0 {
		return errors.New("unexpected EOF")
	}

	if r.bracesCount < 0 {
		return errors.New("unexpected '}'")
	}
	if (r.readAheadPosition == len(r.tokens)) && r.bracesCount > 0 {
		return errors.New("unexpected EOF")
	}

	r.readAheadPosition++
	return nil
}

func (r *reader) next() (*string, error) {
	err := r.readAhead()
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

type MalType struct {
	_integer *int64
	_symbol  *string
	_list    *[]MalType
	_vector  *[]MalType
	_hashmap *[]MalType
}

func readForm(r *reader) (*MalType, error) {
	token, err := r.next()
	if err != nil {
		return nil, err
	}

	if token != nil && *token == "(" {
		return readList(r)
	} else if token != nil && *token == "[" {
		return readVector(r)
	} else if token != nil && *token == "{" {
		return readHashmap(r)
	} else if token != nil && *token == "'" {
		return readPrefixExpansion(r, "quote")
	} else if token != nil && *token == "~" {
		return readPrefixExpansion(r, "unquote")
	} else if token != nil && *token == "`" {
		return readPrefixExpansion(r, "quasiquote")
	} else if token != nil && *token == "@" {
		return readPrefixExpansion(r, "deref")
	} else if token != nil && *token == "~@" {
		return readPrefixExpansion(r, "splice-unquote")
	} else if token != nil && *token != ")" && *token != "]" && *token != "}" {
		return readAtom(token)
	} else {
		return nil, nil
	}
}

func readSequence(r *reader) (*[]MalType, error) {
	sequence := make([]MalType, 0)

	for {
		form, err := readForm(r)
		if err != nil {
			return nil, err
		} else if form != nil {
			sequence = append(sequence, *form)
		} else {
			break
		}
	}

	return &sequence, nil
}

func readAtom(token *string) (*MalType, error) {
	i, err := strconv.ParseInt(*token, 10, 64)
	if err == nil {
		return &MalType{_integer: &i}, nil
	}

	return &MalType{_symbol: token}, nil
}

func readList(r *reader) (*MalType, error) {
	sequence, err := readSequence(r)

	if err != nil {
		return nil, err
	} else {
		return &MalType{_list: sequence}, nil
	}
}

func readVector(r *reader) (*MalType, error) {
	sequence, err := readSequence(r)

	if err != nil {
		return nil, err
	} else {
		return &MalType{_vector: sequence}, nil
	}
}

func readHashmap(r *reader) (*MalType, error) {
	sequence, err := readSequence(r)

	if err != nil {
		return nil, err
	} else {
		return &MalType{_hashmap: sequence}, nil
	}
}

func readPrefixExpansion(r *reader, symbol string) (*MalType, error) {
	form, err := readForm(r)
	if err != nil {
		return nil, err
	}

	sequence := make([]MalType, 0)
	sequence = append(sequence, MalType{_symbol: &symbol})
	sequence = append(sequence, *form)

	return &MalType{_list: &sequence}, nil
}

func tokenize(sexpr string) []string {
	rawTokens := make([]string, 0)
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	for _, group := range re.FindAllStringSubmatch(sexpr, -1) {
		if (group[1] == "") || (group[1][0] == ';') {
			continue
		}
		rawTokens = append(rawTokens, group[1])
	}

	tokens := make([]string, 0)
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

func ReadStr(sexpr string) (*MalType, error) {
	return readForm(&reader{
		position:          0,
		readAheadPosition: 0,
		tokens:            tokenize(sexpr),
		foundString:       false,
	})
}

func PrintStr(t *MalType) string {
	seqToStr := func(seq *[]MalType, lChar string, rChar string) string {
		tokens := make([]string, 0)

		for _, maltype := range *seq {
			token := PrintStr(&maltype)
			if len(token) > 0 {
				tokens = append(tokens, token)
			}
		}

		return fmt.Sprintf("%s%s%s", lChar, strings.Join(tokens, " "), rChar)
	}

	if t != nil {
		if t._integer != nil {
			return fmt.Sprintf("%d", *t._integer)
		} else if t._symbol != nil {
			return *t._symbol
		} else if t._list != nil {
			return seqToStr(t._list, "(", ")")
		} else if t._vector != nil {
			return seqToStr(t._vector, "[", "]")
		} else if t._hashmap != nil {
			return seqToStr(t._hashmap, "{", "}")
		} else {
			return ""
		}
	} else {
		return ""
	}
}
