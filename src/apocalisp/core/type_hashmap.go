package core

type HashmapKey struct {
	Identifier string
	IsSymbol   bool
}

func NewHashmapKey(identifier string, isSymbol bool) HashmapKey {
	return HashmapKey{Identifier: identifier, IsSymbol: isSymbol}
}

func (node *Type) AsHashmapKey() *HashmapKey {
	if node.IsSymbol() {
		return &HashmapKey{Identifier: node.AsSymbol(), IsSymbol: true}
	} else if node.IsString() {
		return &HashmapKey{Identifier: node.AsString(), IsSymbol: false}
	}
	return nil
}

func NewHashmap() *Type {
	m := make(map[HashmapKey]Type)
	return &Type{Hashmap: &m}
}

func NewHashmapFromSequence(sequence []Type) *Type {
	m := make(map[HashmapKey]Type)
	for i := 0; i < len(sequence) && i+1 < len(sequence); i += 2 {
		if key := sequence[i].AsHashmapKey(); key != nil {
			m[*key] = sequence[i+1]
		}
	}
	return &Type{Hashmap: &m}
}

func (node *Type) HashmapSet(key HashmapKey, value Type) {
	node.AsHashmap()[key] = value
}

func (node *Type) AsHashmap() map[HashmapKey]Type {
	return *node.Hashmap
}

func (node *Type) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *Type) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}
