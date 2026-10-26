package types

func NewProcedure(params []string, body Value) Procedure {
	return Procedure{
		params: append([]string(nil), params...),
		body:   body,
	}
}

type Procedure struct {
	params []string
	body   Value
}

func (p Procedure) Params() []string { return p.params }

func (p Procedure) Body() Value { return p.body }
