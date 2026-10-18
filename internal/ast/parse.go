package ast

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

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
		return ParseNumeric(tok.Value)

	case lexer.TokenTypeString:
		return types.NewString(tok.Value), nil

	case lexer.TokenTypSymbol:
		if boolean, ok := ParseBool(tok.Value); ok {
			return boolean, nil
		}
		return types.NewSymbol(tok.Value), nil

	case lexer.TokenTypQuote:
		elem, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("Parse(): %w", err)
		}
		if elem == nil {
			return nil, fmt.Errorf("expected expression after quote")
		}

		return types.List(types.NewSymbol("quote"), elem), nil

	case lexer.TokenTypOpen:
		var elems []types.Value
		for !p.IsEOF() && p.Curr().Typ != lexer.TokenTypClose {
			elem, err := p.Parse()
			if err != nil {
				return nil, fmt.Errorf("Parse(): %w", err)
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

func ParseNumeric(s string) (types.Numeric, error) {
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid number: %s", s)
	}

	if len(parts) == 1 {
		if _, err := strconv.ParseInt(s, 10, 64); err != nil {
			return nil, fmt.Errorf("%s is not a number", s)
		}

		v := new(big.Int)
		v.SetString(s, 10)
		return types.NewIntFrom(v), nil
	}

	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return nil, fmt.Errorf("%s is not a number", s)
	}

	v := new(big.Float)
	v.SetString(s)
	return types.NewFloatFrom(v), nil
}

func ParseBool(s string) (*types.Bool, bool) {
	switch strings.ToLower(s) {
	case "#t", "#true", "true":
		return types.True, true
	case "#f", "#false", "false":
		return types.False, true
	default:
		return nil, false
	}
}
