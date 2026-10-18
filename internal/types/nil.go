package types

var _ Value = NilTyp{}

var Nil = NilTyp{}

type NilTyp struct{}

func (n NilTyp) String() string {
	return "()"
}

func (n NilTyp) Type() Typ {
	return TypeNil
}
