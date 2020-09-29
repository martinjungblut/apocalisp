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

func (node *Type) CoerceNumber() *Type {
	if node.IsInteger() {
		return node
	} else if node.IsFloat() {
		if coerced := int64(node.AsFloat()); float64(coerced) == node.AsFloat() {
			return &Type{Integer: &coerced}
		} else {
			return node
		}
	}

	return nil
}
