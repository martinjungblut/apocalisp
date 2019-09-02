package manalisp

type Environment struct {
	table map[string]ManalispType
}

func NewEnvironment() *Environment {
	table := make(map[string]ManalispType)
	return &Environment{table: table}
}

func (env *Environment) DefineFunction(symbol string, nativeFunction func(...ManalispType) ManalispType) {
	env.table[symbol] = ManalispType{
		NativeFunction: &nativeFunction,
		Symbol:         &symbol,
	}
}

func (env *Environment) Find(symbol string) ManalispType {
	for key, value := range env.table {
		if key == symbol {
			return value
		}
	}

	return ManalispType{Symbol: &symbol}
}
