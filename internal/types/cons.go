package types

import "strings"

var _ Value = (*Cons)(nil)

type Cons struct {
	Car Value
	Cdr Value
}

func NewCons(car, cdr Value) *Cons {
	return &Cons{Car: car, Cdr: cdr}
}

func List(elems ...Value) Value {
	if len(elems) == 0 {
		return Nil
	}

	head := &Cons{Car: elems[0], Cdr: Nil}
	tail := head
	for _, elem := range elems[1:] {
		next := &Cons{Car: elem, Cdr: Nil}
		tail.Cdr = next
		tail = next
	}

	return head
}

func (c *Cons) Type() Typ {
	return TypCons
}

func (t *Cons) String() string {
	var b strings.Builder
	b.WriteByte('(')
	writeValue(&b, t.Car)

	cdr := t.Cdr
	for {
		if _, ok := cdr.(empty); ok {
			b.WriteByte(')')
			return b.String()
		}

		if !IsCons(cdr) {
			b.WriteByte(' ')
			writeValue(&b, cdr)
			b.WriteByte(')')
			return b.String()
		}

		next := cdr.(*Cons)

		b.WriteByte(' ')
		writeValue(&b, next.Car)
		cdr = next.Cdr
	}
}

func writeValue(b *strings.Builder, v Value) {
	if v == nil {
		b.WriteString("<nil>")
		return
	}
	b.WriteString(v.String())
}
