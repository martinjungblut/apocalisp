package typing

func (node *Type) IsCallable() bool {
	return node.Callable != nil
}

func (node *Type) Call(parameters ...Type) Type {
	return (*node.Callable)(parameters...)
}
