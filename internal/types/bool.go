package types

var _ Value = (*Bool)(nil)

var (
	True  = &Bool{val: true}
	False = &Bool{val: false}
)

func NewBool(value bool) *Bool {
	if value {
		return True
	}
	return False
}

type Bool struct {
	val bool
}

func (b *Bool) V() bool {
	return b.val
}

func (b *Bool) Type() Typ {
	return TypeBool
}

func (b Bool) String() string {
	if b.val {
		return "#t"
	}

	return "#f"
}
