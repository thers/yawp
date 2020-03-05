package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseLeftHandSideExpression() ast.Expression {

	var left ast.Expression
	if p.is(token.NEW) {
		left = p.parseNewExpression()
	} else {
		left = p.parsePrimaryExpression()
	}

	for {
		if p.is(token.PERIOD) {
			left = p.parseDotMember(left)
		} else if p.is(token.OPTIONAL_CHAINING) {
			left = p.parseOptionalExpression(left)
		} else if p.is(token.LEFT_BRACKET) {
			left = p.parseBracketMember(left)
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parseLeftHandSideExpressionAllowCall() ast.Expression {

	allowIn := p.scope.allowIn
	p.scope.allowIn = true
	defer func() {
		p.scope.allowIn = allowIn
	}()

	var left ast.Expression
	if p.is(token.NEW) {
		left = p.parseNewExpression()
	} else {
		left = p.parsePrimaryExpression()
	}

	for {
		if p.is(token.PERIOD) {
			left = p.parseDotMember(left)
		} else if p.is(token.OPTIONAL_CHAINING) {
			left = p.parseOptionalExpression(left)
		} else if p.is(token.LEFT_BRACKET) {
			left = p.parseBracketMember(left)
		} else if p.is(token.LEFT_PARENTHESIS) {
			left = p.parseCallExpression(left)
		} else {
			break
		}
	}

	return left
}