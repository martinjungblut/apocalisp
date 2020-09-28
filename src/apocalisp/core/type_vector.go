package core

func NewVector() *Type {
	slice := make([]Type, 0)
	return &Type{Vector: &slice}
}

func (node *Type) IsVector() bool {
	return node.Vector != nil
}
