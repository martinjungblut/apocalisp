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

func Test_ToString_String(t *testing.T) {
	strings := []string{"first", "second"}

	for i := range strings {
		node := ApocalispType{String: &strings[i]}

		if node.ToString() != strings[i] {
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

func Test_IsString_Returns_True_If_String(t *testing.T) {
	value := "value"
	node := ApocalispType{String: &value}
	if !node.IsString() {
		t.Error("IsString() failed.")
	}
}

func Test_IsString_Returns_False_If_Not_String(t *testing.T) {
	node := ApocalispType{}
	if node.IsString() {
		t.Error("IsString() failed.")
	}
}

func Test_IsString_Returns_False_If_Symbol(t *testing.T) {
	value := "value"
	node := ApocalispType{Symbol: &value}
	if node.IsString() {
		t.Error("IsString() failed.")
	}
}

func Test_IsSymbol_Returns_True_If_Symbol(t *testing.T) {
	value := "value"
	node := ApocalispType{Symbol: &value}
	if !node.IsSymbol() {
		t.Error("IsSymbol() failed.")
	}
}

func Test_IsSymbol_Returns_False_If_Not_Symbol(t *testing.T) {
	node := ApocalispType{}
	if node.IsSymbol() {
		t.Error("IsSymbol() failed.")
	}
}

func Test_IsSymbol_Returns_False_If_String(t *testing.T) {
	value := "value"
	node := ApocalispType{String: &value}
	if node.IsSymbol() {
		t.Error("IsSymbol() failed.")
	}
}

func Test_IfBoolean_MustInvokeCallback_IfBoolean(t *testing.T) {
	values := []bool{false, true}

	for _, value := range values {
		called := false
		node := ApocalispType{Boolean: &value}

		node.IfBoolean(func(v bool) {
			called = true
			if v != value {
				t.Error("IfBoolean() failed: inconsistent boolean value.")
			}
		})

		if called != true {
			t.Error("IfBoolean() failed: not called.")
		}
	}
}

func Test_IfBoolean_MustNotInvokeCallback_IfNotBoolean(t *testing.T) {
	called := false
	node := ApocalispType{}

	node.IfBoolean(func(v bool) {
		called = true
	})

	if called != false {
		t.Error("IfBoolean() failed: called.")
	}
}

func Test_IsBoolean_MustReturnTrue_IfBooleanValueIsTheSame(t *testing.T) {
	values := []bool{false, true}

	for _, value := range values {
		node := ApocalispType{Boolean: &value}

		if !node.IsBoolean(value) {
			t.Error("IsBoolean() failed: wrong return value.")
		}
	}

}

func Test_IsBoolean_MustReturnFalse_IfNotBoolean(t *testing.T) {
	node := ApocalispType{}
	if node.IsBoolean(true) || node.IsBoolean(false) {
		t.Error("IsBoolean() failed: wrong return value.")
	}
}
