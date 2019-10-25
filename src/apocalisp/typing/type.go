package typing

import (
	"fmt"
	"strconv"
	"strings"
)

type Type struct {
	Nil            bool
	Boolean        *bool
	Integer        *int64
	Symbol         *string
	String         *string
	List           *[]Type
	Vector         *[]Type
	Hashmap        *[]Type
	NativeFunction *(func(...Type) Type)
}

func (node *Type) ToString() string {
	wrapSequence := func(sequence *[]Type, lWrap string, rWrap string) string {
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
func NewNil() *Type {
	return &Type{Nil: true}
}

func (node *Type) IsNil() bool {
	return node.Nil
}

// boolean
func (node *Type) IfBoolean(callback func(bool)) {
	if node.Boolean != nil {
		callback(*node.Boolean)
	}
}

func (node *Type) IsBoolean(value bool) bool {
	if node.Boolean != nil {
		return (*node.Boolean) == value
	}
	return false
}

// integer
func (node *Type) IsInteger() bool {
	return node.Integer != nil
}

func (node *Type) AsInteger() int64 {
	return *node.Integer
}

// symbol
func (node *Type) IsSymbol() bool {
	return node.Symbol != nil
}

func (node *Type) AsSymbol() string {
	return *node.Symbol
}

// string
func (node *Type) IsString() bool {
	return node.String != nil
}

func (node *Type) AsString() string {
	return *node.String
}

// list
func NewList() *Type {
	l := make([]Type, 1)
	return &Type{List: &l}
}

func (node *Type) AddToList(t Type) {
	*node.List = append(*node.List, t)
}

func (node *Type) AsList() []Type {
	return *node.List
}

func (node *Type) IsList() bool {
	return node.List != nil
}

func (node *Type) IsEmptyList() bool {
	return node.IsList() && (len(*node.List) == 0)
}

// vector
func NewVector() *Type {
	l := make([]Type, 1)
	return &Type{Vector: &l}
}

func (node *Type) AddToVector(t Type) {
	*node.Vector = append(*node.Vector, t)
}

func (node *Type) AsVector() []Type {
	return *node.Vector
}

func (node *Type) IsVector() bool {
	return node.Vector != nil
}

func (node *Type) IsEmptyVector() bool {
	return node.IsVector() && (len(*node.Vector) == 0)
}

// hashmap
func NewHashmap() *Type {
	l := make([]Type, 1)
	return &Type{Hashmap: &l}
}

func (node *Type) AddToHashmap(t Type) {
	*node.Hashmap = append(*node.Hashmap, t)
}

func (node *Type) AsHashmap() []Type {
	return *node.Hashmap
}

func (node *Type) IsHashmap() bool {
	return node.Hashmap != nil
}

func (node *Type) IsEmptyHashmap() bool {
	return node.IsHashmap() && (len(*node.Hashmap) == 0)
}

// native function
func (node *Type) IsNativeFunction() bool {
	return node.NativeFunction != nil
}

func (node *Type) CallNativeFunction(parameters ...Type) Type {
	return (*node.NativeFunction)(parameters...)
}

// TODO: think more about this
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

	return make([]Type, 1)
}
