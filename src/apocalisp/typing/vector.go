package typing

func NewVector() *Type {
	l := make([]Type, 0)
	return &Type{Vector: &l}
}

func (node *Type) AddToVector(t Type) {
	*node.Vector = append(*node.Vector, t)
}

func (node *Type) AsVector() []Type {
	return *node.Vector
}

func (node *Type) IsVector() bool {
	return node.Vector != nil
}

func (node *Type) IsEmptyVector() bool {
	return node.IsVector() && (len(*node.Vector) == 0)
}
