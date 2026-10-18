package evaluator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"robotomize/go-lisp/internal/ast"
	"robotomize/go-lisp/internal/types"
)

type Evaluator interface {
	Eval(Value) (Value, error)
}

var _ Evaluator = (*evalImpl)(nil)

func New(calculator Calculator) Evaluator {
	return NewWithIO(calculator, os.Stdin, os.Stdout)
}

func NewWithIO(calculator Calculator, in io.Reader, out io.Writer) Evaluator {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}

	return &evalImpl{
		calculator: calculator,
		env:        newEnv(nil),
		in:         bufio.NewReader(in),
		out:        out,
	}
}

type evalImpl struct {
	calculator Calculator
	env        *env
	in         *bufio.Reader
	out        io.Writer
}

func (e *evalImpl) Eval(val Value) (Value, error) {
	switch v := val.(type) {
	case *types.Int, *types.Float, *types.String, *types.Bool, types.NilTyp:
		return v, nil
	case *types.Symbol:
		value, ok := e.env.lookup(v.Name())
		if !ok {
			return nil, fmt.Errorf("unbound symbol: %s", v.Name())
		}
		return value, nil
	case *types.Cons:
		return e.evalCons(v)
	default:
		return nil, fmt.Errorf("cannot eval %T", val)
	}
}

func (e *evalImpl) evalCons(cell *types.Cons) (Value, error) {
	head, ok := cell.Car().(*types.Symbol)
	if !ok {
		return e.evalCall(cell)
	}

	op, known := ast.StringToValues[head.Name()]
	if !known {
		return e.evalCall(cell)
	}

	args := cell.Cdr()
	switch op {
	case ast.SfIf:
		return e.evalIf(args)
	case ast.SfQuote:
		return oneArg(args)
	case ast.SfPrint:
		return e.evalPrint(args)
	case ast.SfRead:
		return e.evalRead(args)
	case ast.SfSymbol:
		return e.evalSymbol(args)
	case ast.SfTypeof:
		return e.evalTypeof(args)
	case ast.SfDo:
		return e.evalDo(args)
	case ast.SfBegin:
		return e.evalSeq(args)
	case ast.SfDefine:
		return e.evalDefine(args)
	case ast.OpAdd, ast.OpSub, ast.OpMul, ast.OpDiv, ast.OpMod, ast.OpInc, ast.OpDec,
		ast.PrEq, ast.PrGt, ast.PrLt, ast.PrGte, ast.PrLte:
		nums, err := e.numbers(args)
		if err != nil {
			return nil, err
		}

		switch op {
		case ast.OpAdd:
			return cmpMultInstr(nums, types.NewInt(0), e.calculator.Add, true)
		case ast.OpSub:
			return cmpMultInstr(nums, types.NewInt(0), e.calculator.Sub, false)
		case ast.OpMul:
			return cmpMultInstr(nums, types.NewInt(1), e.calculator.Mul, true)
		case ast.OpDiv:
			return cmpMultInstr(nums, types.NewInt(1), e.calculator.Div, false)
		case ast.OpMod:
			return cmpPairInstr(nums, e.calculator.Mod)
		case ast.OpInc:
			return cmpSingleInstr(nums, e.calculator.Inc)
		case ast.OpDec:
			return cmpSingleInstr(nums, e.calculator.Dec)
		case ast.PrEq, ast.PrGt, ast.PrLt, ast.PrGte, ast.PrLte:
			return cmpPred(nums, e.calculator.Cmp, op)
		default:
			return nil, fmt.Errorf("operator is not implemented")
		}
	default:
		return nil, fmt.Errorf("operator is not implemented: %s", head.Name())
	}
}

func (e *evalImpl) evalDefine(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}
	if len(items) < 2 {
		return nil, fmt.Errorf("expected 2 arguments")
	}

	switch name := items[0].(type) {
	case *types.Symbol:
		if len(items) != 2 {
			return nil, fmt.Errorf("expected 2 arguments")
		}
		value, err := e.Eval(items[1])
		if err != nil {
			return nil, err
		}
		e.env.define(name.Name(), value)
		return name, nil
	case *types.Cons:
		symbol, ok := name.Car().(*types.Symbol)
		if !ok {
			return nil, fmt.Errorf("define name must be a symbol")
		}
		params, err := paramNames(name.Cdr())
		if err != nil {
			return nil, err
		}
		body := make([]types.Value, len(items)-1)
		for i, expr := range items[1:] {
			body[i] = expr
		}
		e.env.define(symbol.Name(), &procedure{
			params: params,
			body:   types.List(body...),
			env:    e.env,
		})
		return symbol, nil
	default:
		return nil, fmt.Errorf("define name must be a symbol")
	}
}

