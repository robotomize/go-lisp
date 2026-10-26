package evaluator

import (
	"fmt"

	"robotomize/go-lisp/internal/evaluator/environment"
	"robotomize/go-lisp/internal/types"
)

func installPrimitives(env environment.Context) {
	if env == nil {
		return
	}

	env.Define("cons", types.NewPrimitive("cons", primCons))
	env.Define("car", types.NewPrimitive("car", primCar))
	env.Define("cdr", types.NewPrimitive("cdr", primCdr))
	env.Define("list", types.NewPrimitive("list", primList))
	env.Define("null?", types.NewPrimitive("null?", primNull))
	env.Define("pair?", types.NewPrimitive("pair?", primPair))
	env.Define("not", types.NewPrimitive("not", primNot))
	env.Define("eq?", types.NewPrimitive("eq?", primEq))
	env.Define("equal?", types.NewPrimitive("equal?", primEqual))
}

func primCons(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments")
	}

	return types.NewCons(args[0], args[1]), nil
}

func primCar(args []types.Value) (types.Value, error) {
	cell, err := onePair(args)
	if err != nil {
		return nil, err
	}

	return cell.Car(), nil
}

func primCdr(args []types.Value) (types.Value, error) {
	cell, err := onePair(args)
	if err != nil {
		return nil, err
	}

	return cell.Cdr(), nil
}

func primList(args []types.Value) (types.Value, error) {
	return types.List(args...), nil
}

func primNull(args []types.Value) (types.Value, error) {
	value, err := oneValue(args)
	if err != nil {
		return nil, err
	}

	return types.NewBool(types.IsNil(value)), nil
}

func primPair(args []types.Value) (types.Value, error) {
	value, err := oneValue(args)
	if err != nil {
		return nil, err
	}

	return types.NewBool(types.IsCons(value)), nil
}

func primNot(args []types.Value) (types.Value, error) {
	value, err := oneValue(args)
	if err != nil {
		return nil, err
	}

	return types.NewBool(isFalse(value)), nil
}

func primEq(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments")
	}

	return types.NewBool(eqValue(args[0], args[1])), nil
}

func primEqual(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments")
	}

	return types.NewBool(equalValue(args[0], args[1])), nil
}

func onePair(args []types.Value) (*types.Cons, error) {
	value, err := oneValue(args)
	if err != nil {
		return nil, err
	}

	cell, ok := value.(*types.Cons)
	if !ok {
		return nil, fmt.Errorf("expected a pair")
	}

	return cell, nil
}

func oneValue(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 argument")
	}

	return args[0], nil
}

func eqValue(a, b types.Value) bool {
	switch x := a.(type) {
	case types.NilTyp:
		_, ok := b.(types.NilTyp)
		return ok
	case *types.Bool:
		y, ok := b.(*types.Bool)
		return ok && x.V() == y.V()
	case *types.Symbol:
		y, ok := b.(*types.Symbol)
		return ok && x.Name() == y.Name()
	default:
		return false
	}
}

func equalValue(a, b types.Value) bool {
	switch x := a.(type) {
	case types.NilTyp:
		_, ok := b.(types.NilTyp)
		return ok
	case *types.Bool:
		y, ok := b.(*types.Bool)
		return ok && x.V() == y.V()
	case *types.Symbol:
		y, ok := b.(*types.Symbol)
		return ok && x.Name() == y.Name()
	case *types.String:
		y, ok := b.(*types.String)
		return ok && x.V() == y.V()
	case *types.Cons:
		y, ok := b.(*types.Cons)
		if !ok {
			return false
		}
		return equalValue(x.Car(), y.Car()) && equalValue(x.Cdr(), y.Cdr())
	case types.Numeric:
		y, ok := b.(types.Numeric)
		return ok && x.Cmp(y) == 0
	default:
		return false
	}
}
