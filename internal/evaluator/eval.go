package evaluator

import (
	"fmt"

	"robotomize/go-lisp/internal/parser"
	"robotomize/go-lisp/internal/types"
)

type Evaluator interface {
	Eval(Value) (Value, error)
}

var _ Evaluator = (*evalImpl)(nil)

func New() Evaluator {
	return &evalImpl{}
}

type evalImpl struct{}

func (e *evalImpl) Eval(val Value) (Value, error) {
	switch v := val.(type) {
	case *types.Number, *types.String, types.NilTyp:
		return v, nil
	case *types.Symbol:
		return nil, fmt.Errorf("unbound symbol: %s", v.Name())
	case *types.Cons:
		return e.evalCons(v)
	default:
		return nil, fmt.Errorf("cannot eval %T", val)
	}
}

func (e *evalImpl) evalCons(cell *types.Cons) (Value, error) {
	head, ok := cell.Car().(*types.Symbol)
	if !ok {
		return nil, fmt.Errorf("operator must be a symbol")
	}

	op, known := parser.StringToValues[head.Name()]
	if !known {
		return nil, fmt.Errorf("unknown operator: %s", head.Name())
	}

	switch op {
	case parser.OpAdd:
		return e.evalAdd(cell.Cdr())
	default:
		return nil, fmt.Errorf("operator is not implemented: %s", head.Name())
	}
}

func (e *evalImpl) evalAdd(args types.Value) (Value, error) {
	var sum int64

	for !types.IsNil(args) {
		cell, ok := args.(*types.Cons)
		if !ok {
			return nil, fmt.Errorf("+: improper list")
		}

		// left side
		value, err := e.Eval(cell.Car())
		if err != nil {
			return nil, err
		}

		// need impl number operation
		number, ok := value.(*types.Number)
		if !ok {
			return nil, fmt.Errorf("+: not a number: %s", value)
		}

		sum += number.Int()
		args = cell.Cdr()
	}

	return types.NewNumber(sum), nil
}
