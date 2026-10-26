package environment

import "robotomize/go-lisp/internal/types"

type Context interface {
	Parent() Context
	Lookup(name string) (types.Value, bool)
	Define(name string, value types.Value)
	Mutate(name string, value types.Value) bool
}

type New func(parent Context) Context
