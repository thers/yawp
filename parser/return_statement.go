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

	if !p.consumePossibleSemicolon() {
		node.Argument = p.parseExpression()

		p.semicolon()
	}

	return node
}
