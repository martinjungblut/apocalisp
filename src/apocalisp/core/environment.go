package core

import (
	"fmt"
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
	for key := range env.table {
		if key == symbol {
			return env
		}
	}

	if env.outer != nil {
		return env.outer.Find(symbol)
	}

	return nil
}

func (env *Environment) Get(symbol string) Type {
	if e := env.Find(symbol); e != nil {
		for key, value := range e.table {
			if key == symbol {
				return value
			}
		}
	}
	return *NewStringException(fmt.Sprintf("'%s' not found", symbol))
}
