package core

import "testing"

func Test_IsKeyword(t *testing.T) {
	if !NewSymbol(":foo").IsKeyword() {
		t.Error("IsKeyword() failed.")
	}

	if !NewSymbol("::foo").IsKeyword() {
		t.Error("IsKeyword() failed.")
	}

	if NewSymbol("foo").IsKeyword() {
		t.Error("IsKeyword() failed.")
	}

	if NewString(":foo").IsKeyword() {
		t.Error("IsKeyword() failed.")
	}
}

func Test_ToKeyword(t *testing.T) {
	if converted, node := NewSymbol(":foo").ToKeyword(); !converted || node.AsSymbol() != ":foo" {
		t.Error("ToKeyword() failed.")
	}

	if converted, node := NewSymbol("foo").ToKeyword(); !converted || node.AsSymbol() != ":foo" {
		t.Error("ToKeyword() failed.")
	}

	if converted, node := NewString(":foo").ToKeyword(); !converted || node.AsSymbol() != ":foo" {
		t.Error("ToKeyword() failed.")
	}

	if converted, node := NewString("foo").ToKeyword(); !converted || node.AsSymbol() != ":foo" {
		t.Error("ToKeyword() failed.")
	}

	node := &Type{}
	if converted, _ := node.ToKeyword(); converted {
		t.Error("ToKeyword() failed.")
	}
}
