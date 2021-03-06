package core

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

func NewString(content string) *Type {
	return &Type{String: &content}
}
