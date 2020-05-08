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
		Loc: p.loc(),
	}

	p.consumeExpected(token.YIELD)

	if p.is(token.MULTIPLY) {
		p.consumeExpected(token.MULTIPLY)
		exp.Delegate = true
	}

	if p.newLineBeforeCurrentToken {
		return exp
	}

	switch p.token {
	case token.RIGHT_BRACE, token.RIGHT_BRACKET, token.RIGHT_PARENTHESIS:
		return exp
	}

	exp.Argument = p.parseAssignmentExpression()

	return exp
}

func (p *Parser) parseYieldStatement() *ast.YieldStatement {
	exp := &ast.YieldStatement{
		Expression: p.parseYieldExpression(),
	}

	p.semicolon()

	return exp
}
