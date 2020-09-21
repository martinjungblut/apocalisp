package core

func NewException() *Type {
	tslice := make([]Type, 0)
	node := Type{Exception: true, Hashmap: &tslice}
	return &node
}

func (node *Type) IsException() bool {
	return node.Exception
}

func (node *Type) SetExceptionMessage(message string) {
	key := "message"
	node.AddToHashmap(Type{String: &key})
	node.AddToHashmap(Type{String: &message})
}

func (node *Type) ExceptionMessage() *string {
	if node.IsException() {
		hashmap := node.AsHashmap()
		for i := 0; i <= len(hashmap); i += 2 {
			if hashmap[i].AsString() == "message" {
				s := hashmap[i+1].AsString()
				return &s
			}
		}
	}
	return nil
}
