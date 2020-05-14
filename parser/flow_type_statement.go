package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeStatement() *ast.FlowTypeStatement {
	opaque := p.is(token.TYPE_OPAQUE)

	if opaque {
		p.next()
	}

	var typeParameters []*ast.FlowTypeParameter

	loc := p.loc()

	p.allowToken(token.TYPE_TYPE)
	p.consumeExpected(token.TYPE_TYPE)

	name := p.parseFlowTypeIdentifier()

	if p.is(token.LESS) {
		typeParameters = p.parseFlowTypeParameters()
	}

	p.consumeExpected(token.ASSIGN)
	value := p.parseFlowType()

	p.implicitSemicolon = true
	p.optionalSemicolon()

	return &ast.FlowTypeStatement{
		StmtNode:       p.stmtNodeAt(loc),
		Name:           name,
		Opaque:         opaque,
		Type:           value,
		TypeParameters: typeParameters,
	}
}
