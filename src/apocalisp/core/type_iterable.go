package core

func (node *Type) EvenIterable() bool {
	if node.IsList() {
		return len(*node.List)%2 == 0 && len(*node.List) > 0
	}

	if node.IsVector() {
		return len(*node.Vector)%2 == 0 && len(*node.Vector) > 0
	}

	return false
}

func (node *Type) Iterable() []Type {
	if node.IsList() {
		return *node.List
	}

	if node.IsVector() {
		return *node.Vector
	}

	return make([]Type, 0)
}
