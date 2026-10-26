package ast

type (
	OpType   int
	PredType int
	SpecForm int
)

const (
	OpAdd OpType = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpInc
	OpDec
)

const (
	PrEq PredType = iota
	PrGt
	PrLt
	PrGte
	PrLte
)

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
	OpAdd: "+",
	OpSub: "-",
	OpMul: "*",
	OpDiv: "/",
	OpMod: "%",
	OpInc: "++",
	OpDec: "--",

	PrEq:  "=",
	PrGt:  ">",
	PrLt:  "<",
	PrGte: ">=",
	PrLte: "<=",

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
	"+":   OpAdd,
	"-":   OpSub,
	"*":   OpMul,
	"/":   OpDiv,
	"%":   OpMod,
	"mod": OpMod,
	"++":  OpInc,
	"inc": OpInc,
	"dec": OpDec,
	"--":  OpDec,

	"=":   PrEq,
	"eq":  PrEq,
	"==":  PrEq,
	">":   PrGt,
	"gt":  PrGt,
	"<":   PrLt,
	"lt":  PrLt,
	">=":  PrGte,
	"gte": PrGte,
	"<=":  PrLte,
	"lte": PrLte,

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
