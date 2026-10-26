package types

var _ Value = (*Macro)(nil)

func NewMacro(params []string, rest string, body Value, env any) *Macro {
	return &Macro{
		Procedure: NewProcedure(params, body),
		rest:      rest,
		env:       env,
	}
}

type Macro struct {
	Procedure
	rest string
	env  any
}

func (m *Macro) Rest() string { return m.rest }

func (m *Macro) Env() any { return m.env }

func (m *Macro) String() string { return "#<macro>" }

func (m *Macro) Type() string { return TypeMacro }
