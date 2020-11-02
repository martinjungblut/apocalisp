package apocalisp

import (
	"apocalisp/core"
	"errors"
	"fmt"
)

func Rep(sexpr string, environment *core.Environment, eval func(*core.Type, *core.Environment) (*core.Type, error), parser core.Parser) (string, error) {
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
	} else if evaluated.IsException() {
		return "", errors.New(evaluated.ToString(false))
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
	processReturn := func() (*core.Type, error) {
		if lexicalError != nil {
			return nil, lexicalError
		} else if lexicalReturnValue != nil {
			return lexicalReturnValue, nil
		}
		return nil, errors.New("Error: Unexpected behavior.")
	}
	wrapReturn := func(node *core.Type, err error) {
		if err != nil {
			lexicalError = err
		} else if node != nil {
			lexicalReturnValue = node
		}
	}

	// TCO loop
	for {
		if lexicalReturnValue != nil || lexicalError != nil {
			return processReturn()
		}

		expanded := macroexpand(*node, *environment)
		node = &expanded

		if !node.IsList() {
			wrapReturn(evalAst(node, environment, Evaluate))
		} else if node.IsList() && node.IsEmptyIterable() {
			wrapReturn(node, nil)
		} else if node.IsList() && !node.IsEmptyIterable() {
			first, rest := node.AsIterable()[0], node.AsIterable()[1:]

			if first.CompareSymbol("def!") {
				wrapReturn(specialFormDef(Evaluate, rest, environment))
			} else if first.CompareSymbol("defmacro!") {
				wrapReturn(specialFormDefmacro(Evaluate, rest, environment))
			} else if first.CompareSymbol("macroexpand") {
				expanded := macroexpand(rest[0], *environment)
				wrapReturn(&expanded, nil)
			} else if first.CompareSymbol("let*") {
				wrapReturn(tcoSpecialFormLet(Evaluate, rest, &node, &environment))
			} else if first.CompareSymbol("do") {
				wrapReturn(tcoSpecialFormDo(Evaluate, rest, &node, &environment))
			} else if first.CompareSymbol("fn*", `\`) {
				wrapReturn(tcoSpecialFormFn(Evaluate, rest, &node, &environment))
			} else if first.CompareSymbol("if") {
				wrapReturn(tcoSpecialFormIf(Evaluate, rest, &node, &environment))
			} else if first.CompareSymbol("quasiquote") {
				wrapReturn(tcoSpecialFormQuasiquote(Evaluate, rest, &node, &environment))
			} else if first.CompareSymbol("quasiquoteexpand") {
				wrapReturn(specialFormQuasiquoteexpand(Evaluate, rest, environment))
			} else if first.CompareSymbol("quote") {
				wrapReturn(specialFormQuote(Evaluate, rest, environment))
			} else if first.CompareSymbol("try*") {
				wrapReturn(specialFormTryCatch(Evaluate, rest, environment))
			} else {
				if container, err := evalAst(node, environment, Evaluate); err != nil {
					wrapReturn(nil, err)
				} else {
					function, parameters := container.AsIterable()[0], container.AsIterable()[1:]
					if function.IsFunction() {
						node = &function.Function.Body
						environment = core.NewEnvironment(&function.Function.Environment, function.Function.Params, parameters)
					} else {
						wrapReturn(evalCallable(container))
					}
				}
			}
		} else {
			return processReturn()
		}
	}
}

func tcoSpecialFormLet(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) (*core.Type, error) {
	if len(rest) != 2 || !rest[0].IsEvenIterable() {
		return nil, errors.New("Error: Invalid syntax for `let*`.")
	} else {
		letEnvironment := core.NewEnvironment(*environment, []string{}, []core.Type{})
		bindings, body := rest[0].AsIterable(), &rest[1]

		for symbol, target := 0, 1; symbol < len(bindings); symbol, target = symbol+2, target+2 {
			if e, ierr := eval(&bindings[target], letEnvironment); ierr != nil {
				return nil, ierr
			} else if e.IsException() {
				*node = e
				return nil, nil
			} else {
				letEnvironment.Set(bindings[symbol].ToString(true), *e)
			}
		}

		*environment, *node = letEnvironment, body
		return nil, nil
	}
}

func tcoSpecialFormDo(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) (*core.Type, error) {
	if len(rest) < 1 {
		return nil, errors.New("Error: Invalid syntax for `do`.")
	} else {
		toEvaluate := rest[:len(rest)-1]
		if _, err := evalAst(&core.Type{List: &toEvaluate}, *environment, eval); err != nil {
			return nil, err
		} else {
			*node = &rest[len(rest)-1]
			return nil, nil
		}
	}
}

func tcoSpecialFormIf(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) (*core.Type, error) {
	length := len(rest)

	if length < 2 || length > 3 {
		return nil, errors.New("Error: Invalid syntax for `if`.")
	} else if condition, err := eval(&rest[0], *environment); err != nil {
		return nil, err
	} else if !condition.IsNil() && !condition.CompareBoolean(false) {
		*node = &rest[1]
	} else if length == 3 {
		*node = &rest[2]
	} else {
		*node = core.NewNil()
	}

	return nil, nil
}

func tcoSpecialFormFn(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) (*core.Type, error) {
	if len(rest) < 2 || !rest[0].IsIterable() {
		return nil, errors.New("Error: Invalid syntax for `fn*`.")
	} else {
		symbols := make([]string, 0)
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
				return *core.NewStringException(err.Error())
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

func tcoSpecialFormQuasiquote(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, node **core.Type, environment **core.Environment) (*core.Type, error) {
	if len(rest) < 1 {
		return nil, errors.New("Error: Invalid syntax for `quasiquote`.")
	} else {
		newNode := quasiquote(rest[0])
		*node = &newNode
	}
	return nil, nil
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
		if e, ierr := eval(&rest[1], environment); ierr != nil {
			return nil, ierr
		} else if e.IsException() {
			return e, nil
		} else {
			environment.Set(rest[0].AsSymbol(), *e)
			return e, nil
		}
	}
}

func specialFormDefmacro(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, environment *core.Environment) (*core.Type, error) {
	if len(rest) != 2 || !rest[0].IsSymbol() {
		return nil, errors.New("Error: Invalid syntax for `defmacro!`.")
	} else {
		if e, ierr := eval(&rest[1], environment); ierr != nil {
			return nil, ierr
		} else if e.IsException() {
			return e, nil
		} else {
			var macro *core.Type
			if e.IsFunction() {
				newFunction := core.Function{
					IsMacro:     true,
					Environment: e.Function.Environment,
					Body:        e.Function.Body,
					Params:      e.Function.Params,
					Callable:    e.Function.Callable,
				}
				macro = &core.Type{Function: &newFunction, Metadata: e.Metadata}
			} else {
				macro = e
			}
			environment.Set(rest[0].AsSymbol(), *macro)
			return macro, nil
		}
	}
}

func specialFormTryCatch(eval func(*core.Type, *core.Environment) (*core.Type, error), rest []core.Type, environment *core.Environment) (*core.Type, error) {
	if len(rest) < 1 {
		return nil, errors.New("Error: Invalid syntax for `try*!`.")
	}

	e, err := eval(&rest[0], environment)

	if err != nil {
		return nil, err
	} else if len(rest) >= 2 {
		catchexp := rest[1].AsIterable()
		if e.IsException() && len(catchexp) == 3 && catchexp[0].CompareSymbol("catch*") && catchexp[1].IsSymbol() {
			symbol, body := catchexp[1].AsSymbol(), catchexp[2]
			return eval(&body, core.NewEnvironment(environment, []string{symbol}, []core.Type{*e.AsException()}))
		}
	}

	return e, nil
}

func evalCallable(node *core.Type) (*core.Type, error) {
	first, rest := node.AsIterable()[0], node.AsIterable()[1:]

	if first.IsCallable() {
		result := first.CallCallable(rest...)
		return &result, nil
	} else if first.IsException() {
		return &first, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Error: '%s' is not a function.", first.ToString(true)))
	}
}

func evalAst(node *core.Type, environment *core.Environment, eval func(*core.Type, *core.Environment) (*core.Type, error)) (*core.Type, error) {
	if node.IsSymbol() && !node.IsKeyword() {
		value := environment.Get(node.AsSymbol())
		return &value, nil
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
		currentHashmap, newHashmap := node.AsHashmap(), core.NewHashmap()
		for key, value := range currentHashmap {
			if evaluated, err := eval(&value, environment); err != nil {
				return nil, err
			} else {
				newHashmap.HashmapSet(key, *evaluated)
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
			if macro := environment.Get(first.AsSymbol()); macro.IsMacroFunction() {
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
