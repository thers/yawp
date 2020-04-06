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

		parameterNameStart := p.idx
		parameter.Name = p.parseFlowTypeIdentifierIncludingKeywords()

		if parameter.Name == nil {
			p.error(parameterNameStart, "Type parameter name is required")
			p.next()
			continue
		}

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

func (p *Parser) isFlowTypeArgumentsStart() bool {
	return p.isAny(token.JSX_FRAGMENT_START, token.LESS)
}

func (p *Parser) parseFlowTypeArguments() []ast.FlowType {
	// nice
	if p.is(token.JSX_FRAGMENT_START) {
		p.next()
		return nil
	}

	p.consumeExpected(token.LESS)
	args := make([]ast.FlowType, 0)


	for p.until(token.GREATER) {
		args = append(args, p.parseFlowType())
		p.consumePossible(token.COMMA)
	}

	p.consumeExpected(token.GREATER)
	return args
}
