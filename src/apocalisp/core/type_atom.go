package core

func NewAtom(value Type) *Type {
	return &Type{Atom: &value}
}

func (node *Type) IsAtom() bool {
	return node.Atom != nil
}

func (node *Type) AsAtom() Type {
	if node.IsAtom() {
		return *node.Atom
	}
	return *NewNil()
}
