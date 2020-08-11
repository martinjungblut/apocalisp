package core

func (node *Type) IsFloat() bool {
	return node.Float != nil
}

func (node *Type) AsFloat() float64 {
	return *node.Float
}
