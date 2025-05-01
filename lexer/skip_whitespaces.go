package lexer

func (l *Lexer) skipWhiteSpaces() {
	for l.char == ' ' || l.char == '\n' || l.char == '\t' || l.char == '\r' {
		l.readChar()
	}
}
