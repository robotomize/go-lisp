package types

var _ Value = NilTyp{}

var Nil = NilTyp{}

type NilTyp struct{}

func (n NilTyp) String() string {
	return "()"
}
