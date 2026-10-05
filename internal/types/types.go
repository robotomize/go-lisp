package types

type Typ int8

const (
	TypNil Typ = iota
	TypCons
	TypSymbol
	TypNumber
)

func IsNil(t Value) bool {
	return t.Type() == TypNil
}

func IsCons(t Value) bool {
	return t.Type() == TypCons
}

func IsSymbol(t Value) bool {
	return t.Type() == TypSymbol
}

func IsNumber(t Value) bool {
	return t.Type() == TypNumber
}

type Value interface {
	String() string
	Type() Typ
}
