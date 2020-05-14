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

func (p *Parser) parseFunctionParameterEndingBy(ending token.Token) ast.FunctionParameter {
	var parameter ast.FunctionParameter
	loc := p.loc()

	switch p.token {
	case token.DOTDOTDOT:
		parameter = p.parseRestParameterFollowedBy(ending)
	case token.IDENTIFIER:
		parameter = &ast.IdentifierParameter{
			Id:           p.parseIdentifier(),
			DefaultValue: nil,
		}
	case token.LEFT_BRACKET, token.LEFT_BRACE:
		parameter = &ast.PatternParameter{
			Binder: p.parseBinder(),
		}
	default:
		p.unexpectedToken()

		return nil
	}

	if parameter == nil {
		p.error(loc, "Unable to parse function argument")
		return nil
	}

	// optional
	if p.is(token.QUESTION_MARK) {
		p.next()
		parameter.SetFlowTypeOptional(true)
	}

	// type annotation
	if p.is(token.COLON) {
		parameter.SetTypeAnnotation(p.parseFlowTypeAnnotation())
	}

	// Parsing default value
	if p.is(token.ASSIGN) {
		p.next()

		parameter.SetDefaultValue(p.parseAssignmentExpression())
	}

	if !p.is(ending) {
		p.consumePossible(token.COMMA)
	}

	return parameter
}

func (p *Parser) parseFunctionParameterList() *ast.FunctionParameters {
	loc := p.loc()

	wasAllowYield := p.scope.allowYield
	p.scope.allowYield = false
	defer func() {
		p.scope.allowYield = wasAllowYield
	}()

	p.consumeExpected(token.LEFT_PARENTHESIS)
	var list []ast.FunctionParameter

	for !p.is(token.RIGHT_PARENTHESIS) && !p.is(token.EOF) {
		list = append(list, p.parseFunctionParameterEndingBy(token.RIGHT_PARENTHESIS))
	}

	loc.End(p.consumeExpected(token.RIGHT_PARENTHESIS))

	return &ast.FunctionParameters{
		Node: p.nodeAt(loc),
		List: list,
	}
}

func (p *Parser) parseFunction(declaration bool, loc *file.Loc, async bool) *ast.FunctionLiteral {
	wasAllowYield := p.scope.allowYield
	wasAllowAwait := p.scope.allowAwait
	p.scope.allowYield = false
	p.scope.allowAwait = false

	p.consumeExpected(token.FUNCTION)

	generator := false
	if p.is(token.MULTIPLY) {
		p.consumeExpected(token.MULTIPLY)
		generator = true
	}

	node := &ast.FunctionLiteral{
		Node:      p.nodeAt(loc),
		Async:     async,
		Generator: generator,
	}

	var name *ast.Identifier
	if p.is(token.IDENTIFIER) {
		name = p.parseIdentifier()
	} else if declaration {
		// Use consumeExpected error handling
		p.consumeExpected(token.IDENTIFIER)
	}

	p.scope.allowYield = wasAllowYield
	p.scope.allowAwait = wasAllowAwait

	// type parameters
	if p.is(token.LESS) {
		node.TypeParameters = p.parseFlowTypeParameters()
	}

	node.Id = name
	node.Parameters = p.parseFunctionParameterList()

	if p.is(token.COLON) {
		node.ReturnType = p.parseFlowTypeAnnotation()
	}

	p.parseFunctionNodeBody(node)

	return node
}
