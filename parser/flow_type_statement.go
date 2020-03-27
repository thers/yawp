package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeStatement() *ast.FlowTypeStatement {
	start := p.consumeExpected(token.TYPE_TYPE)

	name := p.parseFlowTypeIdentifier()
	p.consumeExpected(token.ASSIGN)
	value := p.parseFlowType()

	p.semicolon()

	return &ast.FlowTypeStatement{
		Start: start,
		Name:  name,
		Type:  value,
	}
}
