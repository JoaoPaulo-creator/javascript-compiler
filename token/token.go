package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL" // sera utilizado/chamado quando algum caracter ou algo que nao tenha sido identificado seja utilizado
	EOF     = "EOF"     // end of file

	PLUS = "+"
	INT  = "INT"

	IDENT = "IDENT"

	//keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string // token literal, ex.: 2, 3... numbers in majority
}

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
