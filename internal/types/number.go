package types

import "fmt"

var _ Value = (*Number)(nil)

func NewNumber(value int64) *Number {
	return &Number{value: value}
}

type Number struct {
	value int64
}

func (n Number) Int() int64 {
	return n.value
}

func (n Number) String() string {
	return fmt.Sprintf("%d", n.value)
}
