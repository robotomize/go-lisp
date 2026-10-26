//go:build !cons_lambda

package types

import "strings"

var _ Value = (*Cons)(nil)

func NewCons(car, cdr Value) *Cons {
	return &Cons{car: car, cdr: cdr}
}

type Cons struct {
	car Value
	cdr Value
}

func (c *Cons) Car() Value {
	return c.car
}

func (c *Cons) Cdr() Value {
	return c.cdr
}

func (c *Cons) Type() string {
	return TypeCons
}

func (c *Cons) String() string {
	var b strings.Builder
	b.WriteByte('(')
	writeValue(&b, c.Car())

	cdr := c.Cdr()
	for {
		if IsNil(cdr) {
			b.WriteByte(')')
			return b.String()
		}

		next, ok := cdr.(*Cons)
		if !ok {
			b.WriteString(" . ")
			writeValue(&b, cdr)
			b.WriteByte(')')
			return b.String()
		}

		b.WriteByte(' ')
		writeValue(&b, next.Car())
		cdr = next.Cdr()
	}
}

func writeValue(b *strings.Builder, v Value) {
	if v == nil {
		b.WriteString("<nil>")
		return
	}

	b.WriteString(v.String())
}
