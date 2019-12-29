package apocalisp

import (
	"apocalisp/core"
	"apocalisp/parser"
	"errors"
	"fmt"
	"strings"
)

func Rep(sexpr string, environment *core.Environment, eval func(*core.Type, *core.Environment) (*core.Type, error)) (string, error) {
	// read
	t, err := parser.Parse(sexpr)
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
	return evaluated.ToString(true), nil
}

func NoEval(node *core.Type, environment *core.Environment) (*core.Type, error) {
	return node, nil
}

func Step2Eval(node *core.Type, environment *core.Environment) (*core.Type, error) {
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
			return evalCallable(container)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func Step3Eval(node *core.Type, environment *core.Environment) (*core.Type, error) {
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
			return evalCallable(container)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func Step4Eval(node *core.Type, environment *core.Environment) (*core.Type, error) {
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
			return evalCallable(container)
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func Step5Eval(node *core.Type, environment *core.Environment) (*core.Type, error) {
	for {
		if !node.IsList() {
			if t, err := evalAst(node, environment, Step5Eval); err != nil {
				return nil, err
			} else {
				return t, nil
			}
		} else if node.IsEmptyList() {
			return node, nil
		} else if node.IsList() {
			first, rest := node.AsList()[0], node.AsList()[1:]

			if first.IsSymbol() && first.AsSymbol() == "def!" {
				return specialFormDef(Step5Eval, environment)(rest)
			} else if first.IsSymbol() && first.AsSymbol() == "let*" {
				if err := tcoSpecialFormLet(Step5Eval, &node, &environment)(rest); err != nil {
					return nil, err
				}
			} else if first.IsSymbol() && first.AsSymbol() == "do" {
				if err := tcoSpecialFormDo(Step5Eval, &node, &environment)(rest); err != nil {
					return nil, err
				}
			} else if first.IsSymbol() && (first.AsSymbol() == "fn*" || first.AsSymbol() == "λ") {
				return tcoSpecialFormFn(Step5Eval, &node, &environment)(rest)
			} else if first.IsSymbol() && first.AsSymbol() == "if" {
				if err := tcoSpecialFormIf(Step5Eval, &node, &environment)(rest); err != nil {
					return nil, err
				}
			} else if container, err := evalAst(node, environment, Step5Eval); err == nil {
				f, args := container.AsList()[0], container.AsList()[1:]
				if f.IsFunction() {
					node = &f.Function.Body
					environment = core.NewEnvironment(&f.Function.Environment, f.Function.Params, args)
				} else {
					return evalCallable(container)
				}
			} else {
				return nil, err
			}
		} else {
			break
		}
	}

	return nil, errors.New("Error: Unexpected behavior.")
}

func specialFormDef(eval func(*core.Type, *core.Environment) (*core.Type, error), environment *core.Environment) func([]core.Type) (*core.Type, error) {
	return func(rest []core.Type) (*core.Type, error) {
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

func tcoSpecialFormLet(eval func(*core.Type, *core.Environment) (*core.Type, error), node **core.Type, environment **core.Environment) func([]core.Type) error {
	return func(rest []core.Type) error {
		if len(rest) != 2 || !rest[0].EvenIterable() {
			return errors.New("Error: Invalid syntax for `let*`.")
		} else {
			letEnvironment := core.NewEnvironment(*environment, []string{}, []core.Type{})

			bindings := rest[0].Iterable()
			for i, j := 0, 1; i < len(bindings); i, j = i+2, j+2 {
				s := bindings[i].ToString(true)
				if e, ierr := eval(&bindings[j], letEnvironment); ierr == nil {
					letEnvironment.Set(s, *e)
				} else {
					return ierr
				}
			}

			*environment = letEnvironment
			*node = &rest[1]
			return nil
		}
	}
}

func specialFormLet(eval func(*core.Type, *core.Environment) (*core.Type, error), environment *core.Environment) func([]core.Type) (*core.Type, error) {
	return func(rest []core.Type) (*core.Type, error) {
		if len(rest) != 2 || !rest[0].EvenIterable() {
			return nil, errors.New("Error: Invalid syntax for `let*`.")
		} else {
			letEnvironment := core.NewEnvironment(environment, []string{}, []core.Type{})

			bindings := rest[0].Iterable()
			for i, j := 0, 1; i < len(bindings); i, j = i+2, j+2 {
				s := bindings[i].ToString(true)
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

func tcoSpecialFormDo(eval func(*core.Type, *core.Environment) (*core.Type, error), node **core.Type, environment **core.Environment) func([]core.Type) error {
	return func(rest []core.Type) error {
		if len(rest) < 1 {
			return errors.New("Error: Invalid syntax for `do`.")
		} else {
			toEvaluate := rest[:len(rest)-1]
			if _, err := evalAst(&core.Type{List: &toEvaluate}, *environment, eval); err != nil {
				return err
			} else {
				*node = &rest[len(rest)-1]
				return nil
			}
		}
	}
}

func specialFormDo(eval func(*core.Type, *core.Environment) (*core.Type, error), environment *core.Environment) func([]core.Type) (*core.Type, error) {
	return func(rest []core.Type) (*core.Type, error) {
		if len(rest) < 1 {
			return nil, errors.New("Error: Invalid syntax for `do`.")
		} else {
			toEvaluate := &core.Type{List: &rest}
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

func tcoSpecialFormFn(eval func(*core.Type, *core.Environment) (*core.Type, error), node **core.Type, environment **core.Environment) func([]core.Type) (*core.Type, error) {
	return func(rest []core.Type) (*core.Type, error) {
		if len(rest) < 2 || (rest[0].IsList() && rest[0].IsVector()) {
			return nil, errors.New("Error: Invalid syntax for `fn*`.")
		} else {
			var symbols []string
			for _, node := range rest[0].Iterable() {
				if node.IsSymbol() {
					symbols = append(symbols, node.AsSymbol())
				} else {
					return nil, errors.New("Error: Invalid syntax for `fn*`.")
				}
			}

			callable := func(args ...core.Type) core.Type {
				newEnvironment := core.NewEnvironment(*environment, symbols, args)
				if result, err := eval(&rest[1], newEnvironment); err != nil {
					errorMessage := err.Error()
					return core.Type{String: &errorMessage}
				} else {
					return *result
				}
			}

			function := core.Function{
				Params:      symbols,
				Body:        rest[1],
				Callable:    callable,
				Environment: **environment,
			}

			return &core.Type{Function: &function}, nil
		}
	}
}

func specialFormFn(eval func(*core.Type, *core.Environment) (*core.Type, error), environment *core.Environment) func([]core.Type) (*core.Type, error) {
	return func(rest []core.Type) (*core.Type, error) {
		if len(rest) < 2 || (rest[0].IsList() && rest[0].IsVector()) {
			return nil, errors.New("Error: Invalid syntax for `fn*`.")
		} else {
			var symbols []string
			for _, node := range rest[0].Iterable() {
				if node.IsSymbol() {
					symbols = append(symbols, node.AsSymbol())
				} else {
					return nil, errors.New("Error: Invalid syntax for `fn*`.")
				}
			}

			callable := func(args ...core.Type) core.Type {
				newEnvironment := core.NewEnvironment(environment, symbols, args)
				if result, err := eval(&rest[1], newEnvironment); err != nil {
					errorMessage := err.Error()
					return core.Type{String: &errorMessage}
				} else {
					return *result
				}
			}

			return &core.Type{Callable: &callable}, nil
		}
	}
}

func tcoSpecialFormIf(eval func(*core.Type, *core.Environment) (*core.Type, error), node **core.Type, environment **core.Environment) func([]core.Type) error {
	return func(rest []core.Type) error {
		length := len(rest)

		if length < 2 || length > 3 {
			return errors.New("Error: Invalid syntax for `if`.")
		} else if condition, err := eval(&rest[0], *environment); err != nil {
			return err
		} else if !condition.IsNil() && !condition.IsBoolean(false) {
			*node = &rest[1]
		} else if length == 3 {
			*node = &rest[2]
		} else {
			*node = core.NewNil()
		}
		return nil
	}
}

func specialFormIf(eval func(*core.Type, *core.Environment) (*core.Type, error), environment *core.Environment) func([]core.Type) (*core.Type, error) {
	return func(rest []core.Type) (*core.Type, error) {
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
			return core.NewNil(), nil
		}
	}
}

func evalCallable(node *core.Type) (*core.Type, error) {
	first, rest := node.AsList()[0], node.AsList()[1:]

	if first.IsCallable() {
		result := first.CallCallable(rest...)
		return &result, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Error: '%s' is not a function.", first.ToString(true)))
	}
}

func evalAst(node *core.Type, environment *core.Environment, eval func(*core.Type, *core.Environment) (*core.Type, error)) (*core.Type, error) {
	if node.IsSymbol() && !strings.HasPrefix(node.AsSymbol(), ":") {
		if f, err := environment.Get(node.AsSymbol()); err != nil {
			return nil, err
		} else {
			return &f, nil
		}
	}

	if node.IsList() {
		all := core.NewList()
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
		all := core.NewVector()
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
		newHashmap := core.NewHashmap()
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
