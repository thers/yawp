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

func (p *Parser) parseFor(loc *file.Loc, initializer ast.Expression) *ast.ForStatement {

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

	var left []ast.Expression

	forIn := false
	forOf := false

	if !p.is(token.SEMICOLON) {
		allowIn := p.scope.allowIn
		p.scope.allowIn = false
		if p.isVariableStatementStart() {
			loc := p.loc()
			kind := p.token

			p.next()

			list := p.parseVariableDeclarationList(loc, kind)

			if len(list) == 1 {
				p.allowNext(token.OF)

				if p.is(token.IN) {
					p.consumeExpected(token.IN)
					forIn = true
					left = []ast.Expression{list[0]}
				} else if p.is(token.OF) {
					p.consumeExpected(token.OF)
					forOf = true
					left = []ast.Expression{list[0]}
				} else {
					left = list
				}
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
	} else {
		p.next()
		return p.parseFor(start, nil)
	}

	if left == nil {
		p.error(start, "Invalid for statement initializer")

		return nil
	}

	intoVar := left[0]

	if intoVar == nil {
		p.error(start, "Invalid for statement initializer")

		return nil
	}

	if forIn || forOf {
		switch intoVar.(type) {
		case *ast.Identifier:
		case *ast.DotExpression:
		case *ast.BracketExpression:
		case *ast.VariableExpression:
		case *ast.VariableBinding:
			// These are all acceptable
		default:
			p.error(start, "Invalid left-hand side in for-in")
		}

		if forIn {
			return p.parseForIn(start, intoVar)
		} else {
			return p.parseForOf(start, intoVar)
		}
	}

	p.consumeExpected(token.SEMICOLON)
	return p.parseFor(start, &ast.SequenceExpression{Sequence: left})
}
