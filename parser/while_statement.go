package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseWhileStatement() ast.Statement {
	p.consumeExpected(token.WHILE)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.WhileStatement{
		Test: p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)
	node.Body = p.parseIterationStatement()

	return node
}
