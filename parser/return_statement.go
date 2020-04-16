package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseReturnStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.RETURN)

	if !p.scope.inFunction {
		p.error(loc, "Illegal return statement")
	}

	node := &ast.ReturnStatement{
		Loc: loc,
	}

	if !p.implicitSemicolon && !p.is(token.SEMICOLON) && !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		node.Argument = p.parseExpression()
	}

	p.semicolon()

	return node
}
