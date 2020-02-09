package apocalisp

import (
	"apocalisp/core"
	"apocalisp/parser"
	"testing"
)

func Repl_Test(in string, eout string, t *testing.T) {
	eval := Step6Eval
	environment := core.DefaultEnvironment(parser.Parser{})

	if out, err := Rep(in, environment, eval); err != nil && err.Error() != eout {
		t.Errorf("(output) `%s` != `%s` (expected)", err.Error(), eout)
	} else if out != eout {
		t.Errorf("(output) `%s` != `%s` (expected)", out, eout)
	}
}

func Test_String_Parsing_And_Evaluation(t *testing.T) {
	Repl_Test("\"test\"", "\"test\"", t)
}

func Test_Alternative_Function_Notation(t *testing.T) {
	Repl_Test(`((\ (a b) (+ a b)) 3 4)`, `7`, t)
	Repl_Test(`((\ (a b) (+ a b)) 2 1)`, `3`, t)
	Repl_Test(`((\ (a b) (- a b)) 10 9)`, `1`, t)
}
