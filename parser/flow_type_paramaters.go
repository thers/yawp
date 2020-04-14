package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeParameters() []*ast.FlowTypeParameter {
	closeTypeScope := p.openTypeScope()
	defer closeTypeScope()

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

func (p *Parser) isFlowTypeParametersStart() bool {
	return p.is(token.LESS)
}

func (p *Parser) isFlowTypeArgumentsStart() bool {
	return p.isAny(token.JSX_FRAGMENT_START, token.LESS)
}

func (p *Parser) isFlowTypeStart() bool {
	switch p.token {
	case
		token.IDENTIFIER,
		token.LESS,
		token.LEFT_BRACE,
		token.TYPE_EXACT_OBJECT_START,
		token.LEFT_PARENTHESIS,
		token.TYPE_ANY,
		token.TYPE_BOOLEAN,
		token.TYPE_MIXED,
		token.TYPE_NUMBER,
		token.TYPE_STRING,
		token.TYPE_TYPE,
		token.VOID,
		token.TYPEOF,
		token.BOOLEAN,
		token.STRING,
		token.NUMBER:
			return true
	default:
		return false
	}
}

func (p *Parser) parseFlowTypeArguments() []ast.FlowType {
	closeTypeScope := p.openTypeScope()
	defer closeTypeScope()

	// nice
	if p.is(token.JSX_FRAGMENT_START) {
		p.next()
		return nil
	}

	p.consumeExpected(token.LESS)
	args := make([]ast.FlowType, 0)


	for p.until(token.GREATER) {
		//if !p.isFlowTypeStart() {
		//	p.unexpectedToken()
		//	return nil
		//}

		argStart := p.idx
		arg := p.parseFlowType()

		if arg == nil {
			p.error(argStart, "Unexpected token")
			return nil
		}
		args = append(args, arg)
		p.consumePossible(token.COMMA)
	}

	p.consumeExpected(token.GREATER)
	return args
}

func (p *Parser) tryParseFlowTypeArguments() []ast.FlowType {
	defer func() {
		err := recover()

		if err != nil {
			return
		}
	}()

	return p.parseFlowTypeArguments()
}
