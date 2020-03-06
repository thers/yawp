package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseClassMethodBody() (*ast.BlockStatement, []ast.Declaration) {
	closeFunctionScope := p.openFunctionScope()
	defer closeFunctionScope()

	return p.parseBlockStatement(), p.scope.declarationList
}

func (p *Parser) parseClassBodyStatement() ast.Statement {
	start := p.idx
	async := false
	static := false
	private := false

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

	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		p.unexpectedToken()
		p.next()
		return nil
	}

	p.insertSemicolon = true

	// Could be set/get
	if identifier.Name == "set" || identifier.Name == "get" {
		kind := identifier.Name

		if p.is(token.IDENTIFIER) {
			field := p.parseIdentifier()
			parameterList := p.parseFunctionParameterList()

			node := &ast.FunctionLiteral{
				Start:      identifier.Start,
				Parameters: parameterList,
			}
			p.parseFunctionBlock(node)

			return &ast.ClassAccessorStatement{
				Start: identifier.Start,
				Field: field,
				Kind:  kind,
				Body:  node,
			}
		}
	}

	switch p.token {
	case token.ASSIGN:
		// Field can not be async, huh
		if async {
			p.unexpectedToken()
			return nil
		}

		p.next()
		stmt := &ast.ClassFieldStatement{
			Start:       start,
			Name:        identifier,
			Static:      static,
			Private:     private,
			Initializer: p.parseAssignmentExpression(),
		}
		p.semicolon()

		return stmt
	case token.LEFT_PARENTHESIS:
		// Method declaration
		method := &ast.ClassMethodStatement{
			Method:     start,
			Name:       identifier,
			Static:     static,
			Private:    private,
			Async:      async,
			Generator:  false,
			Parameters: p.parseFunctionParameterList(),
		}

		body, _ := p.parseClassMethodBody()
		source := p.slice(identifier.Start, body.EndAt())

		method.Body = body
		method.Source = source

		return method
	case token.IDENTIFIER, token.ASYNC, token.STATIC, token.HASH:
		p.unexpectedToken()
		p.next()
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
				p.error(0, "Class accessor definition cannot be decorated tho")
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

	if p.is(token.EXTENDS) {
		p.next()
		p.shouldBe(token.IDENTIFIER)

		exp.Extends = p.parseIdentifier()
	}

	p.scope.declare(&ast.ClassDeclaration{
		Name:    exp.Name,
		Extends: exp.Extends,
	})

	exp.Body = p.parseClassBody()

	return exp
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	return &ast.ClassStatement{
		Expression: p.parseClassExpression(),
	}
}
