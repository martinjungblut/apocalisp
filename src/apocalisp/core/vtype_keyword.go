package core

import (
	"fmt"
	"strings"
)

func (node *Type) IsKeyword() bool {
	return node.IsSymbol() && strings.HasPrefix(node.AsSymbol(), ":")
}

func (node *Type) ToKeyword() (bool, *Type) {
	if node.IsKeyword() {
		return true, node
	} else if node.IsSymbol() {
		return true, NewSymbol(fmt.Sprintf(":%s", node.AsSymbol()))
	} else if node.IsString() {
		if strings.HasPrefix(node.AsString(), ":") {
			return true, NewSymbol(node.AsString())
		} else {
			return true, NewSymbol(fmt.Sprintf(":%s", node.AsString()))
		}
	}
	return false, node
}
