package core

func NewBoolean(value bool) *Type {
	return &Type{Boolean: &value}
}

func (node *Type) IsBoolean() bool {
	return node.Boolean != nil
}

func (node *Type) AsBoolean() bool {
	if node.IsBoolean() {
		return *node.Boolean
	}
	return false
}

func (node *Type) CompareBoolean(value bool) bool {
	if node.IsBoolean() {
		return (*node.Boolean) == value
	}
	return false
}
