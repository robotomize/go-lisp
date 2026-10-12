package parser

func NewParser(input string) *Parser {
	return &Parser{
		text: input,
	}
}

type Parser struct {
	text string
}

func (p *Parser) Parse() {

}
