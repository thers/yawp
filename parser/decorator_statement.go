package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseDecoratorsList() []ast.Expression {
	decorators := make([]ast.Expression, 0, 1)

	for p.is(token.AT) {
		p.consumeExpected(token.AT)

		decorators = append(decorators, p.parseAssignmentExpression())
	}

	return decorators
}

func (p *Parser) parseLegacyClassDecoratorStatement() *ast.LegacyDecoratorStatement {
	return &ast.LegacyDecoratorStatement{
		Decorators: p.parseDecoratorsList(),
		Subject:    p.parseClassStatement(),
	}
}
