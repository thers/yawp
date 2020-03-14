package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseYieldExpression() *ast.YieldExpression {
	exp := &ast.YieldExpression{
		Start: p.idx,
	}

	wasImplicitSemicolon := p.implicitSemicolon
	p.consumeExpected(token.YIELD)

	if !wasImplicitSemicolon && p.implicitSemicolon {
		p.error(exp.Start, "No line terminator allowed after yield keyword")
		return nil
	}

	exp.Expression = p.parseAssignmentExpression()

	return exp
}

func (p *Parser) parseYieldStatement() *ast.YieldStatement {
	return &ast.YieldStatement{
		Expression: p.parseYieldExpression(),
	}
}
