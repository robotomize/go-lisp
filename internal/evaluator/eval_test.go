package evaluator_test

import (
	"bytes"
	"strings"
	"testing"

	"robotomize/go-lisp/internal/ast"
	"robotomize/go-lisp/internal/evaluator"
	"robotomize/go-lisp/internal/lexer"
)

func TestEval(t *testing.T) {
	testCases := []struct {
		input    string
		stdin    string
		expected string
		output   string
		err      string
	}{
		{input: "1", expected: "1"},
		{input: "2.5", expected: "2.5"},
		{input: "#t", expected: "#t"},
		{input: "True", expected: "#t"},
		{input: "false", expected: "#f"},
		{input: `"hi"`, expected: `"hi"`},
		{input: "()", expected: "()"},

		{input: "(+)", expected: "0"},
		{input: "(+ 1 2 3)", expected: "6"},
		{input: "(+ 1 2.5)", expected: "3.5"},
		{input: "(*)", expected: "1"},
		{input: "(* 2 3)", expected: "6"},
		{input: "(* 2 3.0)", expected: "6.0"},

		{input: "(- 5)", expected: "-5"},
		{input: "(- 10 3 2)", expected: "5"},
		{input: "(/ 8 2)", expected: "4"},
		{input: "(/ 5 2.0)", expected: "2.5"},
		{input: "(/ 2)", expected: "0"},

		{input: "(mod 5 3)", expected: "2"},
		{input: "(mod 1 2.0)", expected: "1.0"},
		{input: "(% 5 3)", expected: "2"},
		{input: "(inc 1)", expected: "2"},
		{input: "(dec 3)", expected: "2"},

		{input: "(+ (inc 1) 1 2 3)", expected: "8"},
		{input: "(+ (+ 1 2 3) (+ 1 2 3))", expected: "12"},

		{input: "(= 1 1)", expected: "#t"},
		{input: "(= 1 2)", expected: "#f"},
		{input: "(= 1 1.0)", expected: "#t"},
		{input: "(= 1 1 1)", expected: "#t"},
		{input: "(= 1 1 2)", expected: "#f"},
		{input: "(< 1 2 3)", expected: "#t"},
		{input: "(< 1 3 2)", expected: "#f"},
		{input: "(> 3 2 1)", expected: "#t"},
		{input: "(>= 2 2)", expected: "#t"},
		{input: "(<= 2 1)", expected: "#f"},
		{input: "(< 1 2.5)", expected: "#t"},

		{input: "(if #f 1 2)", expected: "2"},
		{input: "(if 0 1 2)", expected: "1"},
		{input: "(if () 1 2)", expected: "1"},
		{input: "(if #f 1)", expected: "()"},
		{input: "(if (= 1 1) (+ 2 3) 0)", expected: "5"},
		{input: "(if (= 1 2) (/ 1 0) 3)", expected: "3"},
		{input: "(if (eq 1 1) 2 3)", expected: "2"},
		{input: "(if (== 1 2) 1 0)", expected: "0"},
		{input: "(if (lt 2 1) 0 5)", expected: "5"},
		{input: "(if (gt 2 1) 5 0)", expected: "5"},
		{input: "(if (lte 2 2) 7 0)", expected: "7"},
		{input: "(if (gte 1 2) 1 0)", expected: "0"},

		{input: "(quote x)", expected: "x"},
		{input: "(quote (+ 1 2))", expected: "(+ 1 2)"},
		{input: "(quote (quote a))", expected: "(quote a)"},
		{input: "'x", expected: "x"},
		{input: "'(+ 1 2)", expected: "(+ 1 2)"},
		{input: "''a", expected: "(quote a)"},
		{input: "'()", expected: "()"},
		{input: "(quote (print 1))", expected: "(print 1)"},

		{input: "(print 1)", expected: "()", output: "1\n"},
		{input: "(display (+ 1 2))", expected: "()", output: "3\n"},
		{input: `(print "hi" 1)`, expected: "()", output: "\"hi\" 1\n"},
		{input: "(begin (print 1) 2)", expected: "2", output: "1\n"},

		{input: "(read)", stdin: "42", expected: "42"},
		{input: "(read)", stdin: "(+ 1 2)", expected: "(+ 1 2)"},
		{input: "(read)", stdin: "'x", expected: "(quote x)"},
		{input: "(read)", stdin: `"hi"`, expected: `"hi"`},

		{input: `(symbol "foo")`, expected: "foo"},
		{input: `(symbol (symbol "foo"))`, expected: "foo"},
		{input: "(typeof 1)", expected: "int"},
		{input: "(typeof 1.5)", expected: "float"},
		{input: "(typeof #t)", expected: "bool"},
		{input: `(typeof "a")`, expected: "string"},
		{input: "(typeof ())", expected: "nil"},
		{input: "(typeof '(1))", expected: "cons"},
		{input: "(type 'a)", expected: "symbol"},
		{input: "(typeof (+ 1 2))", expected: "int"},

		{input: "(begin 1 2 (+ 3 4))", expected: "7"},
		{input: "(begin 1 2)", expected: "2"},
		{input: "(begin)", expected: "()"},
		{input: "(begin (if #f 1 2) 3)", expected: "3"},
		{input: "(do ((i 0 (+ i 1))) ((= i 3) i))", expected: "3"},
		{input: "(do ((a 1 b) (b 2 a) (n 0 (+ n 1))) ((= n 1) b))", expected: "1"},
		{input: "(do ((a 5) (n 0 (+ n 1))) ((= n 2) a))", expected: "5"},
		{input: "(define (fib-do n) (do ((i n (- i 1)) (a 0 b) (b 1 (+ a b))) ((= i 0) a)))\n(fib-do 10)", expected: "55"},

		{input: "(define x 5)", expected: "x"},
		{input: "(define x 5)\nx", expected: "5"},
		{input: "(define x (+ 1 2))\n(+ x x)", expected: "6"},
		{input: "(define x 1)\n(define x 2)\nx", expected: "2"},
		{input: "(define x '(1 2))\nx", expected: "(1 2)"},
		{input: "(define (square x) (* x x))\n(square 4)", expected: "16"},
		{input: "(define (fact n) (if (= n 0) 1 (* n (fact (- n 1)))))\n(fact 5)", expected: "120"},
		{input: "(define x 1)\n(define (f x) x)\n(f 5)", expected: "5"},
		{input: "(define (f) x)\n(define x 3)\n(f)", expected: "3"},
		{input: "(define (make-adder n) (define (add x) (+ x n)) add)\n(define add5 (make-adder 5))\n(add5 3)", expected: "8"},

		{input: "(-)", err: "expected at least 1 argument"},
		{input: "(/)", err: "expected at least 1 argument"},
		{input: "(/ 1 0)", err: "division by zero"},
		{input: "(mod 1)", err: "expected 2 arguments"},
		{input: "(inc)", err: "expected 1 argument"},
		{input: "(= 1)", err: "expected at least 2 arguments"},
		{input: "(if 1)", err: "expected 2 or 3 arguments"},
		{input: "(cond (#t 1))", err: "not implemented"},
		{input: "n", err: "unbound symbol"},
		{input: "(quote)", err: "expected 1 argument"},
		{input: "(quote 1 2)", err: "expected 1 argument"},
		{input: "(typeof)", err: "expected 1 argument"},
		{input: "(symbol 1)", err: "expected a string"},
		{input: "(read 1)", err: "expected 0 arguments"},
		{input: "(begin (/ 1 0) 1)", err: "division by zero"},
		{input: "(do)", err: "expected bindings and a test"},
		{input: "(define)", err: "expected 2 arguments"},
		{input: "(define 1 2)", err: "define name must be a symbol"},
		{input: "(define (f x x) x)", err: "duplicate parameter"},
		{input: "(define (square x) (* x x))\n(square)", err: "expected 1 argument"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var out bytes.Buffer
			eval := evaluator.NewWithIO(evaluator.NewCalculator(), strings.NewReader(tc.stdin), &out)

			got, err := evalSource(t, eval, tc.input)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("error %v, want %q", err, tc.err)
				}

				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got.String() != tc.expected {
				t.Fatalf("got %s, want %s", got.String(), tc.expected)
			}
			if out.String() != tc.output {
				t.Fatalf("output %q, want %q", out.String(), tc.output)
			}
		})
	}
}

