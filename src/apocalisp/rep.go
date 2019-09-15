package apocalisp

import (
	"errors"
	"fmt"
)

func Rep(sexpr string, environment *Environment, eval func(*ApocalispType, *Environment) (*ApocalispType, error)) (string, error) {
	// read
	t, err := Parse(sexpr)
	if err != nil {
		return "", err
	} else if t == nil {
		return "", nil
	}

	// eval
	evaluated, err := eval(t, environment)
	if err != nil {
		return "", err
	}

	// print
	return evaluated.ToString(), nil
}

func NoEval(node *ApocalispType, environment *Environment) (*ApocalispType, error) {
	return node, nil
}

func Step2Eval(node *ApocalispType, environment *Environment) (*ApocalispType, error) {
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

func Step3Eval(node *ApocalispType, environment *Environment) (*ApocalispType, error) {
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

func evalAst(node *ApocalispType, environment *Environment, eval func(*ApocalispType, *Environment) (*ApocalispType, error)) (*ApocalispType, error) {
	if node.IsSymbol() {
		if f, err := environment.Get(node.AsSymbol()); err != nil {
			return nil, err
		} else {
			return &f, nil
		}
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
		currentHashmap := node.AsHashmap()
		newHashmap := NewHashmap()
		for i, j := 0, 1; i < len(currentHashmap); i, j = i+2, j+2 {
			newHashmap.AddToHashmap(currentHashmap[i])
			if evaluated, err := eval(&(currentHashmap[j]), environment); err != nil {
				return nil, err
			} else {
				newHashmap.AddToHashmap(*evaluated)
			}
		}
		return newHashmap, nil
	}

	return node, nil
}
