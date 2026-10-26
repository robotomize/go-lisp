package types

var _ Value = (*Primitive)(nil)

type ApplyFunc func([]Value) (Value, error)

func NewPrimitive(name string, fn ApplyFunc) *Primitive {
	return &Primitive{name: name, fn: fn}
}

type Primitive struct {
	name string
	fn   ApplyFunc
}

func (p *Primitive) Apply(args []Value) (Value, error) {
	return p.fn(args)
}

func (p *Primitive) String() string {
	return "#<primitive:" + p.name + ">"
}

func (p *Primitive) Type() string {
	return TypePrimitive
}
