package core

import (
	"testing"
)

func Test_NewBoolean_Creates_Boolean(t *testing.T) {
	values := []bool{false, true}

	for _, value := range values {
		node := NewBoolean(value)

		if !node.IsBoolean() || node.AsBoolean() != value {
			t.Error("NewBoolean() failed.")
		}
	}
}

func Test_IsBoolean(t *testing.T) {
	node := &Type{}
	if node.IsBoolean() {
		t.Error("IsBoolean() failed: wrong return value.")
	}

	values := []bool{false, true}
	for _, value := range values {
		node := &Type{Boolean: &value}

		if !node.IsBoolean() {
			t.Error("IsBoolean() failed: wrong return value.")
		}
	}
}

func Test_AsBoolean(t *testing.T) {
	node := NewBoolean(false)
	if node.AsBoolean() != false {
		t.Error("AsBoolean() failed: wrong return value.")
	}

	node = NewBoolean(true)
	if node.AsBoolean() != true {
		t.Error("AsBoolean() failed: wrong return value.")
	}

	node = &Type{}
	if node.AsBoolean() != false {
		t.Error("AsBoolean() failed: wrong return value.")
	}
}

func Test_CompareBoolean(t *testing.T) {
	values := []bool{false, true}
	for _, value := range values {
		node := &Type{Boolean: &value}

		if !node.CompareBoolean(value) {
			t.Error("CompareBoolean() failed: wrong return value.")
		}
	}

	node := &Type{}
	if node.CompareBoolean(false) || node.CompareBoolean(true) {
		t.Error("CompareBoolean() failed: wrong return value.")
	}
}
