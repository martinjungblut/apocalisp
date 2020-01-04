package core

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

	if node != nil {
		if node.IsNil() {
			return "nil"
		} else if node.IsBoolean() {
			return strconv.FormatBool(node.AsBoolean())
		} else if node.IsInteger() {
			return fmt.Sprintf("%d", node.AsInteger())
		} else if node.IsCallable() || node.IsFunction() {
			return "#<function>"
		} else if node.IsSymbol() {
			return node.AsSymbol()
		} else if node.IsString() {
			return formatString(node.AsString())
		} else if node.IsList() {
			return formatSequence(node.List, "(", ")")
		} else if node.IsVector() {
			return formatSequence(node.Vector, "[", "]")
		} else if node.IsHashmap() {
			return formatSequence(node.Hashmap, "{", "}")
		}
	}
	return ""
}
