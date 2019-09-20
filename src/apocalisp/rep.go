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
		if t, err := evalAst(node, environment, Step2Eval); err != nil {
			return nil, err
		} else {
			return t, nil
		}
	} else if node.IsEmptyList() {
		return node, nil
	} else if node.IsList() {
		if container, err := evalAst(node, environment, Step2Eval); err == nil {
			first, rest := container.AsList()[1], container.AsList()[2:]

			if first.IsNativeFunction() {
				result := first.CallNativeFunction(rest...)
				return &result, nil
			} else {
				return nil, errors.New(fmt.Sprintf("Symbol is not a function: `%s`.", first.ToString()))
			}
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Unexpected behavior.")
}

func Step3Eval(node *ApocalispType, environment *Environment) (*ApocalispType, error) {
	if !node.IsList() {
		if t, err := evalAst(node, environment, Step3Eval); err != nil {
			return nil, err
		} else {
			return t, nil
		}
	} else if node.IsEmptyList() {
		return node, nil
	} else if node.IsList() {
		first, rest := node.AsList()[0], node.AsList()[1:]

		if first.IsSymbol() && first.AsSymbol() == "def!" {
			if len(rest) != 2 || !rest[0].IsSymbol() {
				return nil, errors.New("Invalid syntax for `def!`.")
			} else {
				if e, ierr := Step3Eval(&rest[1], environment); ierr == nil {
					environment.Set(rest[0].AsSymbol(), *e)
					return e, nil
				} else {
					return nil, ierr
				}
			}
		} else if first.IsSymbol() && first.AsSymbol() == "let*" {
			if len(rest) != 2 || !rest[0].EvenIterable() {
				return nil, errors.New("Invalid syntax for `let*`.")
			} else {
				letEnvironment := NewEnvironment(environment)

				bindings := rest[0].Iterable()
				for i, j := 0, 1; i < len(bindings); i, j = i+2, j+2 {
					s := bindings[i].ToString()
					if e, ierr := Step3Eval(&bindings[j], letEnvironment); ierr == nil {
						letEnvironment.Set(s, *e)
					} else {
						return nil, ierr
					}
				}

				return Step3Eval(&rest[1], letEnvironment)
			}
		} else if evaluated, err := evalAst(node, environment, Step3Eval); err == nil {
			first, rest := evaluated.AsList()[1], evaluated.AsList()[2:]
			if first.IsNativeFunction() {
				result := first.CallNativeFunction(rest...)
				return &result, nil
			} else {
				return nil, errors.New(fmt.Sprintf("Symbol is not a function: `%s`.", first.ToString()))
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
