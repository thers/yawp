package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	node := &ast.BlockStatement{
		StmtNode: p.stmtNode(),
	}

	p.consumeExpected(token.LEFT_BRACE)
	node.List = p.parseStatementList()
	node.Loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return node
}

func (p *Parser) parseBlockStatementOrObjectPatternBinding() ast.IStmt {
	loc := p.loc()
	snapshot := p.snapshot()

	objectBinding, success := p.maybeParseObjectBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.ExpressionStatement{
			StmtNode: p.stmtNode(),
			Expression: &ast.VariableBinding{
				ExprNode:    p.exprNodeAt(loc),
				Binder:      objectBinding,
				Initializer: p.parseAssignmentExpression(),
			},
		}
	}

	p.toSnapshot(snapshot)

	return p.parseBlockStatement()
}
