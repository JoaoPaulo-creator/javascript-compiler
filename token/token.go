package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL" // sera utilizado/chamado quando algum caracter ou algo que nao tenha sido identificado seja utilizado
	EOF     = "EOF"     // end of file

	// identificadores e literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	DOT    = "DOT"

	// operafores
	ASSIGN    = "="
	PLUS      = "+"
	INCREMENT = "++"
	MINUS     = "-"
	DECREMENT = "--"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="

	// delimitadores
	COMMA    = ","
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	WHILE    = "WHILE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string // token literal, ex.: 2, 3... numbers in majority, strings, etc
}

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"const":    CONST,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"return":   RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
