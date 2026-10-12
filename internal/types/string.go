package types

import "strconv"

var _ Value = (*String)(nil)

func NewString(value string) *String {
	return &String{value: value}
}

type String struct {
	value string
}

func (s String) String() string {
	return strconv.Quote(s.value)
}
