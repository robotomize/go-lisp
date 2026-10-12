package parser

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
	OpSConcat
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
	SfPrint
	SfRead
	SfSymbol
	SfDo
)

var ValuesToString = map[any]string{
	OpAdd:     "+",
	OpSub:     "-",
	OpMul:     "*",
	OpDiv:     "/",
	OpMod:     "%",
	OpSConcat: "++",

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
	SfPrint:  "print",
	SfRead:   "read",
	SfSymbol: "symbol",
	SfDo:     "do",
}

var StringToValues = map[string]any{
	"+":  OpAdd,
	"-":  OpSub,
	"*":  OpMul,
	"/":  OpDiv,
	"%":  OpMod,
	"++": OpSConcat,

	"=":  PrEq,
	">":  PrGt,
	"<":  PrLt,
	">=": PrGte,
	"<=": PrLte,

	"quote":  SfQuote,
	"typeof": SfTypeof,
	"cons":   SfCons,
	"car":    SfCar,
	"cdr":    SfCdr,
	"if":     SfIf,
	"print":  SfPrint,
	"read":   SfRead,
	"symbol": SfSymbol,
	"do":     SfDo,
}
