package ast

type SpecForm int

const (
	SfQuote SpecForm = iota
	SfTypeof
	SfIf
	SfCond
	SfPrint
	SfRead
	SfSymbol
	SfDo
	SfDefine
	SfSet
	SfEval
	SfLambda
	SfDambda
	SfLet
	SfDefineMacro
)

var ValuesToString = map[any]string{
	SfQuote:       "quote",
	SfTypeof:      "typeof",
	SfIf:          "if",
	SfCond:        "cond",
	SfPrint:       "print",
	SfRead:        "read",
	SfSymbol:      "symbol",
	SfDo:          "do",
	SfDefine:      "define",
	SfSet:         "set!",
	SfEval:        "eval",
	SfLambda:      "lambda",
	SfDambda:      "dambda",
	SfLet:         "let",
	SfDefineMacro: "define-macro",
}

var StringToValues = map[string]any{
	"quote":        SfQuote,
	"typeof":       SfTypeof,
	"type":         SfTypeof,
	"if":           SfIf,
	"cond":         SfCond,
	"print":        SfPrint,
	"display":      SfPrint,
	"read":         SfRead,
	"symbol":       SfSymbol,
	"do":           SfDo,
	"define":       SfDefine,
	"set!":         SfSet,
	"eval":         SfEval,
	"lambda":       SfLambda,
	"dambda":       SfDambda,
	"let":          SfLet,
	"define-macro": SfDefineMacro,
}
