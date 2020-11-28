package core

import (
	"math/big"
)

func (node *Type) IsFloat() bool {
	return node.Float != nil
}

func (node *Type) AsFloat() *big.Float {
	return new(big.Float).Copy(node.Float)
}
