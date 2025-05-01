package lexer

import (
	"js-compiler/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // current reading position in input
	char         byte // current char under examination
}

// Builds the lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpaces()

	switch l.char {
	case '+':
		tok = newToken(token.PLUS, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			newToken(token.ILLEGAL, l.char)
		}

	}

	l.readChar()
	return tok
}

func (l *Lexer) Tokenize() []token.Token {
	var tokens []token.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}
