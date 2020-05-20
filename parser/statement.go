package parser

import (
	"yawp/ids"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseSourceElement() ast.IStmt {
	return p.parseStatement()
}

func (p *Parser) parseSourceElements() []ast.IStmt {
	body := make([]ast.IStmt, 0)

	for !p.is(token.EOF) {
		stmt := p.parseStatement()

		if _, ok := stmt.(*ast.EmptyStatement); !ok {
			body = append(body, stmt)
		}
	}

	return body
}

func (p *Parser) parseModule() *ast.Module {
	defer p.recover()

	p.openScope()
	defer p.closeScope()
	return &ast.Module{
		Symbols: p.symbolsScope,
		Body:    p.parseSourceElements(),
		File:    p.file,
		Ids:     ids.NewIds(),
	}
}

func (p *Parser) parseEmptyStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.SEMICOLON)

	return &ast.EmptyStatement{
		StmtNode: p.stmtNodeAt(loc),
	}
}

func (p *Parser) parseStatementList() (list []ast.IStmt) {
	for !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		list = append(list, p.parseStatement())
	}

	return
}

func (p *Parser) parseStatement() ast.IStmt {
	if p.is(token.EOF) {
		p.unexpectedToken()
		return nil
	}

	if p.scope.inModuleRoot() {
		p.allowToken(token.TYPE_TYPE)
	}

	switch p.token {
	case token.SEMICOLON:
		return p.parseEmptyStatement()
	case token.LEFT_BRACKET:
		return p.parseArrayPatternAssignmentOrLiteral()
	case token.LEFT_BRACE:
		return p.parseBlockStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.DO:
		return p.parseDoWhileStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.DEBUGGER:
		return p.parseDebuggerStatement()
	case token.WITH:
		p.unexpectedToken()
		return nil
		// with is banned in strict mode
		//return p.parseWithStatement()
	case token.VAR, token.CONST, token.LET:
		return p.parseVariableStatement()
	case token.FUNCTION:
		return p.parseFunction(true, p.loc(), false)
	case token.SWITCH:
		return p.parseSwitchStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.THROW:
		return p.parseThrowStatement()
	case token.TRY:
		return p.parseTryStatement()
	case token.CLASS:
		return p.parseClassStatement()
	case token.IMPORT:
		return p.parseImportDeclaration()
	case token.EXPORT:
		return p.parseExportDeclaration()
	case token.AT:
		return p.parseLegacyClassDecoratorStatement()
	case token.TYPE_TYPE, token.TYPE_OPAQUE:
		if p.scope.inModuleRoot() {
			return p.parseFlowTypeStatement()
		}
	case token.INTERFACE:
		return p.parseFlowInterfaceStatement()
	}

	loc := p.loc()
	expression := p.parseExpression()

	if identifier, isIdentifier := expression.(*ast.Identifier); isIdentifier && p.is(token.COLON) {
		// LabelledStatement
		colon := p.tokenOffset
		p.next() // :
		label := identifier.Name
		for _, value := range p.scope.labels {
			if label == value {
				p.error(identifier.GetLoc(), "Label '%s' already exists", label)
			}
		}
		p.scope.labels = append(p.scope.labels, label) // Push the label
		statement := p.parseStatement()
		p.scope.labels = p.scope.labels[:len(p.scope.labels)-1] // Pop the label
		return &ast.LabelledStatement{
			Label:     identifier,
			Colon:     colon,
			Statement: statement,
		}
	}

	p.optionalSemicolon()

	return &ast.ExpressionStatement{
		StmtNode:   p.stmtNodeAt(loc),
		Expression: expression,
	}
}

func (p *Parser) parseParameterList() (list []string) {
	for !p.is(token.EOF) {
		if !p.is(token.IDENTIFIER) {
			p.consumeExpected(token.IDENTIFIER)
		}
		list = append(list, p.literal)
		p.next()
		if !p.is(token.EOF) {
			p.consumeExpected(token.COMMA)
		}
	}
	return
}

func (p *Parser) parseFunctionBody(generator, async bool) *ast.FunctionBody {
	closeFunctionScope := p.openFunctionScope(generator, async)
	defer closeFunctionScope()

	body := &ast.FunctionBody{
		Node: p.node(),
	}

	p.consumeExpected(token.LEFT_BRACE)
	body.List = p.parseStatementList()
	body.Loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return body
}

func (p *Parser) parseFunctionNodeBody(node *ast.FunctionLiteral) {
	node.Body = p.parseFunctionBody(node.Generator, node.Async)
}

func (p *Parser) parseDebuggerStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.DEBUGGER)

	node := &ast.DebuggerStatement{
		StmtNode: p.stmtNodeAt(loc),
	}

	p.semicolon()

	return node
}

func (p *Parser) parseThrowStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.THROW)

	if p.implicitSemicolon {
		if p.chr == -1 { // Hackish
			p.error(loc, "Unexpected end of input")
		} else {
			p.error(loc, "Illegal newline after throw")
		}

		return nil
	}

	node := &ast.ThrowStatement{
		StmtNode: p.stmtNodeAt(loc),
		Argument: p.parseExpression(),
	}

	p.semicolon()

	return node
}

