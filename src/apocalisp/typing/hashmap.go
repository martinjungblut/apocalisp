package typing

func NewHashmap() *Type {
	l := make([]Type, 1)
	return &Type{Hashmap: &l}
}

func (node *Type) AddToHashmap(t Type) {
	*node.Hashmap = append(*node.Hashmap, t)
}

func (node *Type) AsHashmap() []Type {
	return *node.Hashmap
}

func (node *Type) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *Type) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}
