package typing

func (node *Type) IsString() bool {
	return node.String != nil
}

func (node *Type) AsString() string {
	if node.IsString() {
		return *node.String
	} else {
		return ""
	}
}
