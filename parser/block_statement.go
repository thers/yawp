package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	node := &ast.BlockStatement{
		Loc: p.loc(),
	}

	p.consumeExpected(token.LEFT_BRACE)
	node.List = p.parseStatementList()
	node.Loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return node
}

func (p *Parser) parseBlockStatementOrObjectPatternBinding() ast.Statement {
	loc := p.loc()
	snapshot := p.snapshot()

	objectBinding, success := p.maybeParseObjectBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.ExpressionStatement{
			Expression: &ast.VariableBinding{
				Loc:         loc,
				Binder:      objectBinding,
				Initializer: p.parseAssignmentExpression(),
			},
		}
	}

	p.toSnapshot(snapshot)

	return p.parseBlockStatement()
}
