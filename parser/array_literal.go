package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	var value []ast.Expression

	start := p.consumeExpected(token.LEFT_BRACKET)
	for p.until(token.RIGHT_BRACKET) {
		// [,,,]
		if p.is(token.COMMA) {
			p.next()
			value = append(value, nil)
			continue
		}

		value = append(value, p.parseAssignmentExpression())

		p.consumePossible(token.COMMA)
	}
	end := p.consumeExpected(token.RIGHT_BRACKET)

	return &ast.ArrayLiteral{
		Start: start,
		End:   end,
		Value: value,
	}
}