func paramNames(args Value) ([]string, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(items))
	seen := make(map[string]bool, len(items))
	for i, item := range items {
		symbol, ok := item.(*types.Symbol)
		if !ok {
			return nil, fmt.Errorf("parameter must be a symbol")
		}
		if seen[symbol.Name()] {
			return nil, fmt.Errorf("duplicate parameter: %s", symbol.Name())
		}
		seen[symbol.Name()] = true
		names[i] = symbol.Name()
	}

	return names, nil
}

func (e *evalImpl) evalCall(cell *types.Cons) (Value, error) {
	operator, err := e.Eval(cell.Car())
	if err != nil {
		return nil, err
	}

	proc, ok := operator.(*procedure)
	if !ok {
		return nil, fmt.Errorf("operator must be a procedure")
	}

	args, err := e.evalOperands(cell.Cdr())
	if err != nil {
		return nil, err
	}

	if len(args) != len(proc.params) {
		if len(proc.params) == 1 {
			return nil, fmt.Errorf("expected 1 argument")
		}
		return nil, fmt.Errorf("expected %d arguments", len(proc.params))
	}

	frame := newEnv(proc.env)
	for i, name := range proc.params {
		frame.define(name, args[i])
	}

	saved := e.env
	e.env = frame
	defer func() { e.env = saved }()

	return e.evalSeq(proc.body)
}

func (e *evalImpl) evalOperands(args Value) ([]Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}

	values := make([]Value, len(items))
	for i, expr := range items {
		values[i], err = e.Eval(expr)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

type doVar struct {
	name    string
	init    Value
	step    Value
	hasStep bool
}

func (e *evalImpl) evalDo(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}
	if len(items) < 2 {
		return nil, fmt.Errorf("expected bindings and a test")
	}

	bindings, err := doBindings(items[0])
	if err != nil {
		return nil, err
	}
	exit, err := operands(items[1])
	if err != nil {
		return nil, err
	}
	if len(exit) == 0 {
		return nil, fmt.Errorf("expected a test")
	}

	values := make([]Value, len(bindings))
	for i, binding := range bindings {
		values[i], err = e.Eval(binding.init)
		if err != nil {
			return nil, err
		}
	}

	frame := newEnv(e.env)
	for i, binding := range bindings {
		frame.define(binding.name, values[i])
	}

	saved := e.env
	e.env = frame
	defer func() { e.env = saved }()

	for {
		test, err := e.Eval(exit[0])
		if err != nil {
			return nil, err
		}
		if !isFalse(test) {
			if len(exit) == 1 {
				return test, nil
			}
			return e.evalSeq(valueList(exit[1:]))
		}

		for _, command := range items[2:] {
			if _, err := e.Eval(command); err != nil {
				return nil, err
			}
		}

		next := make([]Value, len(bindings))
		for i, binding := range bindings {
			if !binding.hasStep {
				next[i], _ = e.env.lookup(binding.name)
				continue
			}
			next[i], err = e.Eval(binding.step)
			if err != nil {
				return nil, err
			}
		}
		for i, binding := range bindings {
			e.env.define(binding.name, next[i])
		}
	}
}

func doBindings(spec Value) ([]doVar, error) {
	rows, err := operands(spec)
	if err != nil {
		return nil, err
	}

	bindings := make([]doVar, len(rows))
	seen := make(map[string]bool, len(rows))
	for i, row := range rows {
		parts, err := operands(row)
		if err != nil {
			return nil, err
		}
		if len(parts) < 2 || len(parts) > 3 {
			return nil, fmt.Errorf("do binding must be (var init step)")
		}

		symbol, ok := parts[0].(*types.Symbol)
		if !ok {
			return nil, fmt.Errorf("do variable must be a symbol")
		}
		if seen[symbol.Name()] {
			return nil, fmt.Errorf("duplicate variable: %s", symbol.Name())
		}
		seen[symbol.Name()] = true

		bindings[i] = doVar{name: symbol.Name(), init: parts[1]}
		if len(parts) == 3 {
			bindings[i].step = parts[2]
			bindings[i].hasStep = true
		}
	}

	return bindings, nil
}

func valueList(items []Value) Value {
	body := make([]types.Value, len(items))
	for i, expr := range items {
		body[i] = expr
	}

	return types.List(body...)
}

func (e *evalImpl) evalIf(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}
	if len(items) < 2 || len(items) > 3 {
		return nil, fmt.Errorf("expected 2 or 3 arguments")
	}

	test, err := e.Eval(items[0])
	if err != nil {
		return nil, err
	}

	if isFalse(test) {
		if len(items) == 2 {
			return types.Nil, nil
		}
		return e.Eval(items[2])
	}

	return e.Eval(items[1])
}

