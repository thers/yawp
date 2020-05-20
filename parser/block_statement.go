package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	p.useSymbolsScope(ast.SSTBlock)
	defer p.restoreSymbolsScope()

	node := &ast.BlockStatement{
		StmtNode: p.stmtNode(),
	}

	p.consumeExpected(token.LEFT_BRACE)
	node.List = p.parseStatementList()
	node.Loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return node
}
