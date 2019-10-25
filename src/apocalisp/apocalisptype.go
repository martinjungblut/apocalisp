package apocalisp

import (
	"fmt"
	"strconv"
	"strings"
)

type ApocalispType struct {
	Nil            bool
	Boolean        *bool
	Integer        *int64
	Symbol         *string
	String         *string
	List           *[]ApocalispType
	Vector         *[]ApocalispType
	Hashmap        *[]ApocalispType
	NativeFunction *(func(...ApocalispType) ApocalispType)
}

func (node *ApocalispType) ToString() string {
	wrapSequence := func(sequence *[]ApocalispType, lWrap string, rWrap string) string {
		tokens := []string{}
		for _, element := range *sequence {
			if token := element.ToString(); len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
		return fmt.Sprintf("%s%s%s", lWrap, strings.Join(tokens, " "), rWrap)
	}

	repr := ""
	if node != nil {
		node.IfBoolean(func(value bool) {
			repr = strconv.FormatBool(value)
		})
		if node.IsNil() {
			repr = "nil"
		} else if node.IsInteger() {
			repr = fmt.Sprintf("%d", node.AsInteger())
		} else if node.IsNativeFunction() {
			repr = "#<function>"
		} else if node.IsSymbol() {
			repr = node.AsSymbol()
		} else if node.IsString() {
			repr = node.AsString()
		} else if node.IsList() {
			repr = wrapSequence(node.List, "(", ")")
		} else if node.IsVector() {
			repr = wrapSequence(node.Vector, "[", "]")
		} else if node.IsHashmap() {
			repr = wrapSequence(node.Hashmap, "{", "}")
		}
	}
	return repr
}

// nil
func NewNil() *ApocalispType {
	return &ApocalispType{Nil: true}
}

func (node *ApocalispType) IsNil() bool {
	return node.Nil
}

// boolean
func (node *ApocalispType) IfBoolean(callback func(bool)) {
	if node.Boolean != nil {
		callback(*node.Boolean)
	}
}

func (node *ApocalispType) IsBoolean(value bool) bool {
	if node.Boolean != nil {
		return (*node.Boolean) == value
	}
	return false
}

// integer
func (node *ApocalispType) IsInteger() bool {
	return node.Integer != nil
}

func (node *ApocalispType) AsInteger() int64 {
	return *node.Integer
}

// symbol
func (node *ApocalispType) IsSymbol() bool {
	return node.Symbol != nil
}

func (node *ApocalispType) AsSymbol() string {
	return *node.Symbol
}

// string
func (node *ApocalispType) IsString() bool {
	return node.String != nil
}

func (node *ApocalispType) AsString() string {
	return *node.String
}

// list
func NewList() *ApocalispType {
	l := make([]ApocalispType, 1)
	return &ApocalispType{List: &l}
}

func (node *ApocalispType) AddToList(t ApocalispType) {
	*node.List = append(*node.List, t)
}

func (node *ApocalispType) AsList() []ApocalispType {
	return *node.List
}

func (node *ApocalispType) IsList() bool {
	return node.List != nil
}

func (node *ApocalispType) IsEmptyList() bool {
	return node.IsList() && (len(*node.List) == 0)
}

// vector
func NewVector() *ApocalispType {
	l := make([]ApocalispType, 1)
	return &ApocalispType{Vector: &l}
}

func (node *ApocalispType) AddToVector(t ApocalispType) {
	*node.Vector = append(*node.Vector, t)
}

func (node *ApocalispType) AsVector() []ApocalispType {
	return *node.Vector
}

func (node *ApocalispType) IsVector() bool {
	return node.Vector != nil
}

func (node *ApocalispType) IsEmptyVector() bool {
	return node.IsVector() && (len(*node.Vector) == 0)
}

// hashmap
func NewHashmap() *ApocalispType {
	l := make([]ApocalispType, 1)
	return &ApocalispType{Hashmap: &l}
}

func (node *ApocalispType) AddToHashmap(t ApocalispType) {
	*node.Hashmap = append(*node.Hashmap, t)
}

func (node *ApocalispType) AsHashmap() []ApocalispType {
	return *node.Hashmap
}

func (node *ApocalispType) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *ApocalispType) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}

// native function
func (node *ApocalispType) IsNativeFunction() bool {
	return node.NativeFunction != nil
}

func (node *ApocalispType) CallNativeFunction(parameters ...ApocalispType) ApocalispType {
	return (*node.NativeFunction)(parameters...)
}

// TODO: think more about this
func (node *ApocalispType) EvenIterable() bool {
	if node.IsList() {
		return len(*node.List)%2 == 0 && len(*node.List) > 0
	}

	if node.IsVector() {
		return len(*node.Vector)%2 == 0 && len(*node.Vector) > 0
	}

	return false
}

func (node *ApocalispType) Iterable() []ApocalispType {
	if node.IsList() {
		return *node.List
	}

	if node.IsVector() {
		return *node.Vector
	}

	return make([]ApocalispType, 1)
}
