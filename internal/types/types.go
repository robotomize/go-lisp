package types

type Value interface {
	String() string
	Type() Typ
}

type Typ string

func (s Typ) String() string {
	return string(s)
}

const (
	TypeNil    Typ = "nil"
	TypeCons   Typ = "cons"
	TypeSymbol Typ = "symbol"
	TypeInt    Typ = "int"
	TypeBool   Typ = "bool"
	TypeFloat  Typ = "float"
	TypeString Typ = "string"
)

func IsNil(t Value) bool {
	_, ok := t.(NilTyp)
	return ok
}

func IsCons(t Value) bool {
	_, ok := t.(*Cons)
	return ok
}

func IsSymbol(t Value) bool {
	_, ok := t.(*Symbol)
	return ok
}

func IsNumber(t Value) bool {
	_, ok := t.(*Int)
	return ok
}

func IsBool(t Value) bool {
	_, ok := t.(*Bool)
	return ok
}

func List(elems ...Value) Value {
	var result Value = Nil
	for i := len(elems) - 1; i >= 0; i-- {
		result = NewCons(elems[i], result)
	}

	return result
}
