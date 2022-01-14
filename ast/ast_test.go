package ast_test

import (
	"testing"

	"monkey/ast"
	"monkey/token"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.New(token.LET, "let"),
				Name: &ast.Identifier{
					Token: token.New(token.IDENT, "myVar"),
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.New(token.IDENT, "anotherVar"),
					Value: "anotherVar",
				},
			},
		},
	}

	require.Equal(t, "let myVar = anotherVar;", program.String())

}
