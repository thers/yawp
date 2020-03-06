package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
	next := p.parseUnaryExpression
	left := next()

	for p.is(token.MULTIPLY) || p.is(token.SLASH) ||
		p.is(token.REMAINDER) || p.is(token.EXPONENTIATION) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseAdditiveExpression() ast.Expression {
	next := p.parseMultiplicativeExpression
	left := next()

	for p.is(token.PLUS) || p.is(token.MINUS) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}
