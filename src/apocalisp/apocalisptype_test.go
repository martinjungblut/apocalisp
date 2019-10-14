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

func Test_NewNil_Creates_Nil(t *testing.T) {
	node := NewNil()
	if !node.IsNil() {
		t.Error("NewNil() failed.")
	}
}

func Test_IsNil_Returns_True_If_Nil(t *testing.T) {
	node := ApocalispType{Nil: true}
	if !node.IsNil() {
		t.Error("IsNil() failed.")
	}
}

func Test_IsNil_Returns_False_If_Not_Nil(t *testing.T) {
	node := ApocalispType{Nil: false}
	if node.IsNil() {
		t.Error("IsNil() failed.")
	}
}

func Test_IsFalse_Returns_True_If_False_Boolean(t *testing.T) {
	value := false
	node := ApocalispType{Boolean: &value}
	if !node.IsFalse() {
		t.Error("IsFalse() failed.")
	}
}

func Test_IsFalse_Returns_False_If_True_Boolean(t *testing.T) {
	value := true
	node := ApocalispType{Boolean: &value}
	if node.IsFalse() {
		t.Error("IsFalse() failed.")
	}
}

func Test_IsFalse_Returns_False_If_Not_Boolean(t *testing.T) {
	node := ApocalispType{}
	if node.IsFalse() {
		t.Error("IsFalse() failed.")
	}
}