func (p *Parser) parseSwitchStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.SWITCH)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.SwitchStatement{
		StmtNode:     p.stmtNodeAt(loc),
		Discriminant: p.parseExpression(),
		Default:      -1,
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	p.consumeExpected(token.LEFT_BRACE)

	inSwitch := p.scope.inSwitch
	p.scope.inSwitch = true
	defer func() {
		p.scope.inSwitch = inSwitch
	}()

	for index := 0; !p.is(token.EOF); index++ {
		if p.is(token.RIGHT_BRACE) {
			p.next()
			break
		}

		clause := p.parseCaseStatement()
		if clause.Test == nil {
			if node.Default != -1 {
				p.error(clause.Loc, "Already saw a default in switch")
			}
			node.Default = index
		}
		node.Loc.Add(clause.Loc)
		node.Body = append(node.Body, clause)
	}

	return node
}

func (p *Parser) parseWithStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.WITH)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.WithStatement{
		StmtNode: p.stmtNodeAt(loc),
		Object:   p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	node.Body = p.parseStatement()

	return node
}

func (p *Parser) parseCaseStatement() *ast.CaseStatement {
	node := &ast.CaseStatement{
		StmtNode: p.stmtNode(),
	}

	p.useSymbolsScope(ast.SSTBlock)
	defer p.restoreSymbolsScope()

	if p.is(token.DEFAULT) {
		p.next()
	} else {
		p.consumeExpected(token.CASE)
		node.Test = p.parseExpression()
	}

	p.consumeExpected(token.COLON)

	for {
		if p.is(token.EOF) ||
			p.is(token.RIGHT_BRACE) ||
			p.is(token.CASE) ||
			p.is(token.DEFAULT) {
			break
		}

		consequent := p.parseStatement()
		node.Loc.Add(consequent.GetLoc())
		node.Consequent = append(node.Consequent, consequent)

	}

	return node
}

func (p *Parser) parseDoWhileStatement() ast.IStmt {
	inIteration := p.scope.inIteration
	p.scope.inIteration = true
	defer func() {
		p.scope.inIteration = inIteration
	}()

	loc := p.loc()
	p.consumeExpected(token.DO)

	node := &ast.DoWhileStatement{
		StmtNode: p.stmtNodeAt(loc),
	}

	if p.is(token.LEFT_BRACE) {
		node.Body = p.parseBlockStatement()
	} else {
		p.useSymbolsScope(ast.SSTBlock)
		defer p.restoreSymbolsScope()

		node.Body = p.parseStatement()
	}

	p.consumeExpected(token.WHILE)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node.Test = p.parseExpression()
	p.consumeExpected(token.RIGHT_PARENTHESIS)
	if p.is(token.SEMICOLON) {
		p.next()
	}

	return node
}

func (p *Parser) parseIfStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.IF)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.IfStatement{
		StmtNode: p.stmtNodeAt(loc),
		Test:     p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	if p.is(token.LEFT_BRACE) {
		node.Consequent = p.parseBlockStatement()
	} else {
		p.useSymbolsScope(ast.SSTBlock)
		defer p.restoreSymbolsScope()

		node.Consequent = p.parseStatement()
	}

	if p.is(token.ELSE) {
		p.next()

		p.useSymbolsScope(ast.SSTBlock)
		defer p.restoreSymbolsScope()

		node.Alternate = p.parseStatement()
	}

	return node
}

func (p *Parser) parseBreakStatement() ast.IStmt {
	loc := p.loc()

	p.consumeExpected(token.BREAK)
	semicolon := p.implicitSemicolon

	if p.is(token.SEMICOLON) {
		semicolon = true
		p.next()
	}

	if semicolon || p.is(token.RIGHT_BRACE) {
		p.implicitSemicolon = false
		if !p.scope.inIteration && !p.scope.inSwitch {
			goto illegal
		}
		return &ast.BranchStatement{
			StmtNode: p.stmtNodeAt(loc),
			Token:    token.BREAK,
		}
	}

	if p.is(token.IDENTIFIER) {
		identifier := p.parseIdentifier()

		p.semicolon()

		return &ast.BranchStatement{
			StmtNode: p.stmtNodeAt(loc),
			Token:    token.BREAK,
			Label:    identifier,
		}
	}

	p.consumeExpected(token.IDENTIFIER)

illegal:
	p.error(loc, "Illegal break statement")

	return nil
}

func (p *Parser) parseContinueStatement() ast.IStmt {
	loc := p.loc()
	p.consumeExpected(token.CONTINUE)

	semicolon := p.implicitSemicolon

	if p.is(token.SEMICOLON) {
		semicolon = true
		p.next()
	}

	if semicolon || p.is(token.RIGHT_BRACE) {
		p.implicitSemicolon = false
		if !p.scope.inIteration {
			goto illegal
		}
		return &ast.BranchStatement{
			StmtNode: p.stmtNodeAt(loc),
			Token:    token.CONTINUE,
		}
	}

	if p.is(token.IDENTIFIER) {
		identifier := p.parseIdentifier()

		if !p.scope.inIteration {
			goto illegal
		}
		p.semicolon()
		return &ast.BranchStatement{
			StmtNode: p.stmtNodeAt(loc),
			Token:    token.CONTINUE,
			Label:    identifier,
		}
	}

	p.consumeExpected(token.IDENTIFIER)

illegal:
	p.error(loc, "Illegal continue statement")

	return nil
}
