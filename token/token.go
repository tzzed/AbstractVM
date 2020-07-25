package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	EOI     = ";;"

	// Delimiters
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LF        = "\n"

	IDENT     = "IDENT"
	INT       = "INT"
	FLOAT_NUM = "FLOAT_NUM"

	// keywords
	PUSH   = "push"
	POP    = "pop"
	DUMP   = "dump"
	CLEAR  = "clear"
	DUP    = "dup"
	SWAP   = "swap"
	ASSERT = "assert"
	ADD    = "add"
	SUB    = "sub"
	MUL    = "mul"
	DIV    = "div"
	MOD    = "mod"
	LOAD   = "load"
	STORE  = "store"
	PRINT  = "print"
	EXIT   = "exit"

	// TYPES
	INT8       = "int8"
	INT16      = "int16"
	INT32      = "int32"
	FLOAT      = "float"
	FLOAT32    = "float32"
	FLOAT64    = "float64"
	DOUBLE     = "double"
	BIGDECIMAL = "bigdecimal"
)

var keywords = map[string]TokenType{
	"assert": ASSERT,
	"push":   PUSH,
	"store":  STORE,
	"load":   LOAD,
	"pop":    POP,
	"clear":  CLEAR,
	"add":    ADD,
	"sub":    SUB,
	"dump":   DUMP,
	"mul":    MUL,
	"swap":   SWAP,
	"dup":    DUP,
	"div":    DIV,
	"exit":   EXIT,
	"mod":    MOD,
	"print":  PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

func IsIdent(ident string) bool {
	_, ok := keywords[ident]
	return ok
}