func TestReadTwoForms(t *testing.T) {
	eval := evaluator.NewWithIO(evaluator.NewCalculator(), strings.NewReader("42 'x"), nil)
	first := parseASTFromString(t, "(read)")
	second := parseASTFromString(t, "(read)")

	got, err := eval.Eval(first)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "42" {
		t.Fatalf("got %s, want 42", got.String())
	}

	got, err = eval.Eval(second)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "(quote x)" {
		t.Fatalf("got %s, want (quote x)", got.String())
	}
}

func evalSource(t *testing.T, eval evaluator.Evaluator, src string) (evaluator.Value, error) {
	t.Helper()

	l := lexer.New(src)
	var tokens []lexer.Token
	for {
		token, ok := l.Next()
		if !ok {
			break
		}
		tokens = append(tokens, token)
	}

	parser := ast.New(tokens)
	var last evaluator.Value
	for !parser.IsEOF() {
		form, err := parser.Parse()
		if err != nil {
			t.Fatal(err)
		}
		if form == nil {
			break
		}
		last, err = eval.Eval(form)
		if err != nil {
			return nil, err
		}
	}

	return last, nil
}

func parseASTFromString(t *testing.T, src string) evaluator.Value {
	t.Helper()

	l := lexer.New(src)

	var tokens []lexer.Token
	for {
		token, ok := l.Next()
		if !ok {
			break
		}

		tokens = append(tokens, token)
	}

	tree, err := ast.New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}

	return tree
}
