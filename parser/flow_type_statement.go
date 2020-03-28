package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeStatement() *ast.FlowTypeStatement {
	var typeParameters []*ast.FlowTypeParameter

	start := p.consumeExpected(token.TYPE_TYPE)
	name := p.parseFlowTypeIdentifier()

	if p.is(token.LESS) {
		typeParameters = p.parseFlowTypeParameters()
	}

	p.consumeExpected(token.ASSIGN)
	value := p.parseFlowType()

	p.implicitSemicolon = true
	p.optionalSemicolon()

	return &ast.FlowTypeStatement{
		Start:          start,
		Name:           name,
		Type:           value,
		TypeParameters: typeParameters,
	}
}