package core

import (
	"apocalisp/escaping"
	"fmt"
	"strconv"
	"strings"
)

type Type struct {
	Nil       bool
	Exception *Type
	Boolean   *bool
	Integer   *int64
	Float     *float64
	Symbol    *string
	String    *string
	List      *[]Type
	Vector    *[]Type
	Hashmap   *map[Type]Type
	Callable  *(func(...Type) Type)
	Function  *Function
	Atom      **Type
}

func (node Type) ToString(readably bool) string {
	formatSequence := func(sequence *[]Type, lWrap string, rWrap string) string {
		tokens := []string{}
		for _, element := range *sequence {
			if token := element.ToString(readably); len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
		return fmt.Sprintf("%s%s%s", lWrap, strings.Join(tokens, " "), rWrap)
	}

	hashmapToSequence := func(node Type) *[]Type {
		sequence := make([]Type, 0)
		for key, value := range node.AsHashmap() {
			sequence = append(sequence, key)
			sequence = append(sequence, value)
		}
		return &sequence
	}

	formatString := func(input string) string {
		if readably {
			return fmt.Sprintf("\"%s\"", escaping.EscapeString(input))
		}
		return input
	}

	if node.IsNil() {
		return "nil"
	} else if node.IsException() {
		return fmt.Sprintf("Exception: %s", node.AsException().ToString(readably))
	} else if node.IsBoolean() {
		return strconv.FormatBool(node.AsBoolean())
	} else if node.IsInteger() {
		return fmt.Sprintf("%d", node.AsInteger())
	} else if node.IsFloat() {
		return fmt.Sprintf("%f", node.AsFloat())
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
		return formatSequence(hashmapToSequence(node), "{", "}")
	} else if node.IsAtom() {
		return fmt.Sprintf("(atom %s)", node.AsAtom().ToString(readably))
	}
	return ""
}

func (node Type) Compare(other Type) bool {
	return compare(node, other)
}

func compare(first Type, second Type) bool {
	if (first.IsList() || first.IsVector()) && (second.IsList() || second.IsVector()) {
		return compareIterables(first.AsIterable(), second.AsIterable())
	}

	if first.IsNil() && second.IsNil() {
		return true
	}

	if first.IsBoolean() && second.IsBoolean() {
		return first.AsBoolean() == second.AsBoolean()
	}

	if first.IsInteger() && second.IsInteger() {
		return first.AsInteger() == second.AsInteger()
	}

	if first.IsFloat() && second.IsFloat() {
		return first.AsFloat() == second.AsFloat()
	}

	if first.IsString() && second.IsString() {
		return first.AsString() == second.AsString()
	}

	if first.IsSymbol() && second.IsSymbol() {
		return first.AsSymbol() == second.AsSymbol()
	}

	return false
}

func compareIterables(firstList []Type, secondList []Type) bool {
	if len(firstList) != len(secondList) {
		return false
	} else if len(firstList) == 0 {
		return true
	}

	for index, _ := range firstList {
		if !compare(firstList[index], secondList[index]) {
			return false
		}
	}
	return true
}
