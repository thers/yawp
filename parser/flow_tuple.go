package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTupleType() *ast.FlowTupleType {
	start := p.consumeExpected(token.LEFT_BRACKET)

	elements := make([]ast.FlowType, 0)

	for p.until(token.RIGHT_BRACKET) {
		elements = append(elements, p.parseFlowType())

		if !p.is(token.RIGHT_BRACKET) {
			p.consumeExpected(token.COMMA)
		}
	}

	return &ast.FlowTupleType{
		Start:    start,
		End:      p.consumeExpected(token.RIGHT_BRACKET),
		Elements: elements,
	}
}
