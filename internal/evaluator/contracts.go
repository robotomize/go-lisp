package evaluator

import "robotomize/go-lisp/internal/types"

type Value interface {
	String() string
	Type() types.Typ
}
