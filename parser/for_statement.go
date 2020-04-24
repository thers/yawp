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

func (p *Parser) parseForIn(loc *file.Loc, into ast.Expression) *ast.ForInStatement {
	// Already have consumed "<into> in"

	source := p.parseExpression()
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForInStatement{
		Loc:    loc,
		Into:   into,
		Source: source,
		Body:   p.parseIterationStatement(),
	}
}

func (p *Parser) parseForOf(loc *file.Loc, into ast.Expression) *ast.ForOfStatement {
	// Already have consumed "<into> of"

	source := p.parseExpression()
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForOfStatement{
		Loc:      loc,
		Binder:   into,
		Iterator: source,
		Body:     p.parseIterationStatement(),
	}
}

func (p *Parser) parseFor(loc *file.Loc, initializer *ast.VariableStatement) *ast.ForStatement {

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
		Loc:         loc,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        p.parseIterationStatement(),
	}
}

func (p *Parser) parseForOrForInStatement() ast.Statement {
	start := p.loc()
	p.consumeExpected(token.FOR)
	p.consumeExpected(token.LEFT_PARENTHESIS)

	// for (;
	if p.is(token.SEMICOLON) {
		p.next()
		return p.parseFor(start, nil)
	}

	// for(const/var/let
	if p.isVariableStatementStart() {
		loc := p.loc()
		kind := p.token

		p.next()

		list := p.parseVariableDeclarationList(kind)
		p.allowToken(token.OF)

		switch p.token {
		case token.IN:
			p.next()
			return p.parseForIn(start, list[0])
		case token.OF:
			p.next()
			return p.parseForOf(start, list[0])
		case token.SEMICOLON:
			p.next()
			return p.parseFor(start, &ast.VariableStatement{
				Loc:  loc,
				Kind: kind,
				List: list,
			})
		}
	}

	// for (anything
	intoExpression := p.parseExpression()

	if p.is(token.IN) {
		switch intoExpression.(type) {
		case *ast.Identifier:
		case *ast.DotExpression:
		case *ast.BracketExpression:
			// These are all acceptable
		default:
			p.error(start, "Invalid left-hand side in for-in")
		}

		p.next()
		return p.parseForIn(start, intoExpression)
	}

	p.error(start, "Invalid for statement initializer")
	return nil
}
