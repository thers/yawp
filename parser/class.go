package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseClassBodyStatement() ast.IStmt {
	loc := p.loc()
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

	identifier := p.parseObjectPropertyName()

	p.insertSemicolon = true

	// Could be set/get
	if accessor, ok := identifier.(*ast.Identifier); ok {
		if accessor.Name == "set" || accessor.Name == "get" {
			if async || private || generator {
				p.unexpectedToken()
				return nil
			}

			kind := accessor.Name

			if field := p.parseObjectPropertyName(); field != nil {
				node := &ast.FunctionLiteral{
					Node:       p.nodeAt(loc),
					Parameters: p.parseFunctionParameterList(),
				}

				p.parseFunctionNodeBody(node)

				return &ast.ClassAccessorStatement{
					StmtNode: p.stmtNodeAt(loc),
					Field:    field,
					Kind:     kind,
					Body:     node,
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
			StmtNode: p.stmtNodeAt(loc),
			Name:     identifier,
			Static:   static,
			Private:  private,
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
			StmtNode:  p.stmtNodeAt(loc),
			Name:      identifier,
			Static:    static,
			Private:   private,
			Async:     async,
			Generator: generator,
		}

		if p.is(token.LESS) {
			method.TypeParameters = p.parseFlowTypeParameters()
		}

		p.useSymbolsScope(ast.SSTFunction)
		defer p.restoreSymbolsScope()

		method.Parameters = p.parseFunctionParameterList()

		if p.is(token.COLON) {
			method.ReturnType = p.parseFlowTypeAnnotation()
		}

		body := p.parseFunctionBody(generator, async)

		method.Body = body

		return method
	case token.SEMICOLON:
		p.consumeExpected(token.SEMICOLON)
		return &ast.ClassFieldStatement{
			StmtNode:    p.stmtNodeAt(loc),
			Name:        identifier,
			Static:      static,
			Private:     private,
			Initializer: nil,
		}
	default:
		if p.insertSemicolon {
			p.insertSemicolon = false
			return &ast.ClassFieldStatement{
				StmtNode:    p.stmtNodeAt(loc),
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

func (p *Parser) parseClassBodyStatementList() []ast.IStmt {
	stmts := make([]ast.IStmt, 0)

	for p.until(token.RIGHT_BRACE) {
		if p.is(token.SEMICOLON) {
			p.next()
			continue
		}

		if p.is(token.AT) {
			loc := p.loc()
			decorators := p.parseDecoratorsList()
			statement := p.parseClassBodyStatement()

			var subject ast.IStmt

			switch refinedStatement := statement.(type) {
			case *ast.ClassFieldStatement:
				subject = refinedStatement
			case *ast.ClassMethodStatement:
				subject = refinedStatement
			default:
				p.error(loc, "Class accessor definition can not be decorated")
				return nil
			}

			stmts = append(stmts, &ast.LegacyDecoratorStatement{
				Decorators: decorators,
				Subject:    subject,
			})
		} else {
			stmts = append(stmts, p.parseClassBodyStatement())
		}

		p.consumePossible(token.SEMICOLON)
	}

	return stmts
}

func (p *Parser) parseClassBody() ast.IStmt {
	closeClassScope := p.openClassScope()
	defer closeClassScope()

	p.useSymbolsScope(ast.SSTClass)

	node := &ast.BlockStatement{
		StmtNode: p.stmtNode(),
	}

	p.consumeExpected(token.LEFT_BRACE)
	node.List = p.parseClassBodyStatementList()
	node.Loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return node
}

func (p *Parser) parseClassExpression() *ast.ClassExpression {
	loc := p.loc()
	p.consumeExpected(token.CLASS)

	exp := &ast.ClassExpression{
		ExprNode: p.exprNodeAt(loc),
	}

	if p.is(token.IDENTIFIER) {
		exp.Name = p.symbol(p.parseIdentifier(), ast.SDeclaration, ast.SRClass)
		exp.Name.Symbol.RefType = ast.SRClass
	}

	if p.isFlowTypeParametersStart() {
		exp.TypeParameters = p.parseFlowTypeParameters()
	}

	if p.is(token.EXTENDS) {
		p.next()

		exp.SuperClass = p.parseLeftHandSideExpression()

		if p.isFlowTypeArgumentsStart() {
			exp.SuperTypeArguments = p.parseFlowTypeArguments()
		}
	}

	exp.Body = p.parseClassBody()

	return exp
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	return &ast.ClassStatement{
		StmtNode:   p.stmtNode(),
		Expression: p.parseClassExpression(),
	}
}
