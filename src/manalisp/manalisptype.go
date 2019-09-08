package manalisp

import (
	"fmt"
	"strings"
)

type ManalispType struct {
	Integer        *int64
	Symbol         *string
	List           *[]ManalispType
	Vector         *[]ManalispType
	Hashmap        *[]ManalispType
	NativeFunction *(func(...ManalispType) ManalispType)
}

func (node *ManalispType) ToString() string {
	wrapSequence := func(sequence *[]ManalispType, lWrap string, rWrap string) string {
		tokens := []string{}
		for _, element := range *sequence {
			if token := element.ToString(); len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
		return fmt.Sprintf("%s%s%s", lWrap, strings.Join(tokens, " "), rWrap)
	}

	if node != nil {
		if node.IsInteger() {
			return fmt.Sprintf("%d", node.AsInteger())
		} else if node.IsSymbol() {
			return node.AsSymbol()
		} else if node.IsList() {
			return wrapSequence(node.List, "(", ")")
		} else if node.IsVector() {
			return wrapSequence(node.Vector, "[", "]")
		} else if node.IsHashmap() {
			return wrapSequence(node.Hashmap, "{", "}")
		} else if node.IsNativeFunction() {
			return node.AsSymbol()
		} else {
			return ""
		}
	} else {
		return ""
	}
}

// integer
func (node *ManalispType) IsInteger() bool {
	return node.Integer != nil
}

func (node *ManalispType) AsInteger() int64 {
	return *node.Integer
}

// symbol
func (node *ManalispType) IsSymbol() bool {
	return node.Symbol != nil
}

func (node *ManalispType) AsSymbol() string {
	return *node.Symbol
}

// list
func NewList() *ManalispType {
	l := make([]ManalispType, 1)
	return &ManalispType{List: &l}
}

func (node *ManalispType) AddToList(t ManalispType) {
	*node.List = append(*node.List, t)
}

func (node *ManalispType) AsList() []ManalispType {
	return *node.List
}

func (node *ManalispType) IsList() bool {
	return node.List != nil
}

func (node *ManalispType) IsEmptyList() bool {
	return node.IsList() && (len(*node.List) == 0)
}

// vector
func NewVector() *ManalispType {
	l := make([]ManalispType, 1)
	return &ManalispType{Vector: &l}
}

func (node *ManalispType) AddToVector(t ManalispType) {
	*node.Vector = append(*node.Vector, t)
}

func (node *ManalispType) AsVector() []ManalispType {
	return *node.Vector
}

func (node *ManalispType) IsVector() bool {
	return node.Vector != nil
}

func (node *ManalispType) IsEmptyVector() bool {
	return node.IsVector() && (len(*node.Vector) == 0)
}

// hashmap
func NewHashmap() *ManalispType {
	l := make([]ManalispType, 1)
	return &ManalispType{Hashmap: &l}
}

func (node *ManalispType) AddToHashmap(t ManalispType) {
	*node.Hashmap = append(*node.Hashmap, t)
}

func (node *ManalispType) AsHashmap() []ManalispType {
	return *node.Hashmap
}

func (node *ManalispType) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *ManalispType) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}

// native function
func (node *ManalispType) IsNativeFunction() bool {
	return node.NativeFunction != nil
}

func (node *ManalispType) CallNativeFunction(parameters ...ManalispType) ManalispType {
	return (*node.NativeFunction)(parameters...)
}
