package parser

import (
	"errors"
	"strings"
)

type reader struct {
	position          int
	readAheadPosition int
	readAheadCalled   bool
	tokens            []string
	parensCount       int
	bracketsCount     int
	bracesCount       int
}

func newReader(tokens []string) *reader {
	reader := reader{tokens: tokens}
	return &reader
}

func (r *reader) next() (*string, error) {
	if !r.readAheadCalled {
		r.readAheadCalled = true
		if err := r.readAhead(); err != nil {
			return nil, err
		}
	}

	if r.position < len(r.tokens) {
		token := &(r.tokens[r.position])
		r.position++
		return token, nil
	} else {
		return nil, nil
	}
}

func (r *reader) readAhead() error {
	reachedEnd := func() bool { return r.readAheadPosition == len(r.tokens) }
	currentToken := func() string { return r.tokens[r.readAheadPosition] }
	unclosedString := func(token string) bool {
		return strings.HasPrefix(token, "\"") && (len(token) == 1 || !strings.HasSuffix(token, "\""))
	}

	for !reachedEnd() {
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
		r.readAheadPosition++
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

	return nil
}
