package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseConditionalExpression() ast.Expression {
	left := p.parseNullishCoalescingExpression()

	if p.is(token.QUESTION_MARK) {
		p.next()
		consequent := p.parseAssignmentExpression()
		p.consumeExpected(token.COLON)
		return &ast.ConditionalExpression{
			Test:       left,
			Consequent: consequent,
			Alternate:  p.parseAssignmentExpression(),
		}
	}

	return left
}
