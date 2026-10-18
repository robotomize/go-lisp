package evaluator

import "robotomize/go-lisp/internal/types"

type env struct {
	bindings map[string]Value
	parent   *env
}

func newEnv(parent *env) *env {
	return &env{
		bindings: map[string]Value{},
		parent:   parent,
	}
}

func (e *env) lookup(name string) (Value, bool) {
	for e != nil {
		if value, ok := e.bindings[name]; ok {
			return value, true
		}
		e = e.parent
	}

	return nil, false
}

func (e *env) define(name string, value Value) {
	e.bindings[name] = value
}

type procedure struct {
	params []string
	body   Value
	env    *env
}

func (p *procedure) String() string {
	return "#<procedure>"
}

func (p *procedure) Type() types.Typ {
	return "procedure"
}
