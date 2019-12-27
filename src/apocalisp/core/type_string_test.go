package core

import "testing"

func Test_IsString_Returns_True_If_String(t *testing.T) {
	value := "value"
	node := Type{String: &value}

	if !node.IsString() {
		t.Error("IsString() failed.")
	}
}

func Test_IsString_Returns_False_If_Not_String(t *testing.T) {
	node := Type{}

	if node.IsString() {
		t.Error("IsString() failed.")
	}
}

func Test_IsString_Returns_False_If_Symbol(t *testing.T) {
	value := "value"
	node := Type{Symbol: &value}

	if node.IsString() {
		t.Error("IsString() failed.")
	}
}

func Test_AsString_Returns_Value_If_String(t *testing.T) {
	values := []string{"firstValue", "secondValue"}

	for _, value := range values {
		node := Type{String: &value}

		if node.AsString() != value {
			t.Error("AsString() failed.")
		}
	}
}

func Test_AsString_Returns_Empty_String_If_Not_String(t *testing.T) {
	node := Type{}

	if node.AsString() != "" {
		t.Error("AsString() failed.")
	}
}
