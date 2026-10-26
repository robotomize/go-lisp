package types

var _ Value = (*Dambda)(nil)

func NewDambda(params []string, body Value) *Dambda {
	return &Dambda{Procedure: NewProcedure(params, body)}
}

type Dambda struct {
	Procedure
}

func (d *Dambda) String() string { return "#<dambda>" }

func (d *Dambda) Type() string { return TypeDambda }
