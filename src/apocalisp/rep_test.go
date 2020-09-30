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

func Test_Signed_Float_Support_Mathematical_Expressions(t *testing.T) {
	// note numbers are coerced.
	// integers are preferred, if there is no precision loss.

	Repl_Test(`(+ 1.0 1.0)`, `2`, t)
	Repl_Test(`(+ 1.0 1)`, `2`, t)
	Repl_Test(`(+ 1 1.0)`, `2`, t)

	Repl_Test(`(- 1.0 1)`, `0`, t)
	Repl_Test(`(- 0 2.2)`, `-2.200000`, t)

	Repl_Test(`(* 2.0 3.0)`, `6`, t)
	Repl_Test(`(* 5 -3.0)`, `-15`, t)

	Repl_Test(`(/ 3 2)`, `1.500000`, t)
	Repl_Test(`(/ -3 -3.0)`, `1`, t)
}

func Test_Signed_Float_Support_Comparison_Expressions(t *testing.T) {
	Repl_Test(`(< 5 5.01)`, `true`, t)
	Repl_Test(`(< 5.01 5)`, `false`, t)
	Repl_Test(`(< 5 10.0)`, `true`, t)
	Repl_Test(`(< 10.0 5)`, `false`, t)

	Repl_Test(`(<= 10 10.0)`, `true`, t)
	Repl_Test(`(<= 10.0 10)`, `true`, t)
	Repl_Test(`(<= 5 5.01)`, `true`, t)
	Repl_Test(`(<= 5.01 5)`, `false`, t)

	Repl_Test(`(> 5.01 5)`, `true`, t)
	Repl_Test(`(> 5 5.01)`, `false`, t)
	Repl_Test(`(> 10.0 5)`, `true`, t)
	Repl_Test(`(> 5 10.0)`, `false`, t)

	Repl_Test(`(>= 10.0 10)`, `true`, t)
	Repl_Test(`(>= 10 10.0)`, `true`, t)
	Repl_Test(`(>= 5.01 5)`, `true`, t)
	Repl_Test(`(>= 5 5.01)`, `false`, t)
}

func Test_Signed_Float_Support_Equality_Expressions(t *testing.T) {
	Repl_Test(`(= 5 5.0)`, `true`, t)
	Repl_Test(`(= 0 0.0)`, `true`, t)
	Repl_Test(`(= -1 -1.0)`, `true`, t)
}
