package typing

func NewBoolean(value bool) *Type {
	return &Type{Boolean: &value}
}

func (node *Type) IfBoolean(callback func(bool)) {
	if node.Boolean != nil {
		callback(*node.Boolean)
	}
}

func (node *Type) IsBoolean(value bool) bool {
	if node.Boolean != nil {
		return (*node.Boolean) == value
	}
	return false
}
