package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseNewExpression() ast.Expression {
	loc := p.loc()
	p.consumeExpected(token.NEW)

	if p.is(token.PERIOD) {
		p.consumeExpected(token.PERIOD)

		if !p.is(token.IDENTIFIER) && p.literal != "target" {
			p.unexpectedToken()
			p.next()
			return nil
		}

		loc.End(p.consumeExpected(token.IDENTIFIER))

		return &ast.NewTargetExpression{
			Loc: loc,
		}
	}

	callee := p.parseLeftHandSideExpression()
	node := &ast.NewExpression{
		Loc:    loc,
		Callee: callee,
	}

	if p.isFlowTypeArgumentsStart() {
		node.TypeArguments = p.parseFlowTypeArguments()
	}

	if p.is(token.LEFT_PARENTHESIS) {
		argumentList, _, _ := p.parseArgumentList()
		node.ArgumentList = argumentList
	}

	return node
}
