package apocalisp

import (
	"apocalisp/typing"
	"errors"
	"fmt"
	"strings"
)

type Environment struct {
	outer *Environment
	table map[string]typing.Type
}

func NewEnvironment(outer *Environment, symbols []string, nodes []typing.Type) *Environment {
	table := make(map[string]typing.Type)

	environment := &Environment{
		table: table,
		outer: outer,
	}

	for i := 0; i < len(symbols); i++ {
		if symbols[i] == "&" {
			rest := nodes[i:]
			if i+1 < len(symbols) {
				environment.Set(symbols[i+1], typing.Type{List: &rest})
			} else {
				environment.Set("&", typing.Type{List: &rest})
			}
			break
		} else {
			if i < len(nodes) {
				environment.Set(symbols[i], nodes[i])
			}
		}
	}

	return environment
}

func (env *Environment) Set(symbol string, node typing.Type) {
	env.table[symbol] = node
}

func (env *Environment) SetCallable(symbol string, callable func(...typing.Type) typing.Type) {
	env.table[symbol] = typing.Type{
		Callable: &callable,
		Symbol:   &symbol,
	}
}

func (env *Environment) Find(symbol string) *Environment {
	for key, _ := range env.table {
		if key == symbol {
			return env
		}
	}

	if env.outer != nil {
		return env.outer.Find(symbol)
	}

	return nil
}

func (env *Environment) Get(symbol string) (typing.Type, error) {
	if e := env.Find(symbol); e != nil {
		for key, value := range e.table {
			if key == symbol {
				return value, nil
			}
		}
	}
	return typing.Type{}, errors.New(fmt.Sprintf("Error: '%s' not found.", symbol))
}

func DefaultEnvironment() *Environment {
	env := NewEnvironment(nil, []string{}, []typing.Type{})

	env.SetCallable("+", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return typing.Type{Integer: &r}
	})

	env.SetCallable("-", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return typing.Type{Integer: &r}
	})

	env.SetCallable("/", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return typing.Type{Integer: &r}
	})

	env.SetCallable("*", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return typing.Type{Integer: &r}
	})

	env.SetCallable("list", func(args ...typing.Type) typing.Type {
		list := typing.NewList()
		for _, arg := range args {
			list.AddToList(arg)
		}
		return *list
	})

	env.SetCallable("list?", func(args ...typing.Type) typing.Type {
		return *typing.NewBoolean(args[0].IsList())
	})

	env.SetCallable("empty?", func(args ...typing.Type) typing.Type {
		var value int64 = int64(len(args[0].Iterable()))
		return *typing.NewBoolean(value == 0)
	})

	env.SetCallable("count", func(args ...typing.Type) typing.Type {
		var value int64 = int64(len(args[0].Iterable()))
		return typing.Type{Integer: &value}
	})

	env.SetCallable("=", func(args ...typing.Type) typing.Type {
		if len(args) == 2 {
			return *typing.NewBoolean(compareNodes(args[0], args[1]))
		}
		return *typing.NewBoolean(false)
	})

	env.SetCallable("<", func(args ...typing.Type) typing.Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() < args[1].AsInteger()
			}
		}
		return *typing.NewBoolean(result)
	})

	env.SetCallable("<=", func(args ...typing.Type) typing.Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() <= args[1].AsInteger()
			}
		}
		return *typing.NewBoolean(result)
	})

	env.SetCallable(">", func(args ...typing.Type) typing.Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() > args[1].AsInteger()
			}
		}
		return *typing.NewBoolean(result)
	})

	env.SetCallable(">=", func(args ...typing.Type) typing.Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() >= args[1].AsInteger()
			}
		}
		return *typing.NewBoolean(result)
	})

	env.SetCallable("pr-str", func(args ...typing.Type) typing.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		concatenated := strings.Join(parts, " ")
		return typing.Type{String: &concatenated}
	})

	env.SetCallable("str", func(args ...typing.Type) typing.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		concatenated := strings.Join(parts, "")
		return typing.Type{String: &concatenated}
	})

	env.SetCallable("prn", func(args ...typing.Type) typing.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		fmt.Println(strings.Join(parts, " "))
		return *typing.NewNil()
	})

	env.SetCallable("println", func(args ...typing.Type) typing.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		fmt.Println(strings.Join(parts, " "))
		return *typing.NewNil()
	})

	return env
}

func compareIterables(firstList []typing.Type, secondList []typing.Type) bool {
	if len(firstList) != len(secondList) {
		return false
	} else if len(firstList) == 0 {
		return true
	}

	result := true
	for index, _ := range firstList {
		if !compareNodes(firstList[index], secondList[index]) {
			result = false
			break
		}
	}
	return result
}

func compareNodes(first typing.Type, second typing.Type) bool {
	if (first.IsList() || first.IsVector()) && (second.IsList() || second.IsVector()) {
		return compareIterables(first.Iterable(), second.Iterable())
	}

	result := false

	if first.IsNil() && second.IsNil() {
		result = true
	}

	first.IfBoolean(func(a bool) {
		second.IfBoolean(func(b bool) {
			result = a == b
		})
	})

	if first.IsInteger() && second.IsInteger() {
		result = first.AsInteger() == second.AsInteger()
	}

	if first.IsString() && second.IsString() {
		result = first.AsString() == second.AsString()
	}

	return result
}
