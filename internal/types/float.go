package types

import (
	"math/big"
	"strings"
)

func NewFloat(value float64) *Float {
	return &Float{val: big.NewFloat(value)}
}

func NewFloatFrom(v *big.Float) *Float {
	return &Float{val: v}
}

type Float struct {
	val *big.Float
}

func (f *Float) Add(v Numeric) Numeric {
	return &Float{val: new(big.Float).Add(f.val, FloatOf(v))}
}

func (f *Float) Sub(v Numeric) Numeric {
	return &Float{val: new(big.Float).Sub(f.val, FloatOf(v))}
}

func (f *Float) Mul(v Numeric) Numeric {
	return &Float{val: new(big.Float).Mul(f.val, FloatOf(v))}
}

func (f *Float) Cmp(v Numeric) int {
	return f.val.Cmp(FloatOf(v))
}

func (f *Float) Div(v Numeric) Numeric {
	if IsZero(v) {
		return nil
	}

	return &Float{val: new(big.Float).Quo(f.val, FloatOf(v))}
}

func (f *Float) Mod(v Numeric) Numeric {
	if IsZero(v) {
		return nil
	}

	divisor := FloatOf(v)
	quot := new(big.Float).Quo(f.val, divisor)
	prod := new(big.Float).Mul(divisor, FloorFloat(quot))
	return &Float{val: new(big.Float).Sub(f.val, prod)}
}

func FloorFloat(x *big.Float) *big.Float {
	whole, acc := x.Int(nil)
	if x.Sign() < 0 && acc != big.Exact {
		whole.Sub(whole, big.NewInt(1))
	}

	return new(big.Float).SetInt(whole)
}

func FloatOf(n Numeric) *big.Float {
	switch n := n.(type) {
	case *Float:
		return n.val
	case *Int:
		return new(big.Float).SetInt(n.val)
	default:
		return nil
	}
}

func (f *Float) Type() string {
	return TypeFloat
}

func (f *Float) String() string {
	text := f.val.Text('g', -1)
	if _, acc := f.val.Int(nil); acc == big.Exact && !strings.ContainsAny(text, ".eE") {
		return text + ".0"
	}

	return text
}
