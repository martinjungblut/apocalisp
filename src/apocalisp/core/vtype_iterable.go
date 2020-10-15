package core

func (node *Type) IsIterable() bool {
	return node.IsList() || node.IsVector()
}

func (node *Type) IsEvenIterable() bool {
	return node.IsIterable() && len(node.AsIterable())%2 == 0
}

func (node *Type) AsIterable() []Type {
	if node.IsList() {
		return *node.List
	} else if node.IsVector() {
		return *node.Vector
	}
	return make([]Type, 0)
}

func (node *Type) DeriveIterable() *Type {
	if node.IsList() {
		return NewList()
	} else if node.IsVector() {
		return NewVector()
	}
	return nil
}

func (node *Type) Append(t Type) {
	if node.IsList() {
		*node.List = append(*node.List, t)
	} else if node.IsVector() {
		*node.Vector = append(*node.Vector, t)
	}
}

func (node *Type) Prepend(t Type) {
	if node.IsList() {
		*node.List = append([]Type{t}, (*node.List)...)
	} else if node.IsVector() {
		*node.Vector = append([]Type{t}, (*node.Vector)...)
	}
}
