package typing

import (
	"apocalisp/escaping"
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
	Function *Function
}

func (node *Type) ToString(readably bool) string {
	formatSequence := func(sequence *[]Type, lWrap string, rWrap string) string {
		tokens := []string{}
		for _, element := range *sequence {
			if token := element.ToString(readably); len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
		return fmt.Sprintf("%s%s%s", lWrap, strings.Join(tokens, " "), rWrap)
	}

	formatString := func(input string) string {
		if readably {
			return fmt.Sprintf("\"%s\"", escaping.EscapeString(input))
		}
		return input
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
			repr = formatString(node.AsString())
		} else if node.IsList() {
			repr = formatSequence(node.List, "(", ")")
		} else if node.IsVector() {
			repr = formatSequence(node.Vector, "[", "]")
		} else if node.IsHashmap() {
			repr = formatSequence(node.Hashmap, "{", "}")
		}
	}
	return repr
}
