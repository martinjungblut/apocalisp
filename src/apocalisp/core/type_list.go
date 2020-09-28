package core

func NewList(args ...Type) *Type {
	slice := make([]Type, 0)
	for _, arg := range args {
		slice = append(slice, arg)
	}
	return &Type{List: &slice}
}

func (node *Type) IsList() bool {
	return node.List != nil
}

func (node *Type) IsEmptyList() bool {
	return node.IsList() && len(node.AsIterable()) == 0
}
