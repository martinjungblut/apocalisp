package main

import (
	"apocalisp"
	"apocalisp/parser"
)

func main() {
	apocalisp.Repl(apocalisp.Evaluate, parser.Parser{})
}
