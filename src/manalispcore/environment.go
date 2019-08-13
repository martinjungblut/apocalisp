package manalispcore

type Environment struct {
	table map[string]*(func(...*MalType) *MalType)
}

func NewEnvironment() *Environment {
	table := make(map[string]*(func(...*MalType) *MalType))
	return &Environment{table: table}
}

func (e *Environment) Define(n string, c func(...*MalType) *MalType) {
	e.table[n] = &c
}
