package apocalisp

import (
	"apocalisp/parser"
	"testing"
)

func Repl_Test(in string, eout string, t *testing.T) {
	environment := DefaultEnvironment(parser.Parser{}, Evaluate)

	if out, err := Rep(in, environment, Evaluate); err != nil && err.Error() != eout {
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

func Test_Signed_Float_Support(t *testing.T) {
	Repl_Test(`(+ 1.0 1.0)`, `2.000000`, t)
	Repl_Test(`(+ 1.0 1)`, `2.000000`, t)
	Repl_Test(`(+ 1 1.0)`, `2.000000`, t)
	Repl_Test(`(- 1.0 1)`, `0.000000`, t)
	Repl_Test(`(- 0 2.2)`, `-2.200000`, t)
	Repl_Test(`(* 2.0 3.0)`, `6.000000`, t)
	Repl_Test(`(* 5 -3.0)`, `-15.000000`, t)
	Repl_Test(`(/ 3 3)`, `1.000000`, t)
	Repl_Test(`(/ 3 2)`, `1.500000`, t)
	Repl_Test(`(/ -3 -3)`, `1.000000`, t)
}
