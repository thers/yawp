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

func (p *Parser) parseClassBodyStatementList() []ast.Statement {
	stmt := make([]ast.Statement, 0)

	for !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		idx := p.idx
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

		switch p.token {
		case token.IDENTIFIER:
			identifier := p.parseIdentifier()

			switch p.token {
			case token.ASSIGN:
				// Property can not be async, huh
				if async {
					p.unexpectedToken()
					break
				}

				// Property init
				p.next()
				stmt = append(stmt, &ast.ClassPropertyStatement{
					Property:    idx,
					Name:        identifier,
					Static:      static,
					Private:     private,
					Initializer: p.parseAssignmentExpression(),
				})
				p.semicolon()
			case token.LEFT_PARENTHESIS:
				// Method declaration
				method := &ast.ClassMethodStatement{
					Method:     idx,
					Name:       identifier,
					Static:     static,
					Private:    private,
					Async:      async,
					Generator:  false,
					Parameters: p.parseFunctionParameterList(),
				}

				body, _ := p.parseClassMethodBody()
				source := p.slice(body.Idx0(), body.Idx1())

				method.Body = body
				method.Source = source

				stmt = append(stmt, method)
			default:
				p.unexpectedToken()
				p.next()
			}

		default:
			p.unexpectedToken()
			p.next()
		}
	}

	return stmt
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

func (p *Parser) parseClassStatement() ast.Statement {
	idx := p.consumeExpected(token.CLASS)

	statement := &ast.ClassStatement{
		Class:   idx,
		Name:    nil,
		Extends: nil,
		Body:    nil,
	}

	if p.is(token.IDENTIFIER) {
		statement.Name = p.parseIdentifier()
	}

	if p.is(token.EXTENDS) {
		p.next()
		p.shouldBe(token.IDENTIFIER)

		statement.Extends = p.parseIdentifier()
	}

	p.scope.declare(&ast.ClassDeclaration{
		Name:    statement.Name,
		Extends: statement.Extends,
	})

	statement.Body = p.parseClassBody()

	return statement
}
