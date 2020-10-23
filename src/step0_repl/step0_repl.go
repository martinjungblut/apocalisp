package main

import (
	"apocalisp"
	"apocalisp/core"
)

type parser struct {
}

func (p parser) Parse(sexpr string) (*core.Type, error) {
	return core.NewSymbol(sexpr), nil
}

func main() {
	apocalisp.Repl(apocalisp.NoEval, parser{})
}
