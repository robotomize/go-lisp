package types

type Value interface {
	String() string
	Type() string
}

const (
	TypeNil       string = "nil"
	TypeCons      string = "cons"
	TypeSymbol    string = "symbol"
	TypeInt       string = "int"
	TypeBool      string = "bool"
	TypeFloat     string = "float"
	TypeString    string = "string"
	TypeLambda    string = "lambda"
	TypeDambda    string = "dambda"
	TypePrimitive string = "primitive"
	TypeMacro     string = "macro"
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
