package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeParameters() []*ast.FlowTypeParameter {
	parameters := make([]*ast.FlowTypeParameter, 0)
	defaultValueRequired := false

	p.consumeExpected(token.LESS)
	for p.until(token.GREATER) {
		parameter := &ast.FlowTypeParameter{
			Start:         p.idx,
			Name:          nil,
			Covariant:     false,
			Contravariant: false,
			Boundary:      nil,
			DefaultValue:  nil,
		}

		if p.isAny(token.MINUS, token.PLUS) {
			parameter.Covariant = p.is(token.PLUS)
			parameter.Contravariant = p.is(token.MINUS)

			p.next()
		}

		parameter.Name = p.parseFlowTypeIdentifierIncludingKeywords()

		if p.is(token.COLON) {
			p.next()

			parameter.Boundary = p.parseFlowType()
		}

		if p.is(token.ASSIGN) {
			p.next()
			defaultValueRequired = true

			parameter.DefaultValue = p.parseFlowType()
		} else if defaultValueRequired {
			p.error(p.idx, "Default value required")
		}

		p.consumePossible(token.COMMA)

		parameters = append(parameters, parameter)

	}
	p.consumeExpected(token.GREATER)

	return parameters
}
