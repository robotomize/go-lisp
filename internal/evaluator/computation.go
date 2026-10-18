package evaluator

import (
	"fmt"
	"robotomize/go-lisp/internal/types"
)

type EvalOpFunc func(a, b types.Numeric) (types.Numeric, error)

type Calculator interface {
	Add(a, b types.Numeric) (types.Numeric, error)
	Sub(a, b types.Numeric) (types.Numeric, error)
	Mul(a, b types.Numeric) (types.Numeric, error)
	Div(a, b types.Numeric) (types.Numeric, error)
	Mod(a, b types.Numeric) (types.Numeric, error)
	Inc(a types.Numeric) (types.Numeric, error)
	Dec(a types.Numeric) (types.Numeric, error)
	Cmp(a, b types.Numeric) int
}

type Expression interface {
}

var _ Calculator = (*calcImpl)(nil)

func NewCalculator() Calculator {
	return &calcImpl{}
}

type calcImpl struct{}

func (c *calcImpl) Add(a, b types.Numeric) (types.Numeric, error) {
	return a.Add(b), nil
}

func (c *calcImpl) Sub(a, b types.Numeric) (types.Numeric, error) {
	return a.Sub(b), nil
}

func (c *calcImpl) Mul(a, b types.Numeric) (types.Numeric, error) {
	return a.Mul(b), nil
}

func (c *calcImpl) Div(a, b types.Numeric) (types.Numeric, error) {
	result := a.Div(b)

	if result == nil {
		return nil, fmt.Errorf("division by zero")
	}

	return result, nil
}

func (c *calcImpl) Mod(a, b types.Numeric) (types.Numeric, error) {
	result := a.Mod(b)

	if result == nil {
		return nil, fmt.Errorf("division by zero")
	}

	return result, nil
}

func (c *calcImpl) Cmp(a, b types.Numeric) int {
	return a.Cmp(b)
}

func (c *calcImpl) Inc(a types.Numeric) (types.Numeric, error) {
	return a.Add(types.NewInt(1)), nil
}

func (c *calcImpl) Dec(a types.Numeric) (types.Numeric, error) {
	return a.Sub(types.NewInt(1)), nil
}
