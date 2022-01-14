package parser

import (
	"fmt"
	"strconv"

	"monkey/ast"
	"monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -Xor!X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func (p *Parser) currPrecedence() int {
	if v, ok := precedences[p.currToken.Type]; ok {
		return v
	}
	return LOWEST
}

func (p *Parser) nextPrecedence() int {
	if v, ok := precedences[p.nextToken.Type]; ok {
		return v
	}
	return LOWEST
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.errors = append(p.errors,
			fmt.Sprintf("no prefix parse function for %s found", p.currToken.Type))
		return nil
	}

	leftExp := prefix()

	for p.nextToken.IsNot(token.SEMICOLON) && precedence < p.nextPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExp
		}

		p.getNextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors,
			fmt.Sprintf("could not parse %q as integer", p.currToken.Literal))
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.getNextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.getNextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currToken.Is(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.getNextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectNextToken(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}

	// match (
	if !p.expectNextToken(token.LPAREN) {
		return nil
	}

	// next token
	p.getNextToken()

	// parse condition
	expression.Condition = p.parseExpression(LOWEST)

	// match )
	if !p.expectNextToken(token.RPAREN) {
		return nil
	}

	// match {
	if !p.expectNextToken(token.LBRACE) {
		return nil
	}

	// parse consequence
	expression.Consequence = p.parseBlockStatement()

	if p.nextToken.Is(token.ELSE) {
		p.getNextToken()

		if !p.expectNextToken(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectNextToken(token.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectNextToken(token.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.nextToken.Is(token.RPAREN) {
		p.getNextToken()
		return identifiers
	}

	p.getNextToken()

	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)

	for p.nextToken.Is(token.COMMA) {
		p.getNextToken()
		p.getNextToken()
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectNextToken(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.nextToken.Is(token.RPAREN) {
		p.getNextToken()
		return args
	}

	// first args
	p.getNextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.nextToken.Is(token.COMMA) {
		p.getNextToken()
		p.getNextToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectNextToken(token.RPAREN) {
		return nil
	}

	return args
}
