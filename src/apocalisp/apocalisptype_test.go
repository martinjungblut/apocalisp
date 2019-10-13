package apocalisp

import (
	"fmt"
	"testing"
)

func Test_ToString_NativeFunction(t *testing.T) {
	environment := DefaultEnvironment()
	functionNames := []string{"+", "-", "*", "/"}

	for i := range functionNames {
		if function, err := environment.Get(functionNames[i]); err == nil {
			if function.ToString() != "#<function>" {
				t.Error(fmt.Sprintf("ToString() returned unexpected value: `%s`.", function.ToString()))
			}
		} else {
			t.Error(err)
		}
	}
}

func Test_ToString_Symbol(t *testing.T) {
	symbols := []string{"+", "-", "*", "/"}

	for i := range symbols {
		node := ApocalispType{Symbol: &symbols[i]}

		if node.ToString() != symbols[i] {
			t.Error(fmt.Sprintf("ToString() returned unexpected value: `%s`.", node.ToString()))
		}
	}
}
