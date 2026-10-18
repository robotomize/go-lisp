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
	SfCons
	SfCar
	SfCdr
	SfIf
	SfCond
	SfPrint
	SfRead
	SfSymbol
	SfDo
	SfBegin
	SfDefine
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

	SfQuote:  "quote",
	SfTypeof: "typeof",
	SfCons:   "cons",
	SfCar:    "car",
	SfCdr:    "cdr",
	SfIf:     "if",
	SfCond:   "cond",
	SfPrint:  "print",
	SfRead:   "read",
	SfSymbol: "symbol",
	SfDo:     "do",
	SfBegin:  "begin",
	SfDefine: "define",
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

	"quote":   SfQuote,
	"typeof":  SfTypeof,
	"type":    SfTypeof,
	"cons":    SfCons,
	"car":     SfCar,
	"cdr":     SfCdr,
	"if":      SfIf,
	"cond":    SfCond,
	"print":   SfPrint,
	"display": SfPrint,
	"read":    SfRead,
	"symbol":  SfSymbol,
	"do":      SfDo,
	"begin":   SfBegin,
	"define":  SfDefine,
}
