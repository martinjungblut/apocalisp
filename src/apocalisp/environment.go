package apocalisp

import (
	"apocalisp/typing"
	"errors"
	"fmt"
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
		if i < len(nodes) {
			environment.Set(symbols[i], nodes[i])
		}
	}

	return environment
}

func (env *Environment) Set(symbol string, node typing.Type) {
	env.table[symbol] = node
}

func (env *Environment) SetNativeFunction(symbol string, nativeFunction func(...typing.Type) typing.Type) {
	env.table[symbol] = typing.Type{
		NativeFunction: &nativeFunction,
		Symbol:         &symbol,
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

	env.SetNativeFunction("+", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return typing.Type{Integer: &r}
	})

	env.SetNativeFunction("-", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return typing.Type{Integer: &r}
	})

	env.SetNativeFunction("/", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return typing.Type{Integer: &r}
	})

	env.SetNativeFunction("*", func(inputs ...typing.Type) typing.Type {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return typing.Type{Integer: &r}
	})

	return env
}
