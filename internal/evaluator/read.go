package evaluator

import (
	"bufio"
	"fmt"
	"io"
	"unicode"

	"robotomize/go-lisp/internal/ast"
	"robotomize/go-lisp/internal/lexer"
	"robotomize/go-lisp/internal/types"
)

func readForm(in *bufio.Reader) (types.Value, error) {
	var buf []rune
	err := readExpr(in, &buf)
	if err == io.EOF {
		return nil, fmt.Errorf("end of file")
	}
	if err != nil {
		return nil, err
	}

	return parseForm(string(buf))
}

func parseForm(src string) (types.Value, error) {
	l := lexer.New(src)
	var tokens []lexer.Token
	for {
		token, ok := l.Next()
		if !ok {
			break
		}
		tokens = append(tokens, token)
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("end of file")
	}

	return ast.New(tokens).Parse()
}

func readExpr(in *bufio.Reader, buf *[]rune) error {
	if err := skipSpace(in, buf); err != nil {
		return err
	}

	rn, err := peekRune(in)
	if err != nil {
		return err
	}

	switch rn {
	case '\'':
		if err := takeRune(in, buf); err != nil {
			return err
		}
		err := readExpr(in, buf)
		if err == io.EOF {
			return fmt.Errorf("expected expression after quote")
		}
		return err
	case '(':
		return readList(in, buf)
	case '"':
		return readString(in, buf)
	case ')':
		return fmt.Errorf("unexpected closing parenthesis ')'")
	default:
		return readAtom(in, buf)
	}
}

func readList(in *bufio.Reader, buf *[]rune) error {
	if err := takeRune(in, buf); err != nil {
		return err
	}

	for {
		if err := skipSpace(in, buf); err != nil {
			return err
		}

		rn, err := peekRune(in)
		if err == io.EOF {
			return fmt.Errorf("expected closing parenthesis ')'")
		}
		if err != nil {
			return err
		}
		if rn == ')' {
			return takeRune(in, buf)
		}

		if err := readExpr(in, buf); err != nil {
			if err == io.EOF {
				return fmt.Errorf("expected closing parenthesis ')'")
			}
			return err
		}
	}
}

func readString(in *bufio.Reader, buf *[]rune) error {
	if err := takeRune(in, buf); err != nil {
		return err
	}

	for {
		rn, err := peekRune(in)
		if err == io.EOF {
			return fmt.Errorf("expected closing quote '\"'")
		}
		if err != nil {
			return err
		}
		if err := takeRune(in, buf); err != nil {
			return err
		}
		if rn == '"' {
			return nil
		}
	}
}

func readAtom(in *bufio.Reader, buf *[]rune) error {
	for {
		rn, err := peekRune(in)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if unicode.IsSpace(rn) || rn == '(' || rn == ')' || rn == '\'' || rn == '"' || rn == ';' {
			return nil
		}
		if err := takeRune(in, buf); err != nil {
			return err
		}
	}
}

func skipSpace(in *bufio.Reader, buf *[]rune) error {
	for {
		rn, err := peekRune(in)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if rn == ';' {
			for {
				comment, err := peekRune(in)
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
				if err := takeRune(in, buf); err != nil {
					return err
				}
				if comment == '\n' {
					break
				}
			}
			continue
		}
		if !unicode.IsSpace(rn) {
			return nil
		}
		if err := takeRune(in, buf); err != nil {
			return err
		}
	}
}

func peekRune(in *bufio.Reader) (rune, error) {
	rn, _, err := in.ReadRune()
	if err != nil {
		return 0, err
	}
	if err := in.UnreadRune(); err != nil {
		return 0, err
	}

	return rn, nil
}

func takeRune(in *bufio.Reader, buf *[]rune) error {
	rn, _, err := in.ReadRune()
	if err != nil {
		return err
	}
	*buf = append(*buf, rn)

	return nil
}
