package manalispcore

type MalType struct {
	Integer        *int64
	Symbol         *string
	List           *[]MalType
	Vector         *[]MalType
	Hashmap        *[]MalType
	NativeFunction *(func(...*MalType) *MalType)
}

func NewList() *MalType {
	l := make([]MalType, 1)
	return &MalType{List: &l}
}

func (m *MalType) EachInList(callback func(MalType)) {
	for _, t := range *m.List {
		callback(t)
	}
}

func (m *MalType) AddToList(t MalType) {
	*m.List = append(*m.List, t)
}

func (m *MalType) IsInteger() bool {
	return m.Integer != nil
}

func (m *MalType) IsSymbol() bool {
	return m.Symbol != nil
}

func (m *MalType) AsSymbol() string {
	return *m.Symbol
}

func (m *MalType) IsList() bool {
	return m.List != nil
}
