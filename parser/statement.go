package parser

import (
	"encoding/base64"
	"github.com/go-sourcemap/sourcemap"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseEmptyStatement() ast.Statement {
	idx := p.consumeExpected(token.SEMICOLON)
	return &ast.EmptyStatement{Semicolon: idx}
}

func (p *Parser) parseStatementList() (list []ast.Statement) {
	for !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		list = append(list, p.parseStatement())
	}

	return
}

func (p *Parser) parseStatement() ast.Statement {

	if p.is(token.EOF) {
		p.errorUnexpectedToken(p.token)
		return &ast.BadStatement{From: p.idx, To: p.idx + 1}
	}

	switch p.token {
	case token.SEMICOLON:
		return p.parseEmptyStatement()
	case token.LEFT_BRACKET:
		return p.parseArrayBindingStatement()
	case token.LEFT_BRACE:
		return p.parseBlockStatementOrObjectPatternBinding()
	case token.IF:
		return p.parseIfStatement()
	case token.DO:
		return p.parseDoWhileStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForOrForInStatement()
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
		return p.parseFunction(true, p.idx, false)
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
	case token.TYPE_TYPE:
		return p.parseFlowTypeStatement()
	case token.INTERFACE:
		return p.parseFlowInterfaceStatement()
	}

	expression := p.parseExpression()

	if identifier, isIdentifier := expression.(*ast.Identifier); isIdentifier && p.is(token.COLON) {
		// LabelledStatement
		colon := p.idx
		p.next() // :
		label := identifier.Name
		for _, value := range p.scope.labels {
			if label == value {
				p.error(identifier.StartAt(), "Label '%s' already exists", label)
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
	closeFunctionScope := p.openFunctionScope(node.Generator)
	defer closeFunctionScope()

	node.Body = p.parseBlockStatement()
	node.DeclarationList = p.scope.declarationList
}

func (p *Parser) parseDebuggerStatement() ast.Statement {
	idx := p.consumeExpected(token.DEBUGGER)

	node := &ast.DebuggerStatement{
		Debugger: idx,
	}

	p.semicolon()

	return node
}

func (p *Parser) parseThrowStatement() ast.Statement {
	idx := p.consumeExpected(token.THROW)

	if p.implicitSemicolon {
		if p.chr == -1 { // Hackish
			p.error(idx, "Unexpected end of input")
		} else {
			p.error(idx, "Illegal newline after throw")
		}
		p.nextStatement()
		return &ast.BadStatement{From: idx, To: p.idx}
	}

	node := &ast.ThrowStatement{
		Argument: p.parseExpression(),
	}

	p.semicolon()

	return node
}

func (p *Parser) parseSwitchStatement() ast.Statement {
	p.consumeExpected(token.SWITCH)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.SwitchStatement{
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
				p.error(clause.Case, "Already saw a default in switch")
			}
			node.Default = index
		}
		node.Body = append(node.Body, clause)
	}

	return node
}

func (p *Parser) parseWithStatement() ast.Statement {
	p.consumeExpected(token.WITH)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.WithStatement{
		Object: p.parseExpression(),
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	node.Body = p.parseStatement()

	return node
}

func (p *Parser) parseCaseStatement() *ast.CaseStatement {

	node := &ast.CaseStatement{
		Case: p.idx,
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
		node.Consequent = append(node.Consequent, p.parseStatement())

	}

	return node
}

func (p *Parser) parseDoWhileStatement() ast.Statement {
	inIteration := p.scope.inIteration
	p.scope.inIteration = true
	defer func() {
		p.scope.inIteration = inIteration
	}()

	p.consumeExpected(token.DO)
	node := &ast.DoWhileStatement{}
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
	p.consumeExpected(token.IF)
	p.consumeExpected(token.LEFT_PARENTHESIS)
	node := &ast.IfStatement{
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

func (p *Parser) parseSourceElement() ast.Statement {
	return p.parseStatement()
}

func (p *Parser) parseSourceElements() []ast.Statement {
	body := []ast.Statement(nil)

	for {
		if !p.is(token.STRING) {
			break
		}

		body = append(body, p.parseSourceElement())
	}

	for !p.is(token.EOF) {
		body = append(body, p.parseSourceElement())
	}

	return body
}

func (p *Parser) parseProgram() *ast.Program {
	p.openScope()
	defer p.closeScope()
	return &ast.Program{
		Body:            p.parseSourceElements(),
		DeclarationList: p.scope.declarationList,
		File:            p.file,
		SourceMap:       p.parseSourceMap(),
	}
}

func (p *Parser) parseSourceMap() *sourcemap.Consumer {
	lastLine := p.str[strings.LastIndexByte(p.str, '\n')+1:]
	if strings.HasPrefix(lastLine, "//# sourceMappingURL") {
		urlIndex := strings.Index(lastLine, "=")
		urlStr := lastLine[urlIndex+1:]

		var data []byte
		if strings.HasPrefix(urlStr, "data:application/json") {
			b64Index := strings.Index(urlStr, ",")
			b64 := urlStr[b64Index+1:]
			if d, err := base64.StdEncoding.DecodeString(b64); err == nil {
				data = d
			}
		} else {
			if smUrl, err := url.Parse(urlStr); err == nil {
				if smUrl.Scheme == "" || smUrl.Scheme == "file" {
					if f, err := os.Open(smUrl.Path); err == nil {
						if d, err := ioutil.ReadAll(f); err == nil {
							data = d
						}
					}
				} else {
					// Not implemented - compile error?
					return nil
				}
			}
		}

		if data == nil {
			return nil
		}

		if sm, err := sourcemap.Parse(p.file.Name(), data); err == nil {
			return sm
		}
	}
	return nil
}

func (p *Parser) parseBreakStatement() ast.Statement {
	idx := p.consumeExpected(token.BREAK)
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
			Idx:   idx,
			Token: token.BREAK,
		}
	}

	if p.is(token.IDENTIFIER) {
		identifier := p.parseIdentifier()
		if !p.scope.hasLabel(identifier.Name) {
			p.error(idx, "Undefined label '%s'", identifier.Name)
			return &ast.BadStatement{From: idx, To: identifier.EndAt()}
		}
		p.semicolon()
		return &ast.BranchStatement{
			Idx:   idx,
			Token: token.BREAK,
			Label: identifier,
		}
	}

	p.consumeExpected(token.IDENTIFIER)

illegal:
	p.error(idx, "Illegal break statement")
	p.nextStatement()
	return &ast.BadStatement{From: idx, To: p.idx}
}

func (p *Parser) parseContinueStatement() ast.Statement {
	idx := p.consumeExpected(token.CONTINUE)
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
			Idx:   idx,
			Token: token.CONTINUE,
		}
	}

	if p.is(token.IDENTIFIER) {
		identifier := p.parseIdentifier()
		if !p.scope.hasLabel(identifier.Name) {
			p.error(idx, "Undefined label '%s'", identifier.Name)
			return &ast.BadStatement{From: idx, To: identifier.EndAt()}
		}
		if !p.scope.inIteration {
			goto illegal
		}
		p.semicolon()
		return &ast.BranchStatement{
			Idx:   idx,
			Token: token.CONTINUE,
			Label: identifier,
		}
	}

	p.consumeExpected(token.IDENTIFIER)

illegal:
	p.error(idx, "Illegal continue statement")
	p.nextStatement()
	return &ast.BadStatement{From: idx, To: p.idx}
}

// Find the next statement after an error (recover)
func (p *Parser) nextStatement() {
	for {
		switch p.token {
		case token.BREAK, token.CONTINUE,
			token.FOR, token.IF, token.RETURN, token.SWITCH,
			token.VAR, token.DO, token.TRY, token.WITH,
			token.WHILE, token.THROW, token.CATCH, token.FINALLY:
			// Return only if parser made some progress since last
			// sync or if it has not reached 10 next calls without
			// progress. Otherwise consume at least one token to
			// avoid an endless parser loop
			if p.idx == p.recover.idx && p.recover.count < 10 {
				p.recover.count++
				return
			}
			if p.idx > p.recover.idx {
				p.recover.idx = p.idx
				p.recover.count = 0
				return
			}
			// Reaching here indicates a parser bug, likely an
			// incorrect token list in this function, but it only
			// leads to skipping of possibly correct code if a
			// previous error is present, and thus is preferred
			// over a non-terminating parse.
		case token.EOF:
			return
		}
		p.next()
	}
}
