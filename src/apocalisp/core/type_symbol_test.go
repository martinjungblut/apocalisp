package core

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

func Test_CompareSymbol_Returns_False_If_Not_Symbol(t *testing.T) {
	node := Type{}

	if node.CompareSymbol("") != false {
		t.Error("CompareSymbol() failed.")
	}
}

func Test_CompareSymbol_Returns_False_If_Symbols_Differ(t *testing.T) {
	symbols := map[string]string{
		"":     " ",
		" ":    "  ",
		"john": "JOHN",
	}

	for key, value := range symbols {
		node := Type{Symbol: &key}

		if node.CompareSymbol(value) != false {
			t.Errorf("CompareSymbol() failed. Key: '%s'. Value: '%s'.", key, value)
		}
	}
}

func Test_CompareSymbol_Returns_True_If_Symbols_Equal(t *testing.T) {
	symbols := []string{
		"",
		" ", "  ",
		"john", "JOHN",
	}

	for _, value := range symbols {
		node := Type{Symbol: &value}

		if node.CompareSymbol(value) != true {
			t.Errorf("CompareSymbol() failed. Value: '%s'.", value)
		}
	}
}

func Test_CompareSymbol_Returns_True_If_Symbol_Equals_Any(t *testing.T) {
	symbols := []string{"a", "b", "c"}

	for _, symbol := range symbols {
		node := Type{Symbol: &symbol}
		if node.CompareSymbol(symbols...) != true {
			t.Error("CompareSymbol() failed.")
		}
	}
}
