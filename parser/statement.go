package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectNextToken(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectNextToken(token.ASSIGN) {
		return nil
	}

	p.getNextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.nextToken.Is(token.SEMICOLON) {
		p.getNextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.getNextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.nextToken.Is(token.SEMICOLON) {
		p.getNextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextToken.Is(token.SEMICOLON) {
		p.getNextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	p.getNextToken()

	// match } or EOF
	for p.currToken.IsNot(token.RBRACE) && p.currToken.IsNot(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.getNextToken()
	}

	return block
}
