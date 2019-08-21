package manalispcore

import (
	"fmt"
	"strings"
)

type MalType struct {
	Integer        *int64
	Symbol         *string
	List           *[]MalType
	Vector         *[]MalType
	Hashmap        *[]MalType
	NativeFunction *(func(...MalType) MalType)
}

func (m *MalType) IsInteger() bool {
	return m.Integer != nil
}

func (m *MalType) AsInteger() int64 {
	return *m.Integer
}

func (m *MalType) IsSymbol() bool {
	return m.Symbol != nil
}

func (m *MalType) AsSymbol() string {
	return *m.Symbol
}

func NewList() *MalType {
	l := make([]MalType, 1)
	return &MalType{List: &l}
}

func (m *MalType) AddToList(t MalType) {
	*m.List = append(*m.List, t)
}

func (m *MalType) AsList() []MalType {
	return *m.List
}

func (m *MalType) IsList() bool {
	return m.List != nil
}

func (m *MalType) IsEmptyList() bool {
	if m.IsList() && (len(*m.List) == 0) {
		return true
	} else {
		return false
	}
}

func (m *MalType) IsVector() bool {
	return m.Vector != nil
}

func (m *MalType) IsHashmap() bool {
	return m.Hashmap != nil
}

func (m *MalType) ToString() string {
	wrapSequence := func(sequence *[]MalType, lWrap string, rWrap string) string {
		tokens := []string{}
		for _, element := range *sequence {
			if token := element.ToString(); len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
		return fmt.Sprintf("%s%s%s", lWrap, strings.Join(tokens, " "), rWrap)
	}

	if m != nil {
		if m.IsInteger() {
			return fmt.Sprintf("%d", m.AsInteger())
		} else if m.IsSymbol() {
			return m.AsSymbol()
		} else if m.IsList() {
			return wrapSequence(m.List, "(", ")")
		} else if m.IsVector() {
			return wrapSequence(m.Vector, "[", "]")
		} else if m.IsHashmap() {
			return wrapSequence(m.Hashmap, "{", "}")
		} else {
			return ""
		}
	} else {
		return ""
	}
}
