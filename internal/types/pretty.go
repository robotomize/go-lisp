package types

import "strings"

func Pretty(v Value) string {
	var b strings.Builder
	writePretty(&b, v, 0)
	return b.String()
}

func writePretty(b *strings.Builder, v Value, indent int) {
	c, ok := v.(*Cons)
	if !ok {
		writeAtom(b, v)
		return
	}

	elems, dotted := consElems(c)
	if dotted != nil || !hasNestedList(elems) {
		b.WriteString(v.String())
		return
	}

	b.WriteByte('(')
	writePretty(b, elems[0], indent)
	for _, elem := range elems[1:] {
		b.WriteByte('\n')
		pad(b, indent+2)
		writePretty(b, elem, indent+2)
	}
	b.WriteByte(')')
}

func writeAtom(b *strings.Builder, v Value) {
	if v == nil {
		b.WriteString("<nil>")
		return
	}
	b.WriteString(v.String())
}

func pad(b *strings.Builder, n int) {
	for i := 0; i < n; i++ {
		b.WriteByte(' ')
	}
}

func hasNestedList(elems []Value) bool {
	for _, elem := range elems {
		if _, ok := elem.(*Cons); ok {
			return true
		}
	}
	return false
}

func consElems(c *Cons) (elems []Value, dotted Value) {
	elems = append(elems, c.Car())
	cdr := c.Cdr()
	for {
		if IsNil(cdr) {
			return elems, nil
		}
		next, ok := cdr.(*Cons)
		if !ok {
			return elems, cdr
		}
		elems = append(elems, next.Car())
		cdr = next.Cdr()
	}
}
