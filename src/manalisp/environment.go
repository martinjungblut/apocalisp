package manalisp

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

func (env *Environment) DefineFunction(symbol string, nativeFunction func(...MalType) MalType) {
	env.table[symbol] = MalType{NativeFunction: &nativeFunction}
}

func (env *Environment) Find(symbol string) (MalType, error) {
	for key, value := range env.table {
		if key == symbol {
			return value, nil
		}
	}

	return MalType{}, errors.New(fmt.Sprintf("Symbol not found: %s", symbol))
}
