package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseRestParameterFollowedBy(value token.Token) *ast.RestParameter {
	p.consumeExpected(token.DOTDOTDOT)

	restParameter := &ast.RestParameter{
		Binder: p.parseBinder(),
	}

	p.shouldBe(value)

	return restParameter
}

func (p *Parser) parseObjectDestructureIdentifierParameter() *ast.ObjectPatternIdentifierParameter {
	var parameter ast.FunctionParameter
	var propertyName string

	// Rest parameter
	if p.is(token.DOTDOTDOT) {
		parameter = p.parseRestParameterFollowedBy(token.RIGHT_BRACE)
	} else {
		if p.is(token.IDENTIFIER) {
			identifier := p.parseIdentifier()

			parameter = &ast.IdentifierParameter{
				Name: identifier,
			}
			propertyName = identifier.Name

			if p.is(token.COLON) {
				p.next()

				// property rename
				parameter = &ast.IdentifierParameter{
					Name: p.parseIdentifier(),
				}
			}
		}

		// array destructure
		if p.is(token.LEFT_BRACKET) {
			parameter = p.parseArrayDestructureParameter()
		}

		// object destructure
		if p.is(token.LEFT_BRACE) {
			parameter = p.parseObjectDestructureParameter()
		}

		// Parsing default value
		if p.is(token.ASSIGN) {
			p.next()

			parameter.SetDefaultValue(p.parseAssignmentExpression())
		}
	}

	if !p.is(token.RIGHT_BRACE) {
		p.consumeExpected(token.COMMA)
	}

	return &ast.ObjectPatternIdentifierParameter{
		Parameter:    parameter,
		PropertyName: propertyName,
	}
}

func (p *Parser) parseObjectDestructureParameter() *ast.ObjectPatternParameter {
	p.next()

	if p.is(token.RIGHT_BRACE) {
		p.unexpectedToken()
	}

	param := &ast.ObjectPatternParameter{
		List:         make([]*ast.ObjectPatternIdentifierParameter, 0, 1),
		DefaultValue: nil,
	}

	for !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		param.List = append(
			param.List,
			p.parseObjectDestructureIdentifierParameter(),
		)
	}

	p.consumeExpected(token.RIGHT_BRACE)

	return param
}

func (p *Parser) parseArrayDestructureParameter() *ast.ArrayPatternParameter {
	p.next()

	if p.is(token.RIGHT_BRACKET) {
		p.unexpectedToken()
	}

	param := &ast.ArrayPatternParameter{
		List:         make([]ast.FunctionParameter, 0, 1),
		DefaultValue: nil,
	}

	for !p.is(token.RIGHT_BRACKET) && !p.is(token.EOF) {
		param.List = append(
			param.List,
			p.parseFunctionParameterEndingBy(token.RIGHT_BRACKET),
		)
	}

	p.consumeExpected(token.RIGHT_BRACKET)

	return param
}

func (p *Parser) parseFunctionParameterEndingBy(ending token.Token) ast.FunctionParameter {
	var parameter ast.FunctionParameter
	start := p.idx

	// Rest parameter
	if p.is(token.DOTDOTDOT) {
		parameter = p.parseRestParameterFollowedBy(ending)
	} else {
		if p.is(token.IDENTIFIER) {
			// simple argument binding
			parameter = &ast.IdentifierParameter{
				Name:         p.parseIdentifier(),
				DefaultValue: nil,
			}
		} else if p.is(token.LEFT_BRACKET) {
			// array destructure
			parameter = p.parseArrayDestructureParameter()
		} else if p.is(token.LEFT_BRACE) {
			// object destructure
			parameter = p.parseObjectDestructureParameter()
		} else {
			p.arrowFunctionMode = false
			p.unexpectedToken()
			p.next()
		}
	}

	if parameter == nil {
		p.error(start, "Unable to parse function argument")
		return nil
	}

	// optional
	if p.is(token.QUESTION_MARK) {
		p.next()
		parameter.SetOptional(true)
	}

	// type annotation
	if p.is(token.COLON) {
		parameter.SetTypeAnnotation(p.parseFlowTypeAnnotation())
	}

	// Parsing default value
	if p.is(token.ASSIGN) && parameter != nil {
		p.next()

		parameter.SetDefaultValue(p.parseAssignmentExpression())
	}

	if !p.is(ending) {
		p.consumePossible(token.COMMA)
	}

	return parameter
}

func (p *Parser) parseFunctionParameterList() *ast.FunctionParameters {
	opening := p.consumeExpected(token.LEFT_PARENTHESIS)
	var list []ast.FunctionParameter

	for !p.is(token.RIGHT_PARENTHESIS) && !p.is(token.EOF) {
		list = append(list, p.parseFunctionParameterEndingBy(token.RIGHT_PARENTHESIS))
	}
	closing := p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.FunctionParameters{
		Opening: opening,
		List:    list,
		Closing: closing,
	}
}

func (p *Parser) parseFunction(declaration bool, idx file.Idx, async bool) *ast.FunctionLiteral {
	p.consumeExpected(token.FUNCTION)

	generator := false
	if p.is(token.MULTIPLY) {
		p.consumeExpected(token.MULTIPLY)
		generator = true
	}

	node := &ast.FunctionLiteral{
		Start:     idx,
		Async:     async,
		Generator: generator,
	}

	var name *ast.Identifier
	if p.is(token.IDENTIFIER) {
		name = p.parseIdentifier()
		if declaration {
			p.scope.declare(&ast.FunctionDeclaration{
				Function: node,
			})
		}
	} else if declaration {
		// Use consumeExpected error handling
		p.consumeExpected(token.IDENTIFIER)
	}

	// type parameters
	if p.is(token.LESS) {
		node.TypeParameters = p.parseFlowTypeParameters()
	}

	node.Name = name
	node.Parameters = p.parseFunctionParameterList()

	if p.is(token.COLON) {
		node.ReturnType = p.parseFlowTypeAnnotation()
	}

	p.parseFunctionBlock(node)
	node.Source = p.slice(idx, node.EndAt())

	return node
}
