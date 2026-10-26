package types

import "math/big"

type Numeric interface {
	Add(v Numeric) Numeric
	Sub(v Numeric) Numeric
	Mul(v Numeric) Numeric
	Div(v Numeric) Numeric
	Mod(v Numeric) Numeric
	Cmp(v Numeric) int
	Value
}

var _ Numeric = (*Int)(nil)
var _ Numeric = (*Float)(nil)

func NewInt(value int64) *Int {
	return &Int{val: big.NewInt(value)}
}

func NewIntFrom(v *big.Int) *Int {
	return &Int{val: v}
}

type Int struct {
	val *big.Int
}

func (n *Int) Add(v Numeric) Numeric {
	if m, ok := v.(*Int); ok {
		return &Int{val: new(big.Int).Add(n.val, m.val)}
	}

	return n.ToFloat().Add(v)
}

func (n *Int) Inc() Numeric {
	return n.Add(NewInt(1))
}

func (n *Int) Dec() Numeric {
	return n.Sub(NewInt(1))
}

func (n *Int) Mod(v Numeric) Numeric {
	if IsZero(v) {
		return nil
	}
	if m, ok := v.(*Int); ok {
		return &Int{val: new(big.Int).Mod(n.val, m.val)}
	}
	return n.ToFloat().Mod(v)
}

func (n *Int) Sub(v Numeric) Numeric {
	if m, ok := v.(*Int); ok {
		return &Int{val: new(big.Int).Sub(n.val, m.val)}
	}

	return n.ToFloat().Sub(v)
}

func (n *Int) Mul(v Numeric) Numeric {
	if m, ok := v.(*Int); ok {
		return &Int{val: new(big.Int).Mul(n.val, m.val)}
	}

	return n.ToFloat().Mul(v)
}

func (n *Int) Div(v Numeric) Numeric {
	if IsZero(v) {
		return nil
	}

	if m, ok := v.(*Int); ok {
		return &Int{val: new(big.Int).Div(n.val, m.val)}
	}

	return n.ToFloat().Div(v)
}

func (n *Int) Cmp(v Numeric) int {
	if m, ok := v.(*Int); ok {
		return n.val.Cmp(m.val)
	}

	return n.ToFloat().Cmp(v)
}

func (n *Int) ToFloat() *Float {
	return &Float{val: new(big.Float).SetInt(n.val)}
}

func (n *Int) Type() string {
	return TypeInt
}

func (n Int) String() string {
	return n.val.String()
}

func IsZero(n Numeric) bool {
	switch n := n.(type) {
	case *Int:
		return n.val.Sign() == 0
	case *Float:
		return n.val.Sign() == 0
	default:
		return false
	}
}
