package parser

import (
	"yawp/ids"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseSourceElement() ast.Statement {
	return p.parseStatement()
}

func (p *Parser) parseSourceElements() []ast.Statement {
	body := make([]ast.Statement, 0)

	defer func() {
		err := recover()

		if err == nil {
			return
		}

		switch terr := err.(type) {
		case error:
			p.err = terr
		case *SyntaxError:
			p.err = terr.Error()
		default:
			panic(terr)
		}
	}()

	for !p.is(token.EOF) {
		stmt := p.parseStatement()

		if _, ok := stmt.(*ast.EmptyStatement); !ok {
			body = append(body, stmt)
		}
	}

	return body
}

func (p *Parser) parseProgram() *ast.Module {
	p.openScope()
	defer p.closeScope()
	return &ast.Module{
		Body: p.parseSourceElements(),
		File: p.file,
		Ids:  ids.NewIds(),
	}
}

func (p *Parser) parseEmptyStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.SEMICOLON)
	return &ast.EmptyStatement{Loc: loc}
}

func (p *Parser) parseStatementList() (list []ast.Statement) {
	for !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		list = append(list, p.parseStatement())
	}

	return
}

func (p *Parser) parseStatement() ast.Statement {
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
		return p.parseArrayBindingStatementOrArrayLiteral()
	case token.LEFT_BRACE:
		return p.parseBlockStatementOrObjectPatternBinding()
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
		return p.parseWithStatement()
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
	case token.YIELD:
		return p.parseYieldStatement()
	case token.TYPE_TYPE, token.TYPE_OPAQUE:
		if p.scope.inModuleRoot() {
			return p.parseFlowTypeStatement()
		}
	case token.INTERFACE:
		return p.parseFlowInterfaceStatement()
	}

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

func (p *Parser) parseFunctionBlock(node *ast.FunctionLiteral) {
	closeFunctionScope := p.openFunctionScope(node.Generator, node.Async)
	defer closeFunctionScope()

	node.Body = p.parseBlockStatement()
}

func (p *Parser) parseDebuggerStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.DEBUGGER)

	node := &ast.DebuggerStatement{
		Loc: loc,
	}

	p.semicolon()

	return node
}

func (p *Parser) parseThrowStatement() ast.Statement {
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
		Loc:      loc,
		Argument: p.parseExpression(),
	}

	p.semicolon()

	return node
}

func (p *Parser) parseSwitchStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.SWITCH)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.SwitchStatement{
		Loc:          loc,
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

func (p *Parser) parseWithStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.WITH)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.WithStatement{
		Loc:    loc,
		Object: p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	node.Body = p.parseStatement()

	return node
}

func (p *Parser) parseCaseStatement() *ast.CaseStatement {
	node := &ast.CaseStatement{
		Loc: p.loc(),
	}
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

func (p *Parser) parseDoWhileStatement() ast.Statement {
	inIteration := p.scope.inIteration
	p.scope.inIteration = true
	defer func() {
		p.scope.inIteration = inIteration
	}()

	loc := p.loc()
	p.consumeExpected(token.DO)

	node := &ast.DoWhileStatement{
		Loc: loc,
	}

	if p.is(token.LEFT_BRACE) {
		node.Body = p.parseBlockStatement()
	} else {
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

func (p *Parser) parseIfStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.IF)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.IfStatement{
		Loc:  loc,
		Test: p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	if p.is(token.LEFT_BRACE) {
		node.Consequent = p.parseBlockStatement()
	} else {
		node.Consequent = p.parseStatement()
	}

	if p.is(token.ELSE) {
		p.next()
		node.Alternate = p.parseStatement()
	}

	return node
}

func (p *Parser) parseBreakStatement() ast.Statement {
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
			Loc:   loc,
			Token: token.BREAK,
		}
	}

	if p.is(token.IDENTIFIER) {
		identifier := p.parseIdentifier()

		if !p.scope.hasLabel(identifier.Name) {
			p.error(loc, "Undefined label '%s'", identifier.Name)

			return nil
		}

		p.semicolon()

		return &ast.BranchStatement{
			Loc:   loc,
			Token: token.BREAK,
			Label: identifier,
		}
	}

	p.consumeExpected(token.IDENTIFIER)

illegal:
	p.error(loc, "Illegal break statement")

	return nil
}

func (p *Parser) parseContinueStatement() ast.Statement {
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
			Loc:   loc,
			Token: token.CONTINUE,
		}
	}

	if p.is(token.IDENTIFIER) {
		identifier := p.parseIdentifier()
		if !p.scope.hasLabel(identifier.Name) {
			p.error(loc, "Undefined label '%s'", identifier.Name)

			return nil
		}
		if !p.scope.inIteration {
			goto illegal
		}
		p.semicolon()
		return &ast.BranchStatement{
			Loc:   loc,
			Token: token.CONTINUE,
			Label: identifier,
		}
	}

	p.consumeExpected(token.IDENTIFIER)

illegal:
	p.error(loc, "Illegal continue statement")

	return nil
}
