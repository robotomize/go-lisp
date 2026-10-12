package types

type Value interface {
	String() string
}

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
	_, ok := t.(*Number)
	return ok
}

func List(elems ...Value) Value {
	var result Value = Nil
	for i := len(elems) - 1; i >= 0; i-- {
		result = NewCons(elems[i], result)
	}

	return result
}
