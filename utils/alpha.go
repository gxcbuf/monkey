package utils

func IsLiteral(ch rune) bool {
	return IsLetter(ch) || IsDigit(ch)
}

func IsLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func IsWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
