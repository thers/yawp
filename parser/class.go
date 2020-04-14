package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseClassMethodBody(generator bool, async bool) (*ast.BlockStatement, []ast.Declaration) {
	closeFunctionScope := p.openFunctionScope(generator, async)
	defer closeFunctionScope()

	return p.parseBlockStatement(), p.scope.declarationList
}

func (p *Parser) parseClassFieldName() ast.ClassFieldName {
	var identifier ast.ClassFieldName

	if p.is(token.LEFT_BRACKET) {
		p.consumeExpected(token.LEFT_BRACKET)

		identifier = &ast.ComputedName{
			Expression: p.parseAssignmentExpression(),
		}

		p.consumeExpected(token.RIGHT_BRACKET)

		return identifier
	} else {
		identifier = p.parseIdentifierIncludingKeywords()
	}

	if identifier == nil {
		p.unexpectedToken()
		p.next()
	}

	return identifier
}

func (p *Parser) parseClassBodyStatement() ast.Statement {
	start := p.loc
	async := false
	static := false
	private := false
	generator := false

	if p.is(token.STATIC) {
		static = true
		p.next()
	}

	if p.is(token.ASYNC) {
		async = true
		p.next()
	}

	if p.is(token.HASH) {
		private = true
		p.next()
	}

	if p.is(token.MULTIPLY) {
		generator = true
		p.next()
	}

	identifier := p.parseClassFieldName()

	p.insertSemicolon = true

	// Could be set/get
	if accessor, ok := identifier.(*ast.Identifier); ok {
		if accessor.Name == "set" || accessor.Name == "get" {
			if async || private || generator {
				p.unexpectedToken()
				return nil
			}

			kind := accessor.Name

			if field := p.parseClassFieldName(); field != nil {
				node := &ast.FunctionLiteral{
					Start:      start,
					Parameters: p.parseFunctionParameterList(),
				}

				p.parseFunctionBlock(node)

				return &ast.ClassAccessorStatement{
					Start: start,
					Field: field,
					Kind:  kind,
					Body:  node,
				}
			}
		}
	}

	switch p.token {
	case token.ASSIGN, token.COLON:
		// Field can not be async or generator
		if async || generator {
			p.unexpectedToken()
			return nil
		}

		defer p.semicolon()

		stmt := &ast.ClassFieldStatement{
			Start:       start,
			Name:        identifier,
			Static:      static,
			Private:     private,
		}

		if p.is(token.COLON) {
			stmt.Type = p.parseFlowTypeAnnotation()

			if !p.is(token.ASSIGN) {
				return stmt
			}
		}

		p.consumeExpected(token.ASSIGN)

		stmt.Initializer = p.parseAssignmentExpression()

		return stmt
	case token.LESS, token.LEFT_PARENTHESIS:
		// Method declaration
		method := &ast.ClassMethodStatement{
			Method:    start,
			Name:      identifier,
			Static:    static,
			Private:   private,
			Async:     async,
			Generator: generator,
		}

		if p.is(token.LESS) {
			method.TypeParameters = p.parseFlowTypeParameters()
		}

		method.Parameters = p.parseFunctionParameterList()

		if p.is(token.COLON) {
			method.ReturnType = p.parseFlowTypeAnnotation()
		}

		body, _ := p.parseClassMethodBody(generator, async)

		method.Body = body

		return method
	case token.SEMICOLON:
		p.consumeExpected(token.SEMICOLON)
		return &ast.ClassFieldStatement{
			Start:       start,
			Name:        identifier,
			Static:      static,
			Private:     private,
			Initializer: nil,
		}
	default:
		if p.insertSemicolon {
			p.insertSemicolon = false
			return &ast.ClassFieldStatement{
				Start:       start,
				Name:        identifier,
				Static:      static,
				Private:     private,
				Initializer: nil,
			}
		} else {
			p.unexpectedToken()
			p.next()
		}
	}

	return nil
}

func (p *Parser) parseClassBodyStatementList() []ast.Statement {
	stmts := make([]ast.Statement, 0)

	for p.until(token.RIGHT_BRACE) {
		if p.is(token.AT) {
			decorators := p.parseDecoratorsList()
			statement := p.parseClassBodyStatement()

			var subject ast.LegacyDecoratorSubject

			switch refinedStatement := statement.(type) {
			case *ast.ClassFieldStatement:
				subject = refinedStatement
			case *ast.ClassMethodStatement:
				subject = refinedStatement
			default:
				p.error(0, "Class accessor definition can not be decorated")
				return nil
			}

			stmts = append(stmts, &ast.LegacyDecoratorStatement{
				Decorators: decorators,
				Subject:    subject,
			})
		} else {
			stmts = append(stmts, p.parseClassBodyStatement())
		}
	}

	return stmts
}

func (p *Parser) parseClassBody() ast.Statement {
	closeClassScope := p.openClassScope()
	defer closeClassScope()

	node := &ast.BlockStatement{}
	node.LeftBrace = p.consumeExpected(token.LEFT_BRACE)
	node.List = p.parseClassBodyStatementList()
	node.RightBrace = p.consumeExpected(token.RIGHT_BRACE)

	return node
}

func (p *Parser) parseClassExpression() *ast.ClassExpression {
	start := p.consumeExpected(token.CLASS)

	exp := &ast.ClassExpression{
		Start: start,
	}

	if p.is(token.IDENTIFIER) {
		exp.Name = p.parseIdentifier()
	}

	if p.isFlowTypeParametersStart() {
		exp.TypeParameters = p.parseFlowTypeParameters()
	}

	if p.is(token.EXTENDS) {
		p.next()

		exp.SuperClass = p.parseMemberExpression()

		if p.isFlowTypeArgumentsStart() {
			exp.SuperTypeArguments = p.parseFlowTypeArguments()
		}
	}

	p.scope.declare(&ast.ClassDeclaration{
		Name:    exp.Name,
		Extends: exp.SuperClass,
	})

	exp.Body = p.parseClassBody()

	return exp
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	return &ast.ClassStatement{
		Expression: p.parseClassExpression(),
	}
}
