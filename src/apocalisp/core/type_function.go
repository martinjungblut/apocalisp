package core

type Function struct {
	Params      []string
	Body        Type
	Callable    (func(...Type) Type)
	Environment Environment
}

func (node *Type) IsFunction() bool {
	return node.Function != nil
}

func (node *Type) CallFunction(parameters ...Type) Type {
	return (node.Function.Callable)(parameters...)
}
