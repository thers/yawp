package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseNewExpression() ast.Expression {
	start := p.consumeExpected(token.NEW)

	if p.is(token.PERIOD) {
		p.consumeExpected(token.PERIOD)

		if !p.is(token.IDENTIFIER) && p.literal != "target" {
			p.unexpectedToken()
			p.next()
			return nil
		}

		end := p.consumeExpected(token.IDENTIFIER)

		return &ast.NewTargetExpression{
			Start: start,
			End:   end,
		}
	}

	callee := p.parseLeftHandSideExpression()
	node := &ast.NewExpression{
		Start:  start,
		Callee: callee,
	}

	if p.isFlowTypeArgumentsStart() {
		node.TypeArguments = p.parseFlowTypeArguments()
	}

	if p.is(token.LEFT_PARENTHESIS) {
		argumentList, idx0, idx1 := p.parseArgumentList()
		node.ArgumentList = argumentList
		node.LeftParenthesis = idx0
		node.RightParenthesis = idx1
	}

	return node
}
