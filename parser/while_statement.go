package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseWhileStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.WHILE)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.WhileStatement{
		StmtNode: p.stmtNodeAt(loc),
		Test:     p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	p.useSymbolsScope(ast.SSTBlock)
	defer p.restoreSymbolsScope()

	node.Body = p.parseIterationStatement()

	return node
}
