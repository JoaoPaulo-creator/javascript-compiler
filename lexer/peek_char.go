package lexer

func (l *Lexer) peekChar() byte {
	// first check if the read position is greater the length of the input
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}

}
