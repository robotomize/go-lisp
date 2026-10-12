package parser

import (
	"fmt"
	"strconv"

	"robotomize/go-lisp/internal/lexer"
	"robotomize/go-lisp/internal/types"
)

func New(input []lexer.Token) *Parser {
	return &Parser{
		tokens: input,
		ptr:    0,
	}
}

type Parser struct {
	tokens []lexer.Token
	ptr    int
}

func (p *Parser) Parse() (types.Value, error) {
	if p.IsEOF() {
		return nil, nil
	}

	tok := p.Curr()
	p.ptr++

	switch tok.Typ {
	case lexer.TokenTypNumber:
		val, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("not a number=): %s", tok.Value)
		}
		return types.NewNumber(val), nil

	case lexer.TokenTypeString:
		return types.NewString(tok.Value), nil

	case lexer.TokenTypSymbol:
		return types.NewSymbol(tok.Value), nil

	case lexer.TokenTypOpen:
		var elems []types.Value
		for !p.IsEOF() && p.Curr().Typ != lexer.TokenTypClose {
			elem, err := p.Parse()
			if err != nil {
				return nil, err
			}
			elems = append(elems, elem)
		}

		if p.IsEOF() {
			return nil, fmt.Errorf("expected closing parenthesis ')'")
		}

		p.ptr++
		if len(elems) == 0 {
			return types.Nil, nil
		}

		return types.List(elems...), nil
	case lexer.TokenTypClose:
		return nil, fmt.Errorf("unexpected closing parenthesis ')'")
	}

	return nil, fmt.Errorf("unknown token %s", tok)
}

func (p *Parser) IsEOF() bool {
	return p.ptr >= len(p.tokens)
}

func (p *Parser) Curr() lexer.Token {
	if p.ptr >= len(p.tokens) {
		return lexer.Token{}
	}

	return p.tokens[p.ptr]
}
