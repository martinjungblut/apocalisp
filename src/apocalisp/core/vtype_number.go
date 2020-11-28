package core

import (
	"math/big"
)

func ParseNumber(s string) (*Type, bool) {
	if f, _, err := big.ParseFloat(s, 10, 4096, big.ToNearestEven); err == nil {
		return &Type{Float: f}, true
	} else {
		return nil, false
	}
}

func NewNumber(v float64) *Type {
	return &Type{Float: big.NewFloat(v)}
}

func (node *Type) IsNumber() bool {
	return node.Integer != nil || node.Float != nil
}

func (node *Type) AsNumber() *big.Float {
	if node.IsInteger() {
		istr := node.AsInteger().String()
		if t, ok := ParseNumber(istr); ok {
			return t.AsNumber()
		}
	} else if node.IsFloat() {
		return new(big.Float).Copy(node.AsFloat())
	}

	return big.NewFloat(0)
}

func (node *Type) CoerceNumber() *Type {
	if node.IsFloat() {
		fstr := node.AsFloat().String()
		if i, ok := new(big.Int).SetString(fstr, 10); ok {
			if f, ok := ParseNumber(i.String()); ok {
				if node.NumberEqual(f.AsNumber()) {
					return &Type{Integer: i}
				}
			}
		}
	}

	return node
}

func (node *Type) NumberLessThan(f *big.Float) bool {
	if n := node.AsNumber(); node.IsNumber() {
		r, _ := new(big.Float).Sub(n, f).Float64()
		return r < 0
	}
	return false
}

func (node *Type) NumberLessEqualThan(f *big.Float) bool {
	if n := node.AsNumber(); node.IsNumber() {
		r, _ := new(big.Float).Sub(n, f).Float64()
		return r <= 0
	}
	return false
}

func (node *Type) NumberGreaterThan(f *big.Float) bool {
	if n := node.AsNumber(); node.IsNumber() {
		r, _ := new(big.Float).Sub(n, f).Float64()
		return r > 0
	}
	return false
}

func (node *Type) NumberGreaterEqualThan(f *big.Float) bool {
	if n := node.AsNumber(); node.IsNumber() {
		r, _ := new(big.Float).Sub(n, f).Float64()
		return r >= 0
	}
	return false
}

func (node *Type) NumberEqual(f *big.Float) bool {
	if n := node.AsNumber(); node.IsNumber() {
		r, _ := new(big.Float).Sub(n, f).Float64()
		return r == 0
	}
	return false
}
