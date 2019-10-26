package typing

import (
	"testing"
)

func Test_IsSymbol_Returns_True_If_Symbol(t *testing.T) {
	value := "value"
	node := Type{Symbol: &value}
	if !node.IsSymbol() {
		t.Error("IsSymbol() failed.")
	}
}

func Test_IsSymbol_Returns_False_If_Not_Symbol(t *testing.T) {
	node := Type{}
	if node.IsSymbol() {
		t.Error("IsSymbol() failed.")
	}
}

func Test_IsSymbol_Returns_False_If_String(t *testing.T) {
	value := "value"
	node := Type{String: &value}
	if node.IsSymbol() {
		t.Error("IsSymbol() failed.")
	}
}

func Test_AsSymbol_Returns_Value_If_Symbol(t *testing.T) {
	values := []string{"firstValue", "secondValue"}

	for _, value := range values {
		node := Type{Symbol: &value}

		if node.AsSymbol() != value {
			t.Error("AsSymbol() failed.")
		}
	}
}

func Test_AsSymbol_Returns_Empty_Symbol_If_Not_Symbol(t *testing.T) {
	node := Type{}

	if node.AsSymbol() != "" {
		t.Error("AsSymbol() failed.")
	}
}
