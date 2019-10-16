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
