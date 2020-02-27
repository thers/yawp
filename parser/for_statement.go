package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseIterationStatement() ast.Statement {
	inIteration := p.scope.inIteration
	p.scope.inIteration = true
	defer func() {
		p.scope.inIteration = inIteration
	}()
	return p.parseStatement()
}

func (p *Parser) parseForIn(idx file.Idx, into ast.Expression) *ast.ForInStatement {

	// Already have consumed "<into> in"

	source := p.parseExpression()
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForInStatement{
		For:    idx,
		Into:   into,
		Source: source,
		Body:   p.parseIterationStatement(),
	}
}

func (p *Parser) parseFor(idx file.Idx, initializer ast.Expression) *ast.ForStatement {

	// Already have consumed "<initializer> ;"

	var test, update ast.Expression

	if !p.is(token.SEMICOLON) {
		test = p.parseExpression()
	}
	p.consumeExpected(token.SEMICOLON)

	if !p.is(token.RIGHT_PARENTHESIS) {
		update = p.parseExpression()
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForStatement{
		For:         idx,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        p.parseIterationStatement(),
	}
}

func (p *Parser) parseForOrForInStatement() ast.Statement {
	idx := p.consumeExpected(token.FOR)
	p.consumeExpected(token.LEFT_PARENTHESIS)

	var left []ast.Expression

	forIn := false
	if !p.is(token.SEMICOLON) {

		allowIn := p.scope.allowIn
		p.scope.allowIn = false
		if p.isVariableStatementStart() {
			kind := p.token
			var_ := p.idx
			p.next()
			list := p.parseVariableDeclarationList(var_, kind)
			if len(list) == 1 && p.is(token.IN) {
				p.next() // in
				forIn = true
				left = []ast.Expression{list[0]} // There is only one declaration
			} else {
				left = list
			}
		} else {
			left = append(left, p.parseExpression())
			if p.is(token.IN) {
				p.next()
				forIn = true
			}
		}
		p.scope.allowIn = allowIn
	}

	if forIn {
		switch left[0].(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression, *ast.VariableExpression:
			// These are all acceptable
		default:
			p.error(idx, "Invalid left-hand side in for-in")
			p.nextStatement()
			return &ast.BadStatement{From: idx, To: p.idx}
		}
		return p.parseForIn(idx, left[0])
	}

	p.consumeExpected(token.SEMICOLON)
	return p.parseFor(idx, &ast.SequenceExpression{Sequence: left})
}
