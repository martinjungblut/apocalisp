package core

func NewList() *Type {
	slice := make([]Type, 0)
	return &Type{List: &slice}
}

func NewVector() *Type {
	slice := make([]Type, 0)
	return &Type{Vector: &slice}
}

func (node *Type) IsList() bool {
	return node.List != nil
}

func (node *Type) IsVector() bool {
	return node.Vector != nil
}

func (node *Type) IsEmptyList() bool {
	return node.IsList() && len(node.AsIterable()) == 0
}

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
