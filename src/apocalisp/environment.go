package apocalisp

import (
	"errors"
	"fmt"
)

type Environment struct {
	outer *Environment
	table map[string]ApocalispType
}

func NewEnvironment(outer *Environment, symbols []string, nodes []ApocalispType) *Environment {
	table := make(map[string]ApocalispType)

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

func (env *Environment) Set(symbol string, node ApocalispType) {
	env.table[symbol] = node
}

func (env *Environment) SetNativeFunction(symbol string, nativeFunction func(...ApocalispType) ApocalispType) {
	env.table[symbol] = ApocalispType{
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

func (env *Environment) Get(symbol string) (ApocalispType, error) {
	if e := env.Find(symbol); e != nil {
		for key, value := range e.table {
			if key == symbol {
				return value, nil
			}
		}
	}

	return ApocalispType{}, errors.New(fmt.Sprintf("Error: '%s' not found.", symbol))
}

func DefaultEnvironment() *Environment {
	env := NewEnvironment(nil, []string{}, []ApocalispType{})

	env.SetNativeFunction("+", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return ApocalispType{Integer: &r}
	})

	env.SetNativeFunction("-", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return ApocalispType{Integer: &r}
	})

	env.SetNativeFunction("/", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return ApocalispType{Integer: &r}
	})

	env.SetNativeFunction("*", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return ApocalispType{Integer: &r}
	})

	return env
}
