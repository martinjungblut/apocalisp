package manalispcore

import (
	"errors"
	"fmt"
)

type Environment struct {
	table map[string]MalType
}

func NewEnvironment() *Environment {
	table := make(map[string]MalType)
	return &Environment{table: table}
}

func (e *Environment) DefineFunction(symbol string, nativeFunction func(...MalType) *MalType) {
	e.table[symbol] = MalType{NativeFunction: &nativeFunction}
}

func (e *Environment) Find(symbol string) (*MalType, error) {
	for k, f := range e.table {
		if k == symbol {
			return &f, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Symbol not found: %s", symbol))
}
