package manalisp

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

	readAheadPosition int
}

func (r *reader) readAhead() error {
	reachedEnd := func() bool { return r.readAheadPosition >= len(r.tokens) }
	currentToken := func() string { return r.tokens[r.readAheadPosition] }
	unclosedString := func(token string) bool {
		return strings.HasPrefix(token, "\"") && !strings.HasSuffix(token, "\"")
	}

	if !reachedEnd() {
		switch token := currentToken(); token {
		case "(":
			r.parensCount++
		case ")":
			r.parensCount--
		case "[":
			r.bracketsCount++
		case "]":
			r.bracketsCount--
		case "{":
			r.bracesCount++
		case "}":
			r.bracesCount--
		default:
			if unclosedString(token) {
				return errors.New("unexpected EOF")
			}
		}
	}

	if r.parensCount < 0 {
		return errors.New("unexpected ')'")
	} else if r.bracketsCount < 0 {
		return errors.New("unexpected ']'")
	} else if r.bracesCount < 0 {
		return errors.New("unexpected '}'")
	}
	if reachedEnd() && (r.parensCount > 0 || r.bracketsCount > 0 || r.bracesCount > 0) {
		return errors.New("unexpected EOF")
	}

	r.readAheadPosition++
	return nil
}

func (r *reader) next() (*string, error) {
	if err := r.readAhead(); err != nil {
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

func readForm(r *reader) (*ManalispType, error) {
	token, err := r.next()
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, nil
	}

	if *token == "^" {
		if list, err := readList(r); err == nil {
			symbol := "with-meta"
			subelements := *list.List
			seq := []ManalispType{ManalispType{Symbol: &symbol}, subelements[1], subelements[0]}
			return &ManalispType{List: &seq}, nil
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	}

	if *token == "(" {
		return readList(r)
	} else if *token == "[" {
		return readVector(r)
	} else if *token == "{" {
		return readHashmap(r)
	} else if *token == "'" {
		return readPrefixExpansion(r, "quote")
	} else if *token == "~" {
		return readPrefixExpansion(r, "unquote")
	} else if *token == "`" {
		return readPrefixExpansion(r, "quasiquote")
	} else if *token == "@" {
		return readPrefixExpansion(r, "deref")
	} else if *token == "~@" {
		return readPrefixExpansion(r, "splice-unquote")
	} else if *token != ")" && *token != "]" && *token != "}" {
		return readAtom(token)
	}
	return nil, nil
}

func readSequence(r *reader) (*[]ManalispType, error) {
	sequence := []ManalispType{}
	for form, err := readForm(r); form != nil || err != nil; form, err = readForm(r) {
		if err != nil {
			return nil, err
		} else if form != nil {
			sequence = append(sequence, *form)
		}
	}
	return &sequence, nil
}

func readAtom(token *string) (*ManalispType, error) {
	i, err := strconv.ParseInt(*token, 10, 64)
	if err == nil {
		return &ManalispType{Integer: &i}, nil
	}
	return &ManalispType{Symbol: token}, nil
}

func readList(r *reader) (*ManalispType, error) {
	sequence, err := readSequence(r)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{List: sequence}, nil
	}
}

func readVector(r *reader) (*ManalispType, error) {
	sequence, err := readSequence(r)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{Vector: sequence}, nil
	}
}

func readHashmap(r *reader) (*ManalispType, error) {
	sequence, err := readSequence(r)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{Hashmap: sequence}, nil
	}
}

func readPrefixExpansion(r *reader, symbol string) (*ManalispType, error) {
	if form, err := readForm(r); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []ManalispType{ManalispType{Symbol: &symbol}, *form}
		return &ManalispType{List: &sequence}, nil
	} else {
		return nil, nil
	}
}

func tokenize(sexpr string) []string {
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

func ReadStr(sexpr string) (*ManalispType, error) {
	return readForm(&reader{tokens: tokenize(sexpr)})
}
