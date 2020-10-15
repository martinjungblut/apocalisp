package core

func NewVector(args ...Type) *Type {
	slice := make([]Type, 0)
	for _, arg := range args {
		slice = append(slice, arg)
	}
	return &Type{Vector: &slice}
}

func (node *Type) IsVector() bool {
	return node.Vector != nil
}
