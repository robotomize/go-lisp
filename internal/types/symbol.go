package types

import "unique"

var _ Value = (*Symbol)(nil)

func NewSymbol(name string) *Symbol {
	return &Symbol{
		intern: true,
		handle: unique.Make(name),
	}
}

type Symbol struct {
	intern bool
	handle unique.Handle[string]
}

func (s Symbol) Name() string {
	return s.handle.Value()
}

func (s Symbol) String() string {
	return s.handle.Value()
}

func (s *Symbol) Type() string {
	return TypeSymbol
}
