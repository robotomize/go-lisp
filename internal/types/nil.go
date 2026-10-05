package types

var _ Value = empty{}

var Nil = empty{}

type empty struct{}

func (n empty) Type() Typ {
	return TypNil
}

func (n empty) String() string {
	return "()"
}
