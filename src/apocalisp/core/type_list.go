package core

func NewList() *Type {
	l := make([]Type, 0)
	return &Type{List: &l}
}

func (node *Type) AddToList(t Type) {
	*node.List = append(*node.List, t)
}

func (node *Type) AsList() []Type {
	return *node.List
}

func (node *Type) IsList() bool {
	return node.List != nil
}

func (node *Type) IsEmptyList() bool {
	return node.IsList() && (len(*node.List) == 0)
}
