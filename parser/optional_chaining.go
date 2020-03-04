package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseOptionalExpression(left ast.Expression) ast.Expression {
	exp := &ast.OptionalExpression{
		Left: left,
	}

	period := p.consumeExpected(token.OPTIONAL_CHAINING)

	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		if p.is(token.LEFT_BRACKET) {
			p.consumeExpected(token.LEFT_BRACKET)

			exp.ArrayMember = p.parseAssignmentExpression()

			p.consumeExpected(token.RIGHT_BRACKET)

			return exp
		} else {
			p.consumeExpected(token.IDENTIFIER)
			p.nextStatement()
			return &ast.BadExpression{From: period, To: p.idx}
		}
	}

	exp.Identifier = identifier

	p.next()

	return exp
}
