package core

import "testing"

func Test_IsNumber_Returns_True_For_Integers(t *testing.T) {
	val := int64(34)
	node := Type{Integer: &val}

	if !node.IsNumber() {
		t.Error("IsNumber() should return true for integers.")
	}
}

func Test_IsNumber_Returns_True_For_Floats(t *testing.T) {
	val := float64(34)
	node := Type{Float: &val}

	if !node.IsNumber() {
		t.Error("IsNumber() should return true for floats.")
	}
}

func Test_IsNumber_Returns_False_For_Other_Types(t *testing.T) {
	node := Type{}

	if node.IsNumber() {
		t.Error("IsNumber() should return false for other types.")
	}
}

func Test_AsNumber_Returns_Numeric_Value_For_Integers(t *testing.T) {
	val := int64(34)
	node := Type{Integer: &val}

	if node.AsNumber() != 34.0 {
		t.Error("AsNumber() failed.")
	}
}

func Test_AsNumber_Returns_Numeric_Value_For_Floats(t *testing.T) {
	val := float64(34.5)
	node := Type{Float: &val}

	if node.AsNumber() != 34.5 {
		t.Error("AsNumber() failed.")
	}
}

func Test_AsNumber_Returns_0_For_Other_Types(t *testing.T) {
	node := Type{}

	if node.AsNumber() != 0 {
		t.Error("AsNumber() failed.")
	}
}
