package core

func NewAtom(value Type) *Type {
	ptr := &value
	return &Type{Atom: &ptr}
}

func (node *Type) IsAtom() bool {
	return node.Atom != nil
}

func (node *Type) AsAtom() Type {
	if node.IsAtom() {
		return **node.Atom
	}
	return *NewNil()
}

func (node *Type) SetAtom(value Type) {
	*node.Atom = &value
}
