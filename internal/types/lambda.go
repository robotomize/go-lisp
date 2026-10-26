package types

var _ Value = (*Lambda)(nil)

func NewLambda(params []string, body Value, env any) *Lambda {
	return &Lambda{
		Procedure: NewProcedure(params, body),
		env:       env,
	}
}

type Lambda struct {
	Procedure
	env any
}

func (l *Lambda) Env() any { return l.env }

func (l *Lambda) String() string { return "#<lambda>" }

func (l *Lambda) Type() string { return TypeLambda }
