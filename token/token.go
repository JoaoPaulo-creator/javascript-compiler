package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL" // sera utilizado/chamado quando algum caracter ou algo que nao tenha sido identificado seja utilizado
	EOF     = "EOF"     // end of file

	// operators
	ASSIGN   = "ASSIGN"
	BANG     = "BANG"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	MULTIPLY = "MULTIPLY"
	DIVISION = "DIVISION"
	INT      = "INT"
	LT       = "LESS_THAN"
	GT       = "GREATER_THAN"
	GE       = "GREATER_EQUAL"
	LE       = "LESS_EQUAL"
	NOT_EQ   = "NOT_EQUALS"
	EQ       = "EQUALS"
	IDENT    = "IDENT"
	STRING   = "STRING"

	// delimiters

	// ,
	COMMA = "COMMA"
	// ;
	SEMICOLON = "SEMICOLON"
	// :
	COLON  = "COLON"
	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	// {
	LBRACKET = "LBRACKET"
	// }
	RBRACKET    = "RBRACKET"
	LSQRBRACKET = "LSQRBRACKET"
	RSQRBRACKET = "RSQRBRACKET"

	//keywords
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
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
	"return":   RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
