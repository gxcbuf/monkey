package token

import "fmt"

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 123456
	STRING = "STRING" // "hello", 'world'
	FLOAT  = "FLOAT"  // 0.12, 1.232

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

type Token struct {
	Type    TokenType
	Literal string
}

func New(ttype TokenType, literal string) *Token {
	return &Token{
		Type:    ttype,
		Literal: literal,
	}
}

func NewILLEGAL(literal string) *Token {
	return New(ILLEGAL, literal)
}

func ParseIndent(ident string) *Token {
	if ttype, ok := keywords[ident]; ok {
		return New(ttype, ident)
	}
	return New(IDENT, ident)
}

func (t *Token) Is(ttype TokenType) bool {
	return t.Type == ttype
}

func (t *Token) IsNot(ttype TokenType) bool {
	return t.Type != ttype
}

func (t *Token) IsEOF() bool {
	return t.Is(EOF)
}

func (t *Token) String() string {
	return fmt.Sprintf("type %s, value %s", t.Type, t.Literal)
}
