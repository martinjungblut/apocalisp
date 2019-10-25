package apocalisp

import (
	"apocalisp/typing"
	"errors"
	"fmt"
)

func Rep(sexpr string, environment *Environment, eval func(*typing.Type, *Environment) (*typing.Type, error)) (string, error) {
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

func NoEval(node *typing.Type, environment *Environment) (*typing.Type, error) {
	return node, nil
}

func Step2Eval(node *typing.Type, environment *Environment) (*typing.Type, error) {
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
			return evalNativeFunction(container)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func Step3Eval(node *typing.Type, environment *Environment) (*typing.Type, error) {
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
			return specialFormDef(Step3Eval, environment)(rest)
		} else if first.IsSymbol() && first.AsSymbol() == "let*" {
			return specialFormLet(Step3Eval, environment)(rest)
		} else if container, err := evalAst(node, environment, Step3Eval); err == nil {
			return evalNativeFunction(container)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func Step4Eval(node *typing.Type, environment *Environment) (*typing.Type, error) {
	if !node.IsList() {
		if t, err := evalAst(node, environment, Step4Eval); err != nil {
			return nil, err
		} else {
			return t, nil
		}
	} else if node.IsEmptyList() {
		return node, nil
	} else if node.IsList() {
		first, rest := node.AsList()[0], node.AsList()[1:]

		if first.IsSymbol() && first.AsSymbol() == "def!" {
			return specialFormDef(Step4Eval, environment)(rest)
		} else if first.IsSymbol() && first.AsSymbol() == "let*" {
			return specialFormLet(Step4Eval, environment)(rest)
		} else if first.IsSymbol() && first.AsSymbol() == "do" {
			return specialFormDo(Step4Eval, environment)(rest)
		} else if first.IsSymbol() && (first.AsSymbol() == "fn*" || first.AsSymbol() == "λ") {
			return specialFormFn(Step4Eval, environment)(rest)
		} else if first.IsSymbol() && first.AsSymbol() == "if" {
			return specialFormIf(Step4Eval, environment)(rest)
		} else if container, err := evalAst(node, environment, Step4Eval); err == nil {
			return evalNativeFunction(container)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func specialFormDef(eval func(*typing.Type, *Environment) (*typing.Type, error), environment *Environment) func([]typing.Type) (*typing.Type, error) {
	return func(rest []typing.Type) (*typing.Type, error) {
		if len(rest) != 2 || !rest[0].IsSymbol() {
			return nil, errors.New("Error: Invalid syntax for `def!`.")
		} else {
			if e, ierr := eval(&rest[1], environment); ierr == nil {
				environment.Set(rest[0].AsSymbol(), *e)
				return e, nil
			} else {
				return nil, ierr
			}
		}
	}
}

func specialFormLet(eval func(*typing.Type, *Environment) (*typing.Type, error), environment *Environment) func([]typing.Type) (*typing.Type, error) {
	return func(rest []typing.Type) (*typing.Type, error) {
		if len(rest) != 2 || !rest[0].EvenIterable() {
			return nil, errors.New("Error: Invalid syntax for `let*`.")
		} else {
			letEnvironment := NewEnvironment(environment, []string{}, []typing.Type{})

			bindings := rest[0].Iterable()
			for i, j := 0, 1; i < len(bindings); i, j = i+2, j+2 {
				s := bindings[i].ToString()
				if e, ierr := eval(&bindings[j], letEnvironment); ierr == nil {
					letEnvironment.Set(s, *e)
				} else {
					return nil, ierr
				}
			}

			return eval(&rest[1], letEnvironment)
		}
	}
}

func specialFormDo(eval func(*typing.Type, *Environment) (*typing.Type, error), environment *Environment) func([]typing.Type) (*typing.Type, error) {
	return func(rest []typing.Type) (*typing.Type, error) {
		if len(rest) < 1 {
			return nil, errors.New("Error: Invalid syntax for `do`.")
		} else {
			toEvaluate := &typing.Type{List: &rest}
			if evaluated, err := evalAst(toEvaluate, environment, eval); err != nil {
				return nil, err
			} else {
				list := evaluated.AsList()
				last := list[len(list)-1]
				return &last, nil
			}
		}
	}
}

func specialFormFn(eval func(*typing.Type, *Environment) (*typing.Type, error), environment *Environment) func([]typing.Type) (*typing.Type, error) {
	return func(rest []typing.Type) (*typing.Type, error) {
		if len(rest) < 2 || !rest[0].IsList() || !rest[1].IsList() {
			return nil, errors.New("Error: Invalid syntax for `fn*`.")
		} else {
			var symbols []string
			for _, node := range rest[0].AsList() {
				if node.IsSymbol() {
					symbols = append(symbols, node.AsSymbol())
				} else {
					return nil, errors.New("Error: Invalid syntax for `fn*`.")
				}
			}

			closure := func(args ...typing.Type) typing.Type {
				newEnvironment := NewEnvironment(environment, symbols, args)
				if result, err := eval(&rest[1], newEnvironment); err != nil {
					errorMessage := err.Error()
					return typing.Type{String: &errorMessage}
				} else {
					return *result
				}
			}

			return &typing.Type{NativeFunction: &closure}, nil
		}
	}
}

func specialFormIf(eval func(*typing.Type, *Environment) (*typing.Type, error), environment *Environment) func([]typing.Type) (*typing.Type, error) {
	return func(rest []typing.Type) (*typing.Type, error) {
		length := len(rest)

		if length < 2 || length > 3 {
			return nil, errors.New("Error: Invalid syntax for `if`.")
		} else if condition, err := eval(&rest[0], environment); err != nil {
			return nil, err
		} else if !condition.IsNil() && !condition.IsBoolean(false) {
			if evaluated, err := eval(&rest[1], environment); err != nil {
				return nil, err
			} else {
				return evaluated, nil
			}
		} else if length == 3 {
			if evaluated, err := eval(&rest[2], environment); err != nil {
				return nil, err
			} else {
				return evaluated, nil
			}
		} else {
			return typing.NewNil(), nil
		}
	}
}

func evalNativeFunction(node *typing.Type) (*typing.Type, error) {
	first, rest := node.AsList()[1], node.AsList()[2:]

	if first.IsNativeFunction() {
		result := first.CallNativeFunction(rest...)
		return &result, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Error: '%s' is not a function.", first.ToString()))
	}
}

func evalAst(node *typing.Type, environment *Environment, eval func(*typing.Type, *Environment) (*typing.Type, error)) (*typing.Type, error) {
	if node.IsSymbol() {
		if f, err := environment.Get(node.AsSymbol()); err != nil {
			return nil, err
		} else {
			return &f, nil
		}
	}

	if node.IsList() {
		all := typing.NewList()
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
		all := typing.NewVector()
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
		newHashmap := typing.NewHashmap()
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
