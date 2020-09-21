package core

type Function struct {
	IsMacro     bool
	Params      []string
	Body        Type
	Callable    (func(...Type) Type)
	Environment Environment
}

func (node *Type) IsFunction() bool {
	return node.Function != nil
}

func (node *Type) IsMacroFunction() bool {
	return node.Function != nil && node.Function.IsMacro
}

func (node *Type) CallFunction(parameters ...Type) Type {
	return (node.Function.Callable)(parameters...)
}
