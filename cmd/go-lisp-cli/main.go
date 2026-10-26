package main

import (
	"bufio"
	"fmt"
	"os"

	"robotomize/go-lisp/internal/ast"
	"robotomize/go-lisp/internal/evaluator"
	"robotomize/go-lisp/internal/evaluator/environment"
	"robotomize/go-lisp/internal/evaluator/environment/hme"
	"robotomize/go-lisp/internal/lexer"
	"robotomize/go-lisp/internal/std"
)

func init() {

}

func main() {
	argv := os.Args[1:]

	input := os.Stdin
	isInteractive := true

	if len(argv) > 0 {
		filename := argv[0]
		if _, err := os.Stat(filename); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		defer f.Close()
		input = f
		isInteractive = false
	}

	scanner := bufio.NewScanner(input)
	env := hme.New(nil)
	newEnv := environment.New(hme.New)
	if err := std.Autoload(env, newEnv); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	eval := evaluator.New(evaluator.NewCalculator(), env, newEnv)
	var collected []lexer.Token

	for {
		if isInteractive && len(collected) == 0 {
			fmt.Print("glisp=>")
		}

		if !scanner.Scan() {
			break
		}

		l := lexer.New(scanner.Text())
		for {
			token, ok := l.Next()
			if !ok {
				break
			}
			collected = append(collected, token)
		}

		if !balanced(collected) {
			continue
		}

		if err := evalForms(eval, collected); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			collected = nil
			if !isInteractive {
				os.Exit(1)
			}
			continue
		}
		collected = nil
	}

	if len(collected) > 0 {
		fmt.Fprintln(os.Stderr, "error: expected closing parenthesis ')'")
		os.Exit(1)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func evalForms(eval evaluator.Evaluator, tokens []lexer.Token) error {
	p := ast.New(tokens)
	for !p.IsEOF() {
		form, err := p.Parse()
		if err != nil {
			return err
		}

		if form == nil {
			return nil
		}

		value, err := eval.Eval(form)
		if err != nil {
			return err
		}
		fmt.Println(value)
	}
	return nil
}

// @TODO make this more robust() =()
func balanced(tokens []lexer.Token) bool {
	depth := 0
	for _, token := range tokens {
		switch token.Typ {
		case lexer.TokenTypOpen:
			depth++
		case lexer.TokenTypClose:
			depth--
			if depth < 0 {
				return true
			}
		}
	}
	return depth == 0
}
