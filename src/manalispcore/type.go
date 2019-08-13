package manalispcore

type MalType struct {
	Integer        *int64
	Symbol         *string
	List           *[]MalType
	Vector         *[]MalType
	Hashmap        *[]MalType
	NativeFunction *(func(...*MalType) *MalType)
}

func (m *MalType) IsInteger() bool {
	return m.Integer != nil
}

func (m *MalType) IsSymbol() bool {
	return m.Symbol != nil
}

func (m *MalType) IsList() bool {
	return m.List != nil
}
