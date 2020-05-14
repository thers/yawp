package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) maybeParseArrowFunctionParameterList() *ast.FunctionParameters {
	defer func() {
		_ = recover()
	}()

	return p.parseFunctionParameterList()
}

func (p *Parser) parseArrowFunctionBody(async bool) *ast.FunctionBody {
	// arrow function can not be generators
	closeFunctionScope := p.openFunctionScope(false, async)
	defer closeFunctionScope()

	p.consumeExpected(token.ARROW)

	if p.is(token.LEFT_BRACE) {
		return p.parseFunctionBody(false, async)
	}

	returnStatement := &ast.ReturnStatement{
		StmtNode: p.stmtNode(),
		Argument: p.parseAssignmentExpression(),
	}

	return &ast.FunctionBody{
		Node: (ast.Node)(returnStatement.StmtNode),
		List: ast.Statements{returnStatement},
	}
}

func (p *Parser) parseIdentifierOrSingleArgumentArrowFunction(async bool) ast.IExpr {
	identifier := p.parseIdentifier()

	if p.is(token.ARROW) {
		// Parsing arrow function
		return &ast.ArrowFunctionExpression{
			ExprNode: p.exprNodeAt(identifier.Loc),
			Async:    async,
			Parameters: []ast.FunctionParameter{
				&ast.IdentifierParameter{
					Id:           identifier,
					DefaultValue: nil,
				},
			},
			Body: p.parseArrowFunctionBody(async),
		}
	}

	// Identifier
	if len(identifier.Name) > 1 {
		tkn, strict := token.IsKeyword(identifier.Name)
		if tkn == token.KEYWORD {
			if !strict {
				p.error(identifier.Loc, "Unexpected reserved word")
			}
		}
	}
	return identifier
}

func (p *Parser) parseArrowFunctionOrSequenceExpression(async bool) ast.IExpr {
	snapshot := p.snapshot()

	// First try to parse as arrow function parameters list
	parameters := p.maybeParseArrowFunctionParameterList()
	var returnType ast.FlowType

	// may be also arrow fn return type or type assertion
	if parameters != nil && p.is(token.COLON) {
		wasForbidden := p.forbidUnparenthesizedFunctionType
		p.forbidUnparenthesizedFunctionType = true
		returnType = p.parseFlowTypeAnnotation()
		p.forbidUnparenthesizedFunctionType = wasForbidden
	}

	// If no errors occurred while parsing parameters
	// And next token is => then it's an arrow function
	if parameters != nil && p.is(token.ARROW) {
		return &ast.ArrowFunctionExpression{
			ExprNode:   p.exprNodeAt(parameters.Loc),
			Async:      async,
			ReturnType: returnType,
			Parameters: parameters.List,
			Body:       p.parseArrowFunctionBody(async),
		}
	}

	// It's a sequence expression or flow type assertion
	// restoring parser state like we didn't do shit
	p.toSnapshot(snapshot)

	p.consumeExpected(token.LEFT_PARENTHESIS)
	wasAllowTypeAssertion := p.scope.allowTypeAssertion
	p.scope.allowTypeAssertion = true
	expression := p.parseExpression()
	p.scope.allowTypeAssertion = wasAllowTypeAssertion
	p.consumeExpected(token.RIGHT_PARENTHESIS)
	return expression
}

func (p *Parser) tryParseAsyncArrowFunction(loc *file.Loc, st *ParserSnapshot) ast.IExpr {
	var typeParameters []*ast.FlowTypeParameter

	if p.is(token.LESS) {
		typeParameters = p.parseFlowTypeParameters()
	}

	if p.is(token.IDENTIFIER) {
		if typeParameters != nil {
			p.error(p.loc(), "Parenthesis required around generic arrow function parameters")

			return nil
		}

		identifier := p.parseIdentifier()

		return &ast.ArrowFunctionExpression{
			ExprNode:       p.exprNodeAt(loc),
			Async:          true,
			TypeParameters: typeParameters,
			Parameters: []ast.FunctionParameter{
				&ast.IdentifierParameter{
					Id:           identifier,
					DefaultValue: nil,
				},
			},
			Body: p.parseArrowFunctionBody(true),
		}
	}

	if p.is(token.LEFT_PARENTHESIS) {
		var returnType ast.FlowType

		sloc := p.loc()
		parameters := p.parseFunctionParameterList()

		if p.is(token.COLON) {
			p.forbidUnparenthesizedFunctionType = true
			returnType = p.parseFlowTypeAnnotation()
			p.forbidUnparenthesizedFunctionType = false
		}

		// may be async(bla) fn call
		if !p.is(token.ARROW) {
			p.toSnapshot(st)
			p.token = token.IDENTIFIER

			return p.parseAssignmentExpression()
		}

		return &ast.ArrowFunctionExpression{
			ExprNode:       p.exprNodeAt(sloc),
			Async:          true,
			TypeParameters: typeParameters,
			ReturnType:     returnType,
			Parameters:     parameters.List,
			Body:           p.parseArrowFunctionBody(true),
		}
	}

	p.unexpectedToken()

	return nil
}

func (p *Parser) parseParametrizedArrowFunction() *ast.ArrowFunctionExpression {
	loc := p.loc()
	typeParameters := p.parseFlowTypeParameters()

	if p.is(token.LEFT_PARENTHESIS) {
		var returnType ast.FlowType
		parameters := p.parseFunctionParameterList()

		if p.is(token.COLON) {
			p.forbidUnparenthesizedFunctionType = true
			returnType = p.parseFlowTypeAnnotation()
			p.forbidUnparenthesizedFunctionType = false
		}

		return &ast.ArrowFunctionExpression{
			ExprNode:       p.exprNodeAt(loc),
			Async:          false,
			TypeParameters: typeParameters,
			ReturnType:     returnType,
			Parameters:     parameters.List,
			Body:           p.parseArrowFunctionBody(false),
		}
	}

	p.unexpectedToken()
	p.next()

	return nil
}
