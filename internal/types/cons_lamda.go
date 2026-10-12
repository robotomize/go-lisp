//go:build cons_lambda

package types

import "strings"

type pick func(car, cdr Value) Value

var _ Value = (*Cons)(nil)

func NewCons(car, cdr Value) *Cons {
	return &Cons{
		fn: func(p pick) Value {
			return p(car, cdr)
		},
	}
}

type Cons struct {
	fn func(pick) Value
}

func (c *Cons) Car() Value {
	return c.fn(func(car, _ Value) Value {
		return car
	})
}

func (c *Cons) Cdr() Value {
	return c.fn(func(_, cdr Value) Value {
		return cdr
	})
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
			b.WriteByte(' ')
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
