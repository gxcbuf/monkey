package lexer

import (
	"monkey/token"
	"monkey/utils"
)

type Lexer struct {
	input string

	endtoken *rune

	readPosition int
}

func New(input string) *Lexer {
	ch := rune(0)
	return &Lexer{
		input:    input,
		endtoken: &ch,
	}
}

func (l *Lexer) WithEndToken(ch rune) *Lexer {
	l.endtoken = &ch
	return l
}

func (l *Lexer) NextToken() *token.Token {
	if !l.hasNext() {
		return nil
	}

	l.skipWhiteSpace()

	ch := l.next()

	switch ch {
	case '=':
		if l.peek() == '=' {
			return token.New(token.EQ, string(ch)+string(l.next()))
		}
		return token.New(token.ASSIGN, string(ch))
	case '!':
		if l.peek() == '=' {
			return token.New(token.NOT_EQ, string(ch)+string(l.next()))
		}
		return token.New(token.BANG, string(ch))
	case '+':
		return token.New(token.PLUS, string(ch))
	case '-':
		return token.New(token.MINUS, string(ch))
	case '*':
		return token.New(token.ASTERISK, string(ch))
	case '/':
		return token.New(token.SLASH, string(ch))
	case '<':
		return token.New(token.LT, string(ch))
	case '>':
		return token.New(token.GT, string(ch))
	case ';':
		return token.New(token.SEMICOLON, string(ch))
	case ',':
		return token.New(token.COMMA, string(ch))
	case '{':
		return token.New(token.LBRACE, string(ch))
	case '}':
		return token.New(token.RBRACE, string(ch))
	case '(':
		return token.New(token.LPAREN, string(ch))
	case ')':
		return token.New(token.RPAREN, string(ch))
	case rune(0):
		return token.New(token.EOF, string(ch))
	case '"', '\'':
		return l.readString(ch)
	default:
		if !utils.IsLiteral(ch) {
			return token.NewILLEGAL(string(ch))
		}
		if utils.IsLetter(ch) {
			return l.readIdentifier(ch)
		}
		if utils.IsDigit(ch) {
			return l.readNumber(ch)
		}
	}

	return token.NewILLEGAL(string(ch))
}

func (l *Lexer) hasNext() bool {
	return l.endtoken != nil || l.readPosition < len(l.input)
}

func (l *Lexer) peek() rune {
	if l.readPosition == len(l.input) {
		return *l.endtoken
	}
	return rune(l.input[l.readPosition])
}

func (l *Lexer) next() rune {
	if l.readPosition == len(l.input) {
		ch := *l.endtoken
		l.endtoken = nil
		return ch
	}
	ch := rune(l.input[l.readPosition])
	l.readPosition += 1
	return ch
}

func (l *Lexer) readString(ch rune) *token.Token {
	str, state := string(ch), 0

	if ch == '"' {
		state = 1
	} else if ch == '\'' {
		state = 2
	} else {
		return token.NewILLEGAL(str)
	}

	for l.hasNext() {
		ch := l.next()
		switch state {
		case 1:
			if ch == '"' {
				return token.New(token.STRING, str+string(ch))
			}
		case 2:
			if ch == '\'' {
				return token.New(token.STRING, str+string(ch))
			}
		}
		str += string(ch)
	}

	return token.New(token.ILLEGAL, str)
}

func (l *Lexer) readIdentifier(ch rune) *token.Token {
	str := string(ch)
	for l.hasNext() {
		ch := l.peek()
		if !utils.IsLiteral(ch) {
			break
		}
		l.next()
		str += string(ch)
	}
	return token.ParseIndent(str)
}

func (l *Lexer) readNumber(ch rune) *token.Token {
	str, state := string(ch), 0
	for l.hasNext() {
		ch := l.peek()
		if !utils.IsDigit(ch) && ch != '.' {
			break
		}
		switch state {
		case 0:
			if ch == '.' {
				state = 1
			}
		case 1:
			if ch == '.' {
				return token.NewILLEGAL(str + string(ch))
			}
		}
		l.next()
		str += string(ch)
	}
	if state == 0 {
		return token.New(token.INT, str)
	}
	return token.New(token.FLOAT, str)
}

func (l *Lexer) skipWhiteSpace() {
	for utils.IsWhiteSpace(l.peek()) && l.hasNext() {
		l.next()
	}
}
