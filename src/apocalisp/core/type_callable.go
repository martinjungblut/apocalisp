package core

func (node *Type) IsCallable() bool {
	return node.Callable != nil
}

func (node *Type) CallCallable(parameters ...Type) Type {
	return (*node.Callable)(parameters...)
}