func (e *evalImpl) evalPrint(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}

	parts := make([]string, 0, len(items))
	for _, expr := range items {
		value, err := e.Eval(expr)
		if err != nil {
			return nil, err
		}
		parts = append(parts, value.String())
	}

	if _, err := fmt.Fprintln(e.out, strings.Join(parts, " ")); err != nil {
		return nil, err
	}

	return types.Nil, nil
}

func (e *evalImpl) evalRead(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}

	if len(items) != 0 {
		return nil, fmt.Errorf("expected 0 arguments")
	}

	return readForm(e.in)
}

func (e *evalImpl) evalSymbol(args Value) (Value, error) {
	value, err := e.evalOne(args)
	if err != nil {
		return nil, err
	}

	switch v := value.(type) {
	case *types.String:
		return types.NewSymbol(v.V()), nil
	case *types.Symbol:
		return v, nil
	default:
		return nil, fmt.Errorf("expected a string")
	}
}

func (e *evalImpl) evalTypeof(args Value) (Value, error) {
	value, err := e.evalOne(args)
	if err != nil {
		return nil, err
	}

	return types.NewSymbol(value.Type().String()), nil
}

func (e *evalImpl) evalOne(args Value) (Value, error) {
	expr, err := oneArg(args)
	if err != nil {
		return nil, err
	}

	return e.Eval(expr)
}

func (e *evalImpl) evalSeq(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}

	var last Value = types.Nil
	for _, expr := range items {
		last, err = e.Eval(expr)
		if err != nil {
			return nil, err
		}
	}

	return last, nil
}

func oneArg(args Value) (Value, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}
	if len(items) != 1 {
		return nil, fmt.Errorf("expected 1 argument")
	}

	return items[0], nil
}

func operands(args Value) ([]Value, error) {
	var items []Value
	for !types.IsNil(args) {
		cell, ok := args.(*types.Cons)
		if !ok {
			return nil, fmt.Errorf("improper list")
		}
		items = append(items, cell.Car())
		args = cell.Cdr()
	}

	return items, nil
}

func isFalse(value Value) bool {
	boolean, ok := value.(*types.Bool)
	return ok && !boolean.V()
}

func (e *evalImpl) numbers(args Value) ([]types.Numeric, error) {
	items, err := operands(args)
	if err != nil {
		return nil, err
	}

	nums := make([]types.Numeric, 0, len(items))
	for _, expr := range items {
		value, err := e.Eval(expr)
		if err != nil {
			return nil, err
		}

		number, ok := value.(types.Numeric)
		if !ok {
			return nil, fmt.Errorf("%s is not a number", value)
		}
		nums = append(nums, number)
	}

	return nums, nil
}

func cmpMultInstr(nums []types.Numeric, identity types.Numeric, op EvalOpFunc, allowEmpty bool) (Value, error) {
	if len(nums) == 0 {
		if !allowEmpty {
			return nil, fmt.Errorf("expected at least 1 argument")
		}

		return identity, nil
	}

	acc := identity
	rest := nums
	if len(nums) > 1 {
		acc = nums[0]
		rest = nums[1:]
	}

	for _, number := range rest {
		opResult, err := op(acc, number)
		if err != nil {
			return nil, err
		}

		acc = opResult
	}

	return acc, nil
}

func cmpPred(nums []types.Numeric, cmp func(a, b types.Numeric) int, op any) (Value, error) {
	if len(nums) < 2 {
		return nil, fmt.Errorf("expected at least 2 arguments")
	}

	for i := 1; i < len(nums); i++ {
		order := cmp(nums[i-1], nums[i])
		pass := false
		switch op {
		case ast.PrEq:
			pass = order == 0
		case ast.PrLt:
			pass = order < 0
		case ast.PrGt:
			pass = order > 0
		case ast.PrLte:
			pass = order <= 0
		case ast.PrGte:
			pass = order >= 0
		}

		if !pass {
			return types.False, nil
		}
	}

	return types.True, nil
}

func cmpPairInstr(nums []types.Numeric, op EvalOpFunc) (Value, error) {
	if len(nums) != 2 {
		return nil, fmt.Errorf("expected 2 arguments")
	}

	return op(nums[0], nums[1])
}

func cmpSingleInstr(nums []types.Numeric, op func(types.Numeric) (types.Numeric, error)) (Value, error) {
	if len(nums) != 1 {
		return nil, fmt.Errorf("expected 1 argument")
	}

	return op(nums[0])
}
