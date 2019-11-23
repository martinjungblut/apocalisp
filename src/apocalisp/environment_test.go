package apocalisp

import (
	"apocalisp/typing"
	"testing"
)

func Test_NewEnvironment_Should_Add_Bindings_To_Environment_More_Symbols_Than_Nodes(t *testing.T) {
	firstValue := "firstValue"
	secondValue := "secondValue"
	firstNode := typing.Type{Symbol: &firstValue}
	secondNode := typing.Type{Symbol: &secondValue}

	environment := NewEnvironment(nil, []string{"a", "b", "c"}, []typing.Type{firstNode, secondNode})

	if node, err := environment.Get("a"); err == nil {
		if !node.IsSymbol() || node.AsSymbol() != "firstValue" {
			t.Error("Symbol not set when calling NewEnvironment().")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}

	if node, err := environment.Get("b"); err == nil {
		if !node.IsSymbol() || node.AsSymbol() != "secondValue" {
			t.Error("Symbol not set when calling NewEnvironment().")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}

	if _, err := environment.Get("c"); err == nil {
		t.Error("Symbol should not have been set, but was.")
	}
}

func Test_NewEnvironment_Should_Add_Bindings_To_Environment_More_Nodes_Than_Symbols(t *testing.T) {
	firstValue := "firstValue"
	secondValue := "secondValue"
	firstNode := typing.Type{Symbol: &firstValue}
	secondNode := typing.Type{Symbol: &secondValue}

	environment := NewEnvironment(nil, []string{"a"}, []typing.Type{firstNode, secondNode})

	if node, err := environment.Get("a"); err == nil {
		if !node.IsSymbol() || node.AsSymbol() != "firstValue" {
			t.Error("Symbol not set when calling NewEnvironment().")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}

	if _, err := environment.Get("b"); err == nil {
		t.Error("Symbol should not have been set, but was.")
	}

	if _, err := environment.Get("c"); err == nil {
		t.Error("Symbol should not have been set, but was.")
	}
}

func Test_NewEnvironment_Should_Support_Variadic_Parameters(t *testing.T) {
	firstValue := "firstValue"
	secondValue := "secondValue"
	firstNode := typing.Type{String: &firstValue}
	secondNode := typing.Type{String: &secondValue}

	environment := NewEnvironment(nil, []string{"&", "other"}, []typing.Type{firstNode, secondNode})

	if node, err := environment.Get("other"); err == nil {
		if !node.IsList() {
			t.Error("Symbol should have been set as a list.")
		}

		l := node.AsList()
		if len(l) != 2 {
			t.Error("Incorrect list length.")
		}

		if node.AsList()[0].AsString() != firstValue {
			t.Error("Value mismatch.")
		}

		if node.AsList()[1].AsString() != secondValue {
			t.Error("Value mismatch.")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}
}

func Test_NewEnvironment_Should_Support_Variadic_Parameters_Falling_Back_To_A_Safe_Symbol_If_None_Is_Provided(t *testing.T) {
	firstValue := "firstValue"
	secondValue := "secondValue"
	firstNode := typing.Type{String: &firstValue}
	secondNode := typing.Type{String: &secondValue}

	environment := NewEnvironment(nil, []string{"a", "&"}, []typing.Type{firstNode, secondNode})

	if node, err := environment.Get("&"); err == nil {
		if !node.IsList() {
			t.Error("Symbol should have been set as a list.")
		}

		l := node.AsList()
		if len(l) != 1 {
			t.Error("Incorrect list length.")
		}

		if node.AsList()[0].AsString() != secondValue {
			t.Error("Value mismatch.")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}
}

func Test_NewEnvironment_Should_Set_Symbol_As_Empty_List_If_No_Variadic_Arguments_Are_Provided(t *testing.T) {
	environment := NewEnvironment(nil, []string{"&", "other"}, []typing.Type{})
	if node, err := environment.Get("other"); err == nil {
		if !node.IsList() {
			t.Error("Symbol should have been set as a list.")
		}

		if len(node.AsList()) != 0 {
			t.Error("Incorrect list length.")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}

	environment = NewEnvironment(nil, []string{"&"}, []typing.Type{})
	if node, err := environment.Get("&"); err == nil {
		if !node.IsList() {
			t.Error("Symbol should have been set as a list.")
		}

		if len(node.AsList()) != 0 {
			t.Error("Incorrect list length.")
		}
	} else {
		t.Error("Symbol not set when calling NewEnvironment().")
	}
}

func Test_Set_Should_Add_Binding_To_Environment(t *testing.T) {
	environment := NewEnvironment(nil, []string{}, []typing.Type{})

	emptyType := typing.Type{}

	environment.Set("a", emptyType)

	if fetched, err := environment.Get("a"); err != nil {
		t.Error("Get() should not have returned error.")
	} else if fetched != emptyType {
		t.Error("Set() failed.")
	}
}

func Test_Get_Should_Return_Bindings_From_Outer_Environments(t *testing.T) {
	environment1 := NewEnvironment(nil, []string{}, []typing.Type{})
	environment2 := NewEnvironment(environment1, []string{}, []typing.Type{})
	environment3 := NewEnvironment(environment2, []string{}, []typing.Type{})

	emptyType := typing.Type{}

	environment1.Set("a", emptyType)

	if fetched, err := environment3.Get("a"); err != nil {
		t.Error("Get() should not have returned error.")
	} else if fetched != emptyType {
		t.Error("Set() failed.")
	}
}

func Test_Get_Should_Return_Error_If_Symbol_Not_Found(t *testing.T) {
	environment := NewEnvironment(nil, []string{}, []typing.Type{})

	_, err := environment.Get("a")
	if err == nil {
		t.Error("Get() should have returned error.")
	}
}

func Test_Find_Should_Return_Environment_Containing_Binding(t *testing.T) {
	environment1 := NewEnvironment(nil, []string{}, []typing.Type{})
	environment2 := NewEnvironment(environment1, []string{}, []typing.Type{})
	environment3 := NewEnvironment(environment2, []string{}, []typing.Type{})

	emptyType := typing.Type{}

	environment1.Set("a", emptyType)

	if environment3.Find("a") != environment1 {
		t.Error("Find() failed.")
	}
}

func Test_Find_Should_Return_Nil_If_No_Environment_Contains_Binding(t *testing.T) {
	environment1 := NewEnvironment(nil, []string{}, []typing.Type{})
	environment2 := NewEnvironment(environment1, []string{}, []typing.Type{})
	environment3 := NewEnvironment(environment2, []string{}, []typing.Type{})

	if environment3.Find("a") != nil {
		t.Error("Find() failed.")
	}

	if environment2.Find("a") != nil {
		t.Error("Find() failed.")
	}

	if environment1.Find("a") != nil {
		t.Error("Find() failed.")
	}
}

func Test_SetCallable_Should_Add_Native_Function_Binding_To_Environment(t *testing.T) {
	environment := NewEnvironment(nil, []string{}, []typing.Type{})

	environment.SetCallable("print", func(...typing.Type) typing.Type {
		return typing.Type{}
	})

	if fetched, err := environment.Get("print"); err != nil {
		t.Error("Get() should not have returned error.")
	} else if fetched.Callable == nil || fetched.Symbol == nil {
		t.Error("SetCallable() failed.")
	}
}
