package std

import (
	_ "embed"

	"robotomize/go-lisp/internal/ast"
	"robotomize/go-lisp/internal/evaluator"
	"robotomize/go-lisp/internal/evaluator/environment"
	"robotomize/go-lisp/internal/lexer"
)

//go:embed modules/std.lsp
var std string

func Autoload(env environment.Context, newEnv environment.New) error {
	eval := evaluator.New(evaluator.NewCalculator(), env, newEnv)
	return evalSource(eval, std)
}

func evalSource(eval evaluator.Evaluator, src string) error {
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
	for !parser.IsEOF() {
		form, err := parser.Parse()
		if err != nil {
			return err
		}
		if form == nil {
			break
		}
		if _, err = eval.Eval(form); err != nil {
			return err
		}
	}

	return nil
}
