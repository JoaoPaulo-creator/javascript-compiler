package lexer

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	// moving forward
	l.position = l.readPosition
	l.readPosition++
}
