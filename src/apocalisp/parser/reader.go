package parser

import (
	"errors"
	"strings"
)

type Reader struct {
	position int
	tokens   []string

	parensCount   int
	bracketsCount int
	bracesCount   int

	readAheadPosition int
}

func NewReader(tokens []string) *Reader {
	reader := Reader{tokens: tokens}
	return &reader
}

func (r *Reader) Next() (*string, error) {
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

func (r *Reader) readAhead() error {
	reachedEnd := func() bool { return r.readAheadPosition == len(r.tokens) }
	currentToken := func() string { return r.tokens[r.readAheadPosition] }
	unclosedString := func(token string) bool {
		return strings.HasPrefix(token, "\"") && (len(token) == 1 || !strings.HasSuffix(token, "\""))
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
				return errors.New("Error: unexpected EOF.")
			}
		}
	}

	if r.parensCount < 0 {
		return errors.New("Error: unexpected ')'.")
	} else if r.bracketsCount < 0 {
		return errors.New("Error: unexpected ']'.")
	} else if r.bracesCount < 0 {
		return errors.New("Error: unexpected '}'.")
	} else if reachedEnd() && (r.parensCount > 0 || r.bracketsCount > 0 || r.bracesCount > 0) {
		return errors.New("Error: unexpected EOF.")
	}

	r.readAheadPosition++
	return nil
}
