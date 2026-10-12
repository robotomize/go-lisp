package lexer

import "fmt"

type TokenType int

func (t TokenType) String() string {
	switch t {
	case TokenTypOpen:
		return "o"
	case TokenTypClose:
		return "c"
	case TokenTypSymbol:
		return "sym"
	case TokenTypNumber:
		return "n"
	case TokenTypeString:
		return "str"
	default:
		return "unknown"
	}
}

const (
	TokenTypOpen = iota
	TokenTypClose
	TokenTypSymbol
	TokenTypNumber
	TokenTypeString
)

type Token struct {
	Typ   TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("#%s%s", t.Typ, t.Value)
}
