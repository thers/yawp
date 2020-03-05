package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseReturnStatement() ast.Statement {
	idx := p.consumeExpected(token.RETURN)

	if !p.scope.inFunction {
		p.error(idx, "Illegal return statement")
		p.nextStatement()
		return &ast.BadStatement{From: idx, To: p.idx}
	}

	node := &ast.ReturnStatement{
		Return: idx,
	}

	if !p.implicitSemicolon && !p.is(token.SEMICOLON) && !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		node.Argument = p.parseExpression()
	}

	p.semicolon()

	return node
}
