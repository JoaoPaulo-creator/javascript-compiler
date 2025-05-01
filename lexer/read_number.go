package lexer

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}
