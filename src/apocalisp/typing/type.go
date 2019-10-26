package typing

import (
	"fmt"
	"strconv"
	"strings"
)

type Type struct {
	Nil      bool
	Boolean  *bool
	Integer  *int64
	Symbol   *string
	String   *string
	List     *[]Type
	Vector   *[]Type
	Hashmap  *[]Type
	Callable *(func(...Type) Type)
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
		} else if node.IsCallable() {
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
