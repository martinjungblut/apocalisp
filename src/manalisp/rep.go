package manalisp

import (
	"errors"
	"fmt"
)

func Rep(sexpr string, environment *Environment, eval func(*ManalispType, *Environment) (*ManalispType, error)) (string, error) {
	t, err := read(sexpr)
	if err != nil {
		return "", err
	} else if t == nil {
		return "", nil
	}

	evaluated, err := eval(t, environment)
	if err != nil {
		return "", err
	}

	return print(evaluated), nil
}

func NoEval(node *ManalispType, environment *Environment) (*ManalispType, error) {
	return node, nil
}

func Step2Eval(node *ManalispType, environment *Environment) (*ManalispType, error) {
	if !node.IsList() {
		if t, err := evalAst(node, environment, Step2Eval); err == nil {
			return t, nil
		} else {
			return nil, err
		}
	} else if node.IsEmptyList() {
		return node, nil
	} else if node.IsList() {
		if container, err := evalAst(node, environment, Step2Eval); err == nil {
			objects := container.AsList()
			function := objects[1]
			parameters := objects[2:]

			if function.IsNativeFunction() {
				result := function.CallNativeFunction(parameters...)
				return &result, nil
			} else {
				return nil, errors.New(fmt.Sprintf("Symbol is not a function: %s", function.ToString()))
			}
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Unexpected behavior.")
}

func Step3Eval(node *ManalispType, environment *Environment) (*ManalispType, error) {
	if !node.IsList() {
		if t, err := evalAst(node, environment, Step3Eval); err == nil {
			return t, nil
		} else {
			return nil, err
		}
	} else if node.IsEmptyList() {
		return node, nil
	} else if node.IsList() {
		if container, err := evalAst(node, environment, Step3Eval); err == nil {
			objects := container.AsList()
			function := objects[1]
			parameters := objects[2:]

			if function.IsNativeFunction() {
				result := function.CallNativeFunction(parameters...)
				return &result, nil
			} else if function.IsSymbol() && function.AsSymbol() == "def!" {
				if len(parameters) == 2 && parameters[0].IsSymbol() {
					environment.Set(parameters[0].AsSymbol(), parameters[1])
					return &parameters[1], nil
				} else {
					return nil, errors.New("Invalid syntax.")
				}
			} else {
				return nil, errors.New(fmt.Sprintf("Symbol is not a function: %s", function.ToString()))
			}
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Unexpected behavior.")
}

func read(sexpr string) (*ManalispType, error) {
	return Parse(sexpr)
}

func print(node *ManalispType) string {
	return node.ToString()
}

func evalAst(node *ManalispType, environment *Environment, eval func(*ManalispType, *Environment) (*ManalispType, error)) (*ManalispType, error) {
	if node.IsSymbol() {
		f := environment.Get(node.AsSymbol())
		return &f, nil
	}

	if node.IsList() {
		all := NewList()
		for _, element := range node.AsList() {
			if evaluated, err := eval(&element, environment); err == nil {
				all.AddToList(*evaluated)
			} else {
				return nil, err
			}
		}
		return all, nil
	}

	if node.IsVector() {
		all := NewVector()
		for _, element := range node.AsVector() {
			if evaluated, err := eval(&element, environment); err == nil {
				all.AddToVector(*evaluated)
			} else {
				return nil, err
			}
		}
		return all, nil
	}

	if node.IsHashmap() {
		all := NewHashmap()
		for _, element := range node.AsHashmap() {
			if evaluated, err := eval(&element, environment); err == nil {
				all.AddToHashmap(*evaluated)
			} else {
				return nil, err
			}
		}
		return all, nil
	}

	return node, nil
}
