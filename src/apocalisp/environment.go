package apocalisp

type Environment struct {
	outer *Environment
	table map[string]ApocalispType
}

func NewEnvironment(outer *Environment) *Environment {
	table := make(map[string]ApocalispType)
	return &Environment{
		table: table,
		outer: outer,
	}
}

func DefaultEnvironment() *Environment {
	env := NewEnvironment(nil)

	env.SetFunction("+", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return ApocalispType{Integer: &r}
	})

	env.SetFunction("-", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return ApocalispType{Integer: &r}
	})

	env.SetFunction("/", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return ApocalispType{Integer: &r}
	})

	env.SetFunction("*", func(inputs ...ApocalispType) ApocalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return ApocalispType{Integer: &r}
	})

	return env
}

func (env *Environment) Set(symbol string, node ApocalispType) {
	env.table[symbol] = node
}

func (env *Environment) SetFunction(symbol string, nativeFunction func(...ApocalispType) ApocalispType) {
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

func (env *Environment) Get(symbol string) ApocalispType {
	if e := env.Find(symbol); e != nil {
		for key, value := range e.table {
			if key == symbol {
				return value
			}
		}
	}

	return ApocalispType{Symbol: &symbol}
}
