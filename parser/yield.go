package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseYieldExpression() *ast.YieldExpression {
	if !p.scope.allowYield {
		p.error(p.loc(), "yield can not be used outside of generator function")
		p.next()
		return nil
	}

	exp := &ast.YieldExpression{
		ExprNode: p.exprNode(),
	}

	p.consumeExpected(token.YIELD)

	if !p.consumePossibleSemicolon() {
		if p.is(token.MULTIPLY) {
			p.consumeExpected(token.MULTIPLY)
			exp.Delegate = true
		}

		exp.Argument = p.parseAssignmentExpression()
	}

	return exp
}
