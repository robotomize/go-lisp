package hme

import (
	"robotomize/go-lisp/internal/evaluator/environment"
	"robotomize/go-lisp/internal/types"
)

var _ environment.Context = (*Env)(nil)

func New(parent environment.Context) environment.Context {
	return &Env{
		values: make(map[string]types.Value),
		parent: parent,
	}
}

type Env struct {
	values map[string]types.Value
	parent environment.Context
}

func (e *Env) Parent() environment.Context {
	return e.parent
}

func (e *Env) Lookup(name string) (types.Value, bool) {
	for curr := environment.Context(e); curr != nil; curr = curr.Parent() {
		frame, ok := curr.(*Env)
		if !ok {
			return nil, false
		}
		if value, ok := frame.values[name]; ok {
			return value, true
		}
	}

	return nil, false
}

func (e *Env) Define(name string, value types.Value) {
	e.values[name] = value
}

func (e *Env) Mutate(name string, value types.Value) bool {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return true
	}

	if e.parent == nil {
		return false
	}

	return e.parent.Mutate(name, value)
}
