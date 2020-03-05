package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseNullishCoalescingExpression() ast.Expression {
	left := p.parseLogicalOrExpression()

	if p.is(token.NULLISH_COALESCING) {
		p.consumeExpected(token.NULLISH_COALESCING)

		return &ast.NullishCoalescingExpression{
			Test:       left,
			Consequent: p.parseAssignmentExpression(),
		}
	}

	return left
}
