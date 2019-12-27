package core

import "testing"

func Test_NewNil_Creates_Nil(t *testing.T) {
	node := NewNil()
	if !node.IsNil() {
		t.Error("NewNil() failed.")
	}
}

func Test_IsNil_Returns_True_If_Nil(t *testing.T) {
	node := Type{Nil: true}
	if !node.IsNil() {
		t.Error("IsNil() failed.")
	}
}

func Test_IsNil_Returns_False_If_Not_Nil(t *testing.T) {
	node := Type{Nil: false}
	if node.IsNil() {
		t.Error("IsNil() failed.")
	}
}
