package typing

func (node *Type) IsSymbol() bool {
	return node.Symbol != nil
}

func (node *Type) AsSymbol() string {
	if node.IsSymbol() {
		return *node.Symbol
	} else {
		return ""
	}
}
