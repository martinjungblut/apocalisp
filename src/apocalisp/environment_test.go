package apocalisp

import "testing"

func Test_NewEnvironment_Should_Add_Bindings_To_Environment_More_Symbols_Than_Nodes(t *testing.T) {
	firstValue := "firstValue"
	secondValue := "secondValue"
	firstNode := ApocalispType{Symbol: &firstValue}
	secondNode := ApocalispType{Symbol: &secondValue}

	environment := NewEnvironment(nil, []string{"a", "b", "c"}, []ApocalispType{firstNode, secondNode})
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
	firstNode := ApocalispType{Symbol: &firstValue}
	secondNode := ApocalispType{Symbol: &secondValue}

	environment := NewEnvironment(nil, []string{"a"}, []ApocalispType{firstNode, secondNode})
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

func Test_Set_Should_Add_Binding_To_Environment(t *testing.T) {
	environment := NewEnvironment(nil, []string{}, []ApocalispType{})

	emptyType := ApocalispType{}

	environment.Set("a", emptyType)

	if fetched, err := environment.Get("a"); err != nil {
		t.Error("Get() should not have returned error.")
	} else if fetched != emptyType {
		t.Error("Set() failed.")
	}
}

func Test_Get_Should_Return_Bindings_From_Outer_Environments(t *testing.T) {
	environment1 := NewEnvironment(nil, []string{}, []ApocalispType{})
	environment2 := NewEnvironment(environment1, []string{}, []ApocalispType{})
	environment3 := NewEnvironment(environment2, []string{}, []ApocalispType{})

	emptyType := ApocalispType{}

	environment1.Set("a", emptyType)

	if fetched, err := environment3.Get("a"); err != nil {
		t.Error("Get() should not have returned error.")
	} else if fetched != emptyType {
		t.Error("Set() failed.")
	}
}

func Test_Get_Should_Return_Error_If_Symbol_Not_Found(t *testing.T) {
	environment := NewEnvironment(nil, []string{}, []ApocalispType{})

	_, err := environment.Get("a")
	if err == nil {
		t.Error("Get() should have returned error.")
	}
}

func Test_Find_Should_Return_Environment_Containing_Binding(t *testing.T) {
	environment1 := NewEnvironment(nil, []string{}, []ApocalispType{})
	environment2 := NewEnvironment(environment1, []string{}, []ApocalispType{})
	environment3 := NewEnvironment(environment2, []string{}, []ApocalispType{})

	emptyType := ApocalispType{}

	environment1.Set("a", emptyType)

	if environment3.Find("a") != environment1 {
		t.Error("Find() failed.")
	}
}

func Test_Find_Should_Return_Nil_If_No_Environment_Contains_Binding(t *testing.T) {
	environment1 := NewEnvironment(nil, []string{}, []ApocalispType{})
	environment2 := NewEnvironment(environment1, []string{}, []ApocalispType{})
	environment3 := NewEnvironment(environment2, []string{}, []ApocalispType{})

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

func Test_SetNativeFunction_Should_Add_Native_Function_Binding_To_Environment(t *testing.T) {
	environment := NewEnvironment(nil, []string{}, []ApocalispType{})

	environment.SetNativeFunction("print", func(...ApocalispType) ApocalispType {
		return ApocalispType{}
	})

	if fetched, err := environment.Get("print"); err != nil {
		t.Error("Get() should not have returned error.")
	} else if fetched.NativeFunction == nil || fetched.Symbol == nil {
		t.Error("SetNativeFunction() failed.")
	}
}
