package core

func NewHashmap() *Type {
	m := make(map[Type]Type)
	return &Type{Hashmap: &m}
}

func NewHashmapFromSequence(sequence []Type) *Type {
	m := make(map[Type]Type)

	for i := 0; i < len(sequence) && i+1 < len(sequence); i += 2 {
		m[sequence[i]] = sequence[i+1]
	}

	return &Type{Hashmap: &m}
}

func (node *Type) HashmapSet(key Type, value Type) {
	node.AsHashmap()[key] = value
}

func (node *Type) AsHashmap() map[Type]Type {
	return *node.Hashmap
}

func (node *Type) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *Type) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}
