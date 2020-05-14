package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseNewExpression() ast.IExpr {
	loc := p.loc()
	p.consumeExpected(token.NEW)

	if p.is(token.PERIOD) {
		periodLoc := p.loc()
		p.next()

		if !p.is(token.IDENTIFIER) || p.literal != "target" {
			p.error(periodLoc, "Unexpected period")
			return nil
		}

		loc.End(p.consumeExpected(token.IDENTIFIER))

		return &ast.NewTargetExpression{
			ExprNode: p.exprNodeAt(loc),
		}
	}

	callee := p.parseLeftHandSideExpression()
	node := &ast.NewExpression{
		ExprNode: p.exprNodeAt(loc),
		Callee:   callee,
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
