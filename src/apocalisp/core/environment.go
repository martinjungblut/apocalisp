package core

import (
	"errors"
	"fmt"
	"strings"
)

type Environment struct {
	outer *Environment
	table map[string]Type
}

func NewEnvironment(outer *Environment, symbols []string, nodes []Type) *Environment {
	table := make(map[string]Type)

	environment := &Environment{
		table: table,
		outer: outer,
	}

	for i := 0; i < len(symbols); i++ {
		if symbols[i] == "&" {
			rest := nodes[i:]
			if i+1 < len(symbols) {
				environment.Set(symbols[i+1], Type{List: &rest})
			} else {
				environment.Set("&", Type{List: &rest})
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

func (env *Environment) Set(symbol string, node Type) {
	env.table[symbol] = node
}

func (env *Environment) SetCallable(symbol string, callable func(...Type) Type) {
	env.table[symbol] = Type{
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

func (env *Environment) Get(symbol string) (Type, error) {
	if e := env.Find(symbol); e != nil {
		for key, value := range e.table {
			if key == symbol {
				return value, nil
			}
		}
	}
	return Type{}, errors.New(fmt.Sprintf("Error: '%s' not found.", symbol))
}

func DefaultEnvironment() *Environment {
	env := NewEnvironment(nil, []string{}, []Type{})

	env.SetCallable("+", func(inputs ...Type) Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return Type{Integer: &r}
	})

	env.SetCallable("-", func(inputs ...Type) Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return Type{Integer: &r}
	})

	env.SetCallable("/", func(inputs ...Type) Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return Type{Integer: &r}
	})

	env.SetCallable("*", func(inputs ...Type) Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return Type{Integer: &r}
	})

	env.SetCallable("list", func(args ...Type) Type {
		list := NewList()
		for _, arg := range args {
			list.AddToList(arg)
		}
		return *list
	})

	env.SetCallable("list?", func(args ...Type) Type {
		return *NewBoolean(args[0].IsList())
	})

	env.SetCallable("empty?", func(args ...Type) Type {
		var value int64 = int64(len(args[0].AsIterable()))
		return *NewBoolean(value == 0)
	})

	env.SetCallable("count", func(args ...Type) Type {
		var value int64 = int64(len(args[0].AsIterable()))
		return Type{Integer: &value}
	})

	env.SetCallable("=", func(args ...Type) Type {
		if len(args) == 2 {
			return *NewBoolean(args[0].Compare(args[1]))
		}
		return *NewBoolean(false)
	})

	env.SetCallable("<", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() < args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable("<=", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() <= args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable(">", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() > args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable(">=", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() >= args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable("pr-str", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		concatenated := strings.Join(parts, " ")
		return Type{String: &concatenated}
	})

	env.SetCallable("str", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		concatenated := strings.Join(parts, "")
		return Type{String: &concatenated}
	})

	env.SetCallable("prn", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		fmt.Println(strings.Join(parts, " "))
		return *NewNil()
	})

	env.SetCallable("println", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		fmt.Println(strings.Join(parts, " "))
		return *NewNil()
	})

	return env
}
