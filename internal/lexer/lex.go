package lexer

import (
	"strconv"
	"unicode"
)

func New(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
		ptr:    0,
	}
}

type Lexer struct {
	source []rune
	ptr    int
}

func (l *Lexer) Next() (Token, bool) {
	for !l.IsEOF() {
		curr := l.Curr()
		// ignore whitespace
		if unicode.IsSpace(curr) {
			l.ptr++
			continue
		}

		// ignore comment symbol
		if curr == ';' {
			for !l.IsEOF() && l.Curr() != '\n' {
				l.ptr++
			}
			continue
		}

		switch curr {
		case '(':
			l.ptr++
			return Token{Typ: TokenTypOpen, Value: "("}, true
		case ')':
			l.ptr++
			return Token{Typ: TokenTypClose, Value: ")"}, true
		}

		// process string value @TODO need move to separate func
		if curr == '"' {
			l.ptr++
			var output []rune // @TODO optimize, trye to prealloc
			for !l.IsEOF() && l.Curr() != '"' {
				output = append(output, l.Curr())
				l.ptr++
			}

			// pass "
			if !l.IsEOF() {
				l.ptr++
			}

			return Token{
				Typ:   TokenTypeString,
				Value: string(output),
			}, true
		}

		var output []rune
		for !l.IsEOF() {
			curr = l.Curr()
			if unicode.IsSpace(curr) || curr == '"' || curr == '(' || curr == ')' || curr == ';' {
				break
			}
			output = append(output, curr)
			l.ptr++
		}

		strVal := string(output)
		if _, err := strconv.ParseFloat(strVal, 64); err == nil {
			return Token{Typ: TokenTypNumber, Value: strVal}, true
		}

		return Token{Typ: TokenTypSymbol, Value: strVal}, true
	}

	return Token{}, false
}

func (l *Lexer) IsEOF() bool {
	return l.ptr >= len(l.source)
}

func (l *Lexer) Curr() rune {
	if l.ptr >= len(l.source) {
		return 0
	}

	return l.source[l.ptr]
}
