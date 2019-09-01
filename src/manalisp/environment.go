package manalisp

import (
	"errors"
	"fmt"
)

type Environment struct {
	table map[string]ManalispType
}

func NewEnvironment() *Environment {
	table := make(map[string]ManalispType)
	return &Environment{table: table}
}

func (env *Environment) DefineFunction(symbol string, nativeFunction func(...ManalispType) ManalispType) {
	env.table[symbol] = ManalispType{NativeFunction: &nativeFunction}
}

func (env *Environment) Find(symbol string) (ManalispType, error) {
	for key, value := range env.table {
		if key == symbol {
			return value, nil
		}
	}

	return ManalispType{}, errors.New(fmt.Sprintf("Symbol not found: %s", symbol))
}
