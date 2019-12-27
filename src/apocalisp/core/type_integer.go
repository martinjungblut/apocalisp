package core

func (node *Type) IsInteger() bool {
	return node.Integer != nil
}

func (node *Type) AsInteger() int64 {
	return *node.Integer
}
