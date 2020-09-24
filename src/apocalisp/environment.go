package apocalisp

import (
	"apocalisp/core"
	"apocalisp/parser"
)

// Expose DefaultEnvironment() through the 'apocalisp' namespace.
func DefaultEnvironment(eval func(*core.Type, *core.Environment, bool) (*core.Type, error)) *core.Environment {
	environment := core.DefaultEnvironment(parser.Parser{})

	environment.SetCallable("map", func(args ...core.Type) core.Type {
		result := core.NewList()

		if len(args) >= 2 && args[1].IsIterable() {
			f, iterable := args[0], args[1].AsIterable()
			if f.IsFunction() {
				for _, e := range iterable {
					result.Append(f.CallFunction(e))
				}
			}
		}

		return *result
	})

	environment.SetCallable("eval", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			node := args[0]
			if r, err := eval(&node, environment, true); err == nil {
				return *r
			}
		}
		return *core.NewNil()
	})

	return environment
}
