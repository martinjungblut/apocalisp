package typing

import (
	"testing"
)

func Test_IfBoolean_Must_Invoke_Callback_If_Boolean(t *testing.T) {
	values := []bool{false, true}

	for _, value := range values {
		called := false
		node := Type{Boolean: &value}

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

func Test_IfBoolean_Must_Not_Invoke_Callback_If_Not_Boolean(t *testing.T) {
	called := false
	node := Type{}

	node.IfBoolean(func(v bool) {
		called = true
	})

	if called != false {
		t.Error("IfBoolean() failed: called.")
	}
}

func Test_IsBoolean_Must_Return_True_If_Boolean_Value_Equals(t *testing.T) {
	values := []bool{false, true}

	for _, value := range values {
		node := Type{Boolean: &value}

		if !node.IsBoolean(value) {
			t.Error("IsBoolean() failed: wrong return value.")
		}
	}
}

func Test_IsBoolean_Must_Return_False_If_Not_Boolean(t *testing.T) {
	node := Type{}
	if node.IsBoolean(true) || node.IsBoolean(false) {
		t.Error("IsBoolean() failed: wrong return value.")
	}
}

func Test_NewBoolean_Creates_Boolean(t *testing.T) {
	values := []bool{false, true}

	for _, value := range values {
		node := NewBoolean(value)

		if !node.IsBoolean(value) {
			t.Error("NewBoolean() failed.")
		}
	}
}
