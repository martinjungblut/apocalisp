package core

import (
	"math/big"
)

func (node *Type) IsInteger() bool {
	return node.Integer != nil
}

func (node *Type) AsInteger() *big.Int {
	return node.Integer
}
