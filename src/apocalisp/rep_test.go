package apocalisp

import (
	"testing"
)

func Repl_Test(in string, eout string, t *testing.T) {
	eval := Step4Eval
	environment := DefaultEnvironment()

	if out, err := Rep(in, environment, eval); err != nil && err.Error() != eout {
		t.Errorf("(output) `%s` != `%s` (expected)", err.Error(), eout)
	} else if out != eout {
		t.Errorf("(output) `%s` != `%s` (expected)", out, eout)
	}
}

func Test_String_Parsing_And_Evaluation(t *testing.T) {
	Repl_Test("\"test\"", "\"test\"", t)
}

func Test_Lambda_Works_As_Function(t *testing.T) {
	Repl_Test("((λ (a b) (+ a b)) 3 4)", "7", t)
	Repl_Test("((λ (a b) (+ a b)) 2 1)", "3", t)
	Repl_Test("((λ (a b) (- a b)) 10 9)", "1", t)
}
