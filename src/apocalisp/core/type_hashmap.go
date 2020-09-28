package core

func NewHashmap() *Type {
	m := make(map[string]Type)
	return &Type{Hashmap: &m}
}

func NewHashmapFromSequence(sequence []Type) *Type {
	m := make(map[string]Type)

	for i := 0; i < len(sequence) && i+1 < len(sequence); i += 2 {
		if sequence[i].IsString() {
			m[sequence[i].AsString()] = sequence[i+1]
		} else if sequence[i].IsSymbol() {
			sequence[i+1].HashmapSymbolValue = true
			m[sequence[i].AsSymbol()] = sequence[i+1]
		}
	}

	return &Type{Hashmap: &m}
}

func (node *Type) HashmapSet(key Type, value Type) {
	if key.IsString() {
		node.AsHashmap()[key.AsString()] = value
	} else if key.IsSymbol() {
		node.AsHashmap()[key.AsSymbol()] = value
	}
}

func (node *Type) AsHashmap() map[string]Type {
	return *node.Hashmap
}

func (node *Type) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *Type) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}
