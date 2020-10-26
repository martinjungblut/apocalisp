package core

func NewException(exception Type) *Type {
	return &Type{Exception: &exception}
}

func NewStringException(message string) *Type {
	return &Type{Exception: &Type{String: &message}}
}

func (node *Type) IsException() bool {
	return node.Exception != nil
}

func (node *Type) AsException() *Type {
	return node.Exception
}
