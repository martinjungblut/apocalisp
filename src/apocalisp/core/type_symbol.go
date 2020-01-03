package core

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

func (node *Type) CompareSymbol(others ...string) bool {
	if node.IsSymbol() {
		for _, other := range others {
			if node.AsSymbol() == other {
				return true
			}
		}
	}
	return false
}
