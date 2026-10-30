package evaluator_test

import (
	"bytes"
	"strings"
	"testing"

	"robotomize/go-lisp/internal/ast"
	"robotomize/go-lisp/internal/evaluator"
	"robotomize/go-lisp/internal/evaluator/environment"
	"robotomize/go-lisp/internal/evaluator/environment/hme"
	"robotomize/go-lisp/internal/lexer"
	"robotomize/go-lisp/internal/std"
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
		{input: "(cond (#t 1))", expected: "1"},
		{input: "(cond (#f 1) (#t 2) (#t 3))", expected: "2"},
		{input: "(cond (#f (/ 1 0)) (#t 2))", expected: "2"},
		{input: "(cond (#t 1) (#f (/ 1 0)))", expected: "1"},
		{input: "(cond (#f 1))", expected: "()"},
		{input: "(cond (5))", expected: "5"},
		{input: "(cond ((> 2 1) 10) (else 0))", expected: "10"},
		{input: "(cond (#f 1) (else 2))", expected: "2"},
		{input: "(cond (#t 1 2))", expected: "2"},
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

		{input: "(do ((i 0 (+ i 1))) ((< i 3) i))", expected: "3"},
		{input: "(do ((a 1 b) (b 2 a) (n 0 (+ n 1))) ((= n 0) b))", expected: "1"},
		{input: "(do ((a 5) (n 0 (+ n 1))) ((< n 2) a))", expected: "5"},
		{input: "(do ((i 0 (+ i 1))) ((< i 2)) (display i))", expected: "()", output: "0\n1\n"},
		{input: "(define (fib-do n) (do ((i n (- i 1)) (a 0 b) (b 1 (+ a b))) ((> i 0) a)))\n(fib-do 10)", expected: "55"},

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
		{input: "(define x 1)\n(set! x 5)\nx", expected: "5"},
		{input: "(define x 1)\n(define (f) (set! x (+ x 1)) x)\n(f)\n(f)\nx", expected: "3"},
		{input: "(define x 1)\n((lambda (x) (set! x 9)) 0)\nx", expected: "1"},
		{input: "(define x 1)\n((lambda (x) (set! x 9) x) 0)", expected: "9"},
		{input: "(eval '(+ 1 2))", expected: "3"},
		{input: "(define expr '(+ 1 2))\n(eval expr)", expected: "3"},
		{input: "(eval 5)", expected: "5"},
		{input: "(define x 10)\n(eval 'x)", expected: "10"},
		{input: "(define x 1)\n((lambda (x) (eval 'x)) 5)", expected: "5"},
		{input: "(eval '(define z 4))\nz", expected: "4"},

		{input: "(typeof (lambda (x) x))", expected: "lambda"},
		{input: "(typeof (dambda (x) x))", expected: "dambda"},
		{input: "((lambda (x) (* x x)) 4)", expected: "16"},
		{input: "((lambda () 1))", expected: "1"},
		{input: "(define (make n) (lambda (x) (+ x n)))\n((make 5) 3)", expected: "8"},
		{input: "(define x 1)\n(define f (lambda () x))\n((lambda (x) (f)) 9)", expected: "1"},
		{input: "(define x 1)\n(define g (dambda () x))\n((lambda (x) (g)) 9)", expected: "9"},
		{input: "(define n 10)\n(define add (lambda (x) (+ x n)))\n((lambda (n) (add 3)) 5)", expected: "13"},
		{input: "(define n 10)\n(define add (dambda (x) (+ x n)))\n((lambda (n) (add 3)) 5)", expected: "8"},

		{input: "(cons 1 2)", expected: "(1 . 2)"},
		{input: "(cons 1 ())", expected: "(1)"},
		{input: "(car (cons 1 2))", expected: "1"},
		{input: "(cdr (cons 1 '(2)))", expected: "(2)"},
		{input: "(car '(1 2))", expected: "1"},
		{input: "(cdr '(1 2))", expected: "(2)"},
		{input: "((lambda (f x) (f x)) car '(1 2))", expected: "1"},
		{input: "(list)", expected: "()"},
		{input: "(list 1 (+ 1 1))", expected: "(1 2)"},

		{input: "(null? ())", expected: "#t"},
		{input: "(null? '(1))", expected: "#f"},
		{input: "(pair? '(1))", expected: "#t"},
		{input: "(pair? (cons 1 2))", expected: "#t"},
		{input: "(pair? ())", expected: "#f"},
		{input: "(pair? 1)", expected: "#f"},
		{input: "(not #f)", expected: "#t"},
		{input: "(not #t)", expected: "#f"},
		{input: "(not 0)", expected: "#f"},
		{input: "(not ())", expected: "#f"},
		{input: "(eq? 'a 'a)", expected: "#t"},
		{input: "(eq? 'a 'b)", expected: "#f"},
		{input: "(eq? #t #t)", expected: "#t"},
		{input: "(eq? () ())", expected: "#t"},
		{input: "(eq? 1 1)", expected: "#f"},
		{input: "(equal? 1 1)", expected: "#t"},
		{input: "(equal? 1 1.0)", expected: "#t"},
		{input: `(equal? "a" "a")`, expected: "#t"},
		{input: "(equal? '(1 2) '(1 2))", expected: "#t"},
		{input: "(equal? '(1 2) '(1 3))", expected: "#f"},
		{input: "(equal? (cons 1 2) (cons 1 2))", expected: "#t"},
		{input: "(equal? (cons 1 2) (cons 1 3))", expected: "#f"},

		{input: "(define add +)\n(add 1 2)", expected: "3"},
		{input: "(let ((+ (lambda (x y) 1))) (+ 2 3))", expected: "1"},
		{input: "(let ((x 2) (y 3)) (+ x y))", expected: "5"},
		{input: "(let () 1)", expected: "1"},
		{input: "(let ((x 1)) x (+ x 2))", expected: "3"},
		{input: "(define x 2)\n(let ((x 1) (y x)) y)", expected: "2"},

		{input: "(define-macro (when test body) (list 'if test body))\n(when #t 1)", expected: "1"},
		{input: "(define-macro (when test body) (list 'if test body))\n(when #f (/ 1 0))", expected: "()"},
		{input: "(define-macro (quote-it x) (list 'quote x))\n(quote-it (+ 1 2))", expected: "(+ 1 2)"},

		{input: "(-)", err: "expected at least 1 argument"},
		{input: "(/)", err: "expected at least 1 argument"},
		{input: "(/ 1 0)", err: "division by zero"},
		{input: "(mod 1)", err: "expected 2 arguments"},
		{input: "(inc)", err: "expected 1 argument"},
		{input: "(= 1)", err: "expected at least 2 arguments"},
		{input: "(if 1)", err: "expected 2 or 3 arguments"},
		{input: "(cond (else 1) (#t 2))", err: "else clause must be last"},
		{input: "(cond ())", err: "cond clause must be (test body)"},
		{input: "n", err: "unbound symbol"},
		{input: "(quote)", err: "expected 1 argument"},
		{input: "(quote 1 2)", err: "expected 1 argument"},
		{input: "(typeof)", err: "expected 1 argument"},
		{input: "(symbol 1)", err: "expected a string"},
		{input: "(read 1)", err: "expected 0 arguments"},
		{input: "(begin 1)", err: "unbound symbol: begin"},
		{input: "(do)", err: "expected bindings and a test"},
		{input: "(define)", err: "expected 2 arguments"},
		{input: "(set!)", err: "expected 2 arguments"},
		{input: "(set! 1 2)", err: "set! name must be a symbol"},
		{input: "(set! missing 1)", err: "unbound symbol: missing"},
		{input: "(eval)", err: "expected 1 argument"},
		{input: "(eval '(+ missing 1))", err: "unbound symbol: missing"},
		{input: "(define 1 2)", err: "define name must be a symbol"},
		{input: "(define (f x x) x)", err: "duplicate parameter"},
		{input: "(define (square x) (* x x))\n(square)", err: "expected 1 argument"},
		{input: "(lambda)", err: "expected a parameter list and a body"},
		{input: "(lambda (x x) x)", err: "duplicate parameter"},
		{input: "(dambda 1 2)", err: "expected a parameter list"},
		{input: "(car 1)", err: "expected a pair"},
		{input: "(car ())", err: "expected a pair"},
		{input: "(cons 1)", err: "expected 2 arguments"},
		{input: "(null?)", err: "expected 1 argument"},
		{input: "(eq? 1)", err: "expected 2 arguments"},
		{input: "(let)", err: "expected bindings and a body"},
		{input: "(let ((x 1)))", err: "expected bindings and a body"},
		{input: "(let ((x 1) (x 2)) x)", err: "duplicate parameter"},
		{input: "(let ((1 2)) 1)", err: "parameter must be a symbol"},
		{input: "(define-macro)", err: "expected a name and a body"},
		{input: "(define-macro when (lambda () 1))", err: "define-macro name must be a list"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var out bytes.Buffer
			eval := evaluator.NewWithIO(evaluator.NewCalculator(), hme.New(nil), environment.New(hme.New), strings.NewReader(tc.stdin), &out)

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

func TestStdLibrary(t *testing.T) {
	env := hme.New(nil)
	newEnv := environment.New(hme.New)
	if err := std.Autoload(env, newEnv); err != nil {
		t.Fatal(err)
	}

	eval := evaluator.New(evaluator.NewCalculator(), env, newEnv)
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "(length ())", expected: "0"},
		{input: "(length '(1 2 3))", expected: "3"},
		{input: "(append '(1) '(2 3))", expected: "(1 2 3)"},
		{input: "(append () '(1))", expected: "(1)"},
		{input: "(reverse '(1 2 3))", expected: "(3 2 1)"},
		{input: "(map (lambda (x) (+ x 1)) '(1 2 3))", expected: "(2 3 4)"},
		{input: "(map car '((1) (2 3)))", expected: "(1 2)"},
		{input: "(filter (lambda (x) (= (% x 2) 0)) '(1 2 3 4))", expected: "(2 4)"},
		{input: "(accumulate (lambda (a b) (+ a b)) 0 '(1 2 3))", expected: "6"},
		{input: "(accumulate + 0 '(1 2 3 4))", expected: "10"},
		{input: "(map inc '(1 2 3))", expected: "(2 3 4)"},
		{input: "(caar '((1 2) 3))", expected: "1"},
		{input: "(cadr '(1 2 3))", expected: "2"},
		{input: "(when #t 1)", expected: "1"},
		{input: "(when #f (/ 1 0))", expected: "()"},
		{input: "(and)", expected: "#t"},
		{input: "(and 1 2 3)", expected: "3"},
		{input: "(and #f (/ 1 0))", expected: "#f"},
		{input: "(and 1 #f 2)", expected: "#f"},
		{input: "(and 0 1)", expected: "1"},
		{input: "(or)", expected: "#f"},
		{input: "(or #f 2 3)", expected: "2"},
		{input: "(or #f #f)", expected: "#f"},
		{input: "(or 0 1)", expected: "0"},
		{input: "(or 1 (/ 1 0))", expected: "1"},
		{input: "(or #f #f 4)", expected: "4"},
		{input: "(unless #f 1)", expected: "1"},
		{input: "(unless #t (/ 1 0))", expected: "()"},
		{input: "(len '(a b c))", expected: "3"},
		{input: "(len ())", expected: "0"},
		{input: "(range 1 5)", expected: "(1 2 3 4 5)"},
		{input: "(range 3 3)", expected: "(3)"},
		{input: "(range 3 1)", expected: "()"},
		{input: "(foldl (lambda (acc x) (+ acc x)) 0 '(1 2 3))", expected: "6"},
		{input: "(foldl + 0 '(1 2 3))", expected: "6"},
		{input: "(reduce + '(1 2 3))", expected: "6"},
		{input: "(foldr (lambda (x acc) (cons x acc)) () '(1 2 3))", expected: "(1 2 3)"},
		{input: "(reduce (lambda (a b) (+ a b)) '(1 2 3))", expected: "6"},
		{input: "(reduce (lambda (a b) (+ a b)) ())", expected: "()"},
		{input: "(take 2 '(1 2 3 4))", expected: "(1 2)"},
		{input: "(take 0 '(1 2))", expected: "()"},
		{input: "(drop 2 '(1 2 3 4))", expected: "(3 4)"},
		{input: "(drop 9 '(1 2))", expected: "()"},
		{input: "(take-while (lambda (x) (< x 3)) '(1 2 3 0))", expected: "(1 2)"},
		{input: "(drop-while (lambda (x) (< x 3)) '(1 2 3 0))", expected: "(3 0)"},
		{input: "(defn square (x) (* x x))\n(square 4)", expected: "16"},
		{input: "(defmacro quote-it (x) (list 'quote x))\n(quote-it (+ 1 2))", expected: "(+ 1 2)"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := evalSource(t, eval, tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != tc.expected {
				t.Fatalf("got %s, want %s", got.String(), tc.expected)
			}
		})
	}
}

func TestReadTwoForms(t *testing.T) {
	eval := evaluator.NewWithIO(evaluator.NewCalculator(), hme.New(nil), environment.New(hme.New), strings.NewReader("42 'x"), nil)
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
