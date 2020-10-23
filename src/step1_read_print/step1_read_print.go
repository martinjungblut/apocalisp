package main

import (
	"apocalisp"
	"apocalisp/parser"
)

func main() {
	apocalisp.Repl(apocalisp.NoEval, parser.Parser{})
}
