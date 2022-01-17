package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currToken.Is(token.TRUE)}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for p.nextToken.IsNot(token.RBRACE) {
		p.getNextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectNextToken(token.COLON) {
			return nil
		}

		p.getNextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if p.nextToken.IsNot(token.RBRACE) && !p.expectNextToken(token.COMMA) {
			return nil
		}
	}

	if !p.expectNextToken(token.RBRACE) {
		return nil
	}

	return hash
}
