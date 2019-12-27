package core

func NewNil() *Type {
	return &Type{Nil: true}
}

func (node *Type) IsNil() bool {
	return node.Nil
}
