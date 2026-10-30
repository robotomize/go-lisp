package evaluator

import (
	"fmt"

	"robotomize/go-lisp/internal/evaluator/environment"
	"robotomize/go-lisp/internal/types"
)

func installPrimitives(env environment.Context, calc Calculator) {
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

	if calc == nil {
		return
	}

	bindFold(env, "+", types.NewInt(0), calc.Add, true)
	bindFold(env, "-", types.NewInt(0), calc.Sub, false)
	bindFold(env, "*", types.NewInt(1), calc.Mul, true)
	bindFold(env, "/", types.NewInt(1), calc.Div, false)
	bindPair(env, "%", calc.Mod)
	bindPair(env, "mod", calc.Mod)
	bindSingle(env, "++", calc.Inc)
	bindSingle(env, "inc", calc.Inc)
	bindSingle(env, "--", calc.Dec)
	bindSingle(env, "dec", calc.Dec)
	bindCmp(env, "=", calc.Cmp, func(order int) bool { return order == 0 })
	bindCmp(env, "eq", calc.Cmp, func(order int) bool { return order == 0 })
	bindCmp(env, "==", calc.Cmp, func(order int) bool { return order == 0 })
	bindCmp(env, ">", calc.Cmp, func(order int) bool { return order > 0 })
	bindCmp(env, "gt", calc.Cmp, func(order int) bool { return order > 0 })
	bindCmp(env, "<", calc.Cmp, func(order int) bool { return order < 0 })
	bindCmp(env, "lt", calc.Cmp, func(order int) bool { return order < 0 })
	bindCmp(env, ">=", calc.Cmp, func(order int) bool { return order >= 0 })
	bindCmp(env, "gte", calc.Cmp, func(order int) bool { return order >= 0 })
	bindCmp(env, "<=", calc.Cmp, func(order int) bool { return order <= 0 })
	bindCmp(env, "lte", calc.Cmp, func(order int) bool { return order <= 0 })
}

func bindFold(env environment.Context, name string, identity types.Numeric, op EvalOpFunc, allowEmpty bool) {
	env.Define(name, types.NewPrimitive(name, func(args []types.Value) (types.Value, error) {
		nums, err := asNumbers(args)
		if err != nil {
			return nil, err
		}

		return cmpMultInstr(name, nums, identity, op, allowEmpty)
	}))
}

func bindPair(env environment.Context, name string, op EvalOpFunc) {
	env.Define(name, types.NewPrimitive(name, func(args []types.Value) (types.Value, error) {
		nums, err := asNumbers(args)
		if err != nil {
			return nil, err
		}

		return cmpPairInstr(name, nums, op)
	}))
}

func bindSingle(env environment.Context, name string, op func(types.Numeric) (types.Numeric, error)) {
	env.Define(name, types.NewPrimitive(name, func(args []types.Value) (types.Value, error) {
		nums, err := asNumbers(args)
		if err != nil {
			return nil, err
		}

		return cmpSingleInstr(name, nums, op)
	}))
}

func bindCmp(env environment.Context, name string, cmp func(a, b types.Numeric) int, pass func(int) bool) {
	env.Define(name, types.NewPrimitive(name, func(args []types.Value) (types.Value, error) {
		nums, err := asNumbers(args)
		if err != nil {
			return nil, err
		}

		return cmpPred(name, nums, cmp, pass)
	}))
}

func asNumbers(args []types.Value) ([]types.Numeric, error) {
	nums := make([]types.Numeric, len(args))
	for i, value := range args {
		number, ok := value.(types.Numeric)
		if !ok {
			return nil, fmt.Errorf("%s is not a number", value)
		}

		nums[i] = number
	}

	return nums, nil
}

func primCons(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, arityError("cons", "expected 2 arguments", args)
	}

	return types.NewCons(args[0], args[1]), nil
}

func primCar(args []types.Value) (types.Value, error) {
	cell, err := onePair("car", args)
	if err != nil {
		return nil, err
	}

	return cell.Car(), nil
}

func primCdr(args []types.Value) (types.Value, error) {
	cell, err := onePair("cdr", args)
	if err != nil {
		return nil, err
	}

	return cell.Cdr(), nil
}

func primList(args []types.Value) (types.Value, error) {
	return types.List(args...), nil
}

func primNull(args []types.Value) (types.Value, error) {
	value, err := oneValue("null?", args)
	if err != nil {
		return nil, err
	}

	return types.NewBool(types.IsNil(value)), nil
}

func primPair(args []types.Value) (types.Value, error) {
	value, err := oneValue("pair?", args)
	if err != nil {
		return nil, err
	}

	return types.NewBool(types.IsCons(value)), nil
}

func primNot(args []types.Value) (types.Value, error) {
	value, err := oneValue("not", args)
	if err != nil {
		return nil, err
	}

	return types.NewBool(isFalse(value)), nil
}

func primEq(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, arityError("eq?", "expected 2 arguments", args)
	}

	return types.NewBool(eqValue(args[0], args[1])), nil
}

func primEqual(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return nil, arityError("equal?", "expected 2 arguments", args)
	}

	return types.NewBool(equalValue(args[0], args[1])), nil
}

func onePair(name string, args []types.Value) (*types.Cons, error) {
	value, err := oneValue(name, args)
	if err != nil {
		return nil, err
	}

	cell, ok := value.(*types.Cons)
	if !ok {
		return nil, fmt.Errorf("%s: expected a pair, got %s: (%s %s)", name, value, name, value)
	}

	return cell, nil
}

func oneValue(name string, args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, arityError(name, "expected 1 argument", args)
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
