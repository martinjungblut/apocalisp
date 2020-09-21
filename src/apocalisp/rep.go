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

func Evaluate(node *core.Type, environment *core.Environment) (*core.Type, error) {
	var lexicalReturnValue *core.Type
	var lexicalError error
	wrapReturn := func(node *core.Type, err error) {
		if node != nil {
			lexicalReturnValue = node
		}
		if err != nil {
			lexicalError = err
		}
	}
	processReturn := func() (*core.Type, error) {
		if lexicalReturnValue != nil {
			if message := lexicalReturnValue.ExceptionMessage(); message != nil {
				return nil, errors.New(fmt.Sprintf("Exception: %s", *message))
			} else {
				return lexicalReturnValue, nil
			}
		}
		if lexicalError != nil {
			return nil, lexicalError
		}
		return nil, errors.New("Error: Unexpected behavior.")
	}

	// TCO loop
	for {
		if lexicalReturnValue != nil || lexicalError != nil {
			return processReturn()
		}

		expanded := macroexpand(*node, *environment)
		node = &expanded

		if !node.IsList() {
			if evaluated, err := evalAst(node, environment, Evaluate); err != nil {
				wrapReturn(nil, err)
			} else {
				wrapReturn(evaluated, nil)
			}
		} else if node.IsEmptyList() {
			wrapReturn(node, nil)
		} else if node.IsList() {
			first, rest := node.AsIterable()[0], node.AsIterable()[1:]

			if first.CompareSymbol("def!") {
				wrapReturn(specialFormDef(Evaluate, rest, environment))
			} else if first.CompareSymbol("defmacro!") {
				wrapReturn(specialFormDefmacro(Evaluate, rest, environment))
			} else if first.CompareSymbol("macroexpand") {
				expanded := macroexpand(rest[0], *environment)
				wrapReturn(&expanded, nil)
			} else if first.CompareSymbol("let*") {
				if err := tcoSpecialFormLet(Evaluate, rest, &node, &environment); err != nil {
					wrapReturn(nil, err)
				}
			} else if first.CompareSymbol("do") {
				if err := tcoSpecialFormDo(Evaluate, rest, &node, &environment); err != nil {
					wrapReturn(nil, err)
				}
			} else if first.CompareSymbol("fn*", `\`) {
				wrapReturn(tcoSpecialFormFn(Evaluate, rest, &node, &environment))
			} else if first.CompareSymbol("if") {
				if err := tcoSpecialFormIf(Evaluate, rest, &node, &environment); err != nil {
					wrapReturn(nil, err)
				}
			} else if first.CompareSymbol("quasiquote") {
				if err := tcoSpecialFormQuasiquote(Evaluate, rest, &node, &environment); err != nil {
					wrapReturn(nil, err)
				}
			} else if first.CompareSymbol("quasiquoteexpand") {
				wrapReturn(specialFormQuasiquoteexpand(Evaluate, rest, environment))
			} else if first.CompareSymbol("quote") {
				wrapReturn(specialFormQuote(Evaluate, rest, environment))
			} else if container, err := evalAst(node, environment, Evaluate); err != nil {
				wrapReturn(nil, err)
			} else {
				first, rest := container.AsIterable()[0], container.AsIterable()[1:]
				if first.IsFunction() {
					node = &first.Function.Body
					environment = core.NewEnvironment(&first.Function.Environment, first.Function.Params, rest)
				} else {
					wrapReturn(evalCallable(container))
				}
			}
		} else {
			return processReturn()
		}
	}
}

func tcoSpecialFormLet(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) error {
	if len(rest) != 2 || !rest[0].IsEvenIterable() {
		return errors.New("Error: Invalid syntax for `let*`.")
	} else {
		letEnvironment := core.NewEnvironment(*environment, []string{}, []core.Type{})

		bindings := rest[0].AsIterable()
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

func tcoSpecialFormDo(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) error {
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

func tcoSpecialFormIf(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) error {
	length := len(rest)

	if length < 2 || length > 3 {
		return errors.New("Error: Invalid syntax for `if`.")
	} else if condition, err := eval(&rest[0], *environment); err != nil {
		return err
	} else if !condition.IsNil() && !condition.CompareBoolean(false) {
		*node = &rest[1]
	} else if length == 3 {
		*node = &rest[2]
	} else {
		*node = core.NewNil()
	}

	return nil
}

func tcoSpecialFormFn(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) (*core.Type, error) {
	if len(rest) < 2 || !rest[0].IsIterable() {
		return nil, errors.New("Error: Invalid syntax for `fn*`.")
	} else {
		var symbols []string
		for _, node := range rest[0].AsIterable() {
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

func tcoSpecialFormQuasiquote(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) error {
	if len(rest) < 1 {
		return errors.New("Error: Invalid syntax for `quasiquote`.")
	} else {
		newNode := quasiquote(rest[0])
		*node = &newNode
	}
	return nil
}

func specialFormQuote(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, environment *core.Environment) (*core.Type, error) {
	if len(rest) < 1 {
		return nil, errors.New("Error: Invalid syntax for `quote`.")
	} else {
		return &rest[0], nil
	}
}

func specialFormQuasiquoteexpand(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, environment *core.Environment) (*core.Type, error) {
	if len(rest) >= 1 {
		newNode := quasiquote(rest[0])
		return &newNode, nil
	}
	return nil, errors.New("Error: Invalid syntax for `quasiquoteexpand`.")
}

func specialFormDef(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, environment *core.Environment) (*core.Type, error) {
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

func specialFormDefmacro(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, environment *core.Environment) (*core.Type, error) {
	if len(rest) != 2 || !rest[0].IsSymbol() {
		return nil, errors.New("Error: Invalid syntax for `def!`.")
	} else {
		if e, ierr := eval(&rest[1], environment); ierr == nil {
			if e.IsFunction() {
				e.Function.IsMacro = true
			}
			environment.Set(rest[0].AsSymbol(), *e)
			return e, nil
		} else {
			return nil, ierr
		}
	}
}

func evalCallable(node *core.Type) (*core.Type, error) {
	first, rest := node.AsIterable()[0], node.AsIterable()[1:]

	if first.IsCallable() {
		result := first.CallCallable(rest...)
		return &result, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Error: '%s' is not a function.", first.ToString(true)))
	}
}

func evalAst(node *core.Type, environment *core.Environment, eval func(*core.Type, *core.Environment) (*core.Type, error)) (*core.Type, error) {
	if node.IsSymbol() && !strings.HasPrefix(node.AsSymbol(), ":") {
		if t, err := environment.Get(node.AsSymbol()); err != nil {
			return nil, err
		} else {
			return &t, nil
		}
	}

	if node.IsIterable() {
		newIterable := node.DeriveIterable()
		for _, element := range node.AsIterable() {
			if evaluated, err := eval(&element, environment); err != nil {
				return nil, err
			} else {
				newIterable.Append(*evaluated)
			}
		}
		return newIterable, nil
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

func quasiquote(node core.Type) core.Type {
	iterable := node.AsIterable()
	unquoted := len(iterable) >= 2 && iterable[0].CompareSymbol("unquote")

	if node.IsSymbol() || node.IsHashmap() || (node.IsVector() && unquoted) {
		return *core.NewList(*core.NewSymbol("quote"), node)
	} else if node.IsVector() {
		return *core.NewList(*core.NewSymbol("vec"), quasiquote(*core.NewList(iterable...)))
	} else if unquoted {
		return iterable[1]
	} else if len(iterable) >= 1 {
		result := *core.NewList()

		for i := len(iterable) - 1; i >= 0; i-- {
			el := iterable[i]
			eli := el.AsIterable()
			if len(eli) >= 2 && eli[0].CompareSymbol("splice-unquote") {
				result = *core.NewList(*core.NewSymbol("concat"), eli[1], result)
			} else {
				result = *core.NewList(*core.NewSymbol("cons"), quasiquote(el), result)
			}
		}

		return result
	}

	return node
}

func isMacroCall(node core.Type, environment core.Environment, capture func(core.Type)) bool {
	if iterable := node.AsIterable(); node.IsList() && len(iterable) >= 1 {
		if first := iterable[0]; first.IsSymbol() {
			if macro, err := environment.Get(first.AsSymbol()); err == nil && macro.IsMacroFunction() {
				capture(macro)
				return true
			}
		}
	}
	return false
}

func macroexpand(node core.Type, environment core.Environment) core.Type {
	var macro core.Type
	capture := func(m core.Type) {
		macro = m
	}

	for isMacroCall(node, environment, capture) {
		parameters := node.AsIterable()[1:]
		node = macro.CallFunction(parameters...)
	}

	return node
}
