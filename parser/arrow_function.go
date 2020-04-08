package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) maybeParseArrowFunctionParameterList() (*ast.FunctionParameters, bool) {
	wasArrowMode := p.arrowFunctionMode
	p.arrowFunctionMode = true

	defer func() {
		p.arrowFunctionMode = wasArrowMode
	}()

	params := p.parseFunctionParameterList()

	return params, p.arrowFunctionMode
}

func (p *Parser) parseArrowFunctionBody() ast.Statement {
	// arrow function could not be generators
	closeFunctionScope := p.openFunctionScope(false)
	defer closeFunctionScope()

	if p.is(token.LEFT_BRACE) {
		return p.parseBlockStatement()
	}

	return &ast.ReturnStatement{
		Return:   p.idx,
		Argument: p.parseAssignmentExpression(),
	}
}

func (p *Parser) parseIdentifierOrSingleArgumentArrowFunction(async bool) ast.Expression {
	identifier := p.parseIdentifier()

	if p.is(token.ARROW) {
		// Parsing arrow function
		p.next()

		return &ast.ArrowFunctionExpression{
			Start: identifier.Start,
			Async: async,
			Parameters: []ast.FunctionParameter{
				&ast.IdentifierParameter{
					Name:         identifier,
					DefaultValue: nil,
				},
			},
			Body: p.parseArrowFunctionBody(),
		}
	}

	// Identifier
	if len(identifier.Name) > 1 {
		tkn, strict := token.IsKeyword(identifier.Name)
		if tkn == token.KEYWORD {
			if !strict {
				p.error(identifier.Start, "Unexpected reserved word")
			}
		}
	}
	return identifier
}

func (p *Parser) parseArrowFunctionOrSequenceExpression(async bool) ast.Expression {
	partialState := p.getPartialState()

	// First try to parse as arrow function parameters list
	parameters, success := p.maybeParseArrowFunctionParameterList()
	var returnType ast.FlowType

	success = success && len(partialState.errors) == len(p.errors)

	// may be also arrow fn return type or type assertion
	if success && p.is(token.COLON) {
		wasForbidden := p.forbidUnparenthesizedFunctionType
		p.forbidUnparenthesizedFunctionType = true
		returnType = p.parseFlowTypeAnnotation()
		p.forbidUnparenthesizedFunctionType = wasForbidden
	}

	// If no errors occurred while parsing parameters
	// And next token is => then it's an arrow function
	if success && p.is(token.ARROW) {
		p.next()
		return &ast.ArrowFunctionExpression{
			Start:      parameters.Opening,
			Async:      async,
			ReturnType: returnType,
			Parameters: parameters.List,
			Body:       p.parseArrowFunctionBody(),
		}
	}

	// It's a sequence expression or flow type assertion
	// restoring parser state like we didn't do shit
	p.restorePartialState(partialState)

	p.consumeExpected(token.LEFT_PARENTHESIS)
	wasAllowTypeAssertion := p.scope.allowTypeAssertion
	p.scope.allowTypeAssertion = true
	expression := p.parseExpression()
	p.scope.allowTypeAssertion = wasAllowTypeAssertion
	p.consumeExpected(token.RIGHT_PARENTHESIS)
	return expression
}

func (p *Parser) tryParseAsyncArrowFunction(idx file.Idx) ast.Expression {
	var typeParameters []*ast.FlowTypeParameter

	if p.is(token.LESS) {
		typeParameters = p.parseFlowTypeParameters()
	}

	if p.is(token.IDENTIFIER) {
		if typeParameters != nil {
			p.error(p.idx, "Parenthesis required around generic arrow function parameters")
			goto End
		}

		identifier := p.parseIdentifier()
		p.consumeExpected(token.ARROW)

		return &ast.ArrowFunctionExpression{
			Start:          identifier.Start,
			Async:          true,
			TypeParameters: typeParameters,
			Parameters: []ast.FunctionParameter{
				&ast.IdentifierParameter{
					Name:         identifier,
					DefaultValue: nil,
				},
			},
			Body: p.parseArrowFunctionBody(),
		}
	}

	if p.is(token.LEFT_PARENTHESIS) {
		var returnType ast.FlowType
		parameters := p.parseFunctionParameterList()

		if p.is(token.COLON) {
			p.forbidUnparenthesizedFunctionType = true
			returnType = p.parseFlowTypeAnnotation()
			p.forbidUnparenthesizedFunctionType = false
		}

		p.consumeExpected(token.ARROW)

		return &ast.ArrowFunctionExpression{
			Start:          parameters.Opening,
			Async:          true,
			TypeParameters: typeParameters,
			ReturnType:     returnType,
			Parameters:     parameters.List,
			Body:           p.parseArrowFunctionBody(),
		}
	}

End:
	return &ast.BadExpression{
		From: idx,
		To:   p.idx,
	}
}
