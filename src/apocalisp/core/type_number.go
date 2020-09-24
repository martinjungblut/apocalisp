package core

func (node *Type) IsNumber() bool {
	return node.Integer != nil || node.Float != nil
}

func (node *Type) AsNumber() float64 {
	if node.IsInteger() {
		return float64(node.AsInteger())
	}
	if node.IsFloat() {
		return node.AsFloat()
	}
	return 0
}
