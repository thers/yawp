package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseTryStatement() ast.Statement {
	loc := p.loc()
	p.consumeExpected(token.TRY)

	node := &ast.TryStatement{
		Loc:  loc,
		Body: p.parseBlockStatement(),
	}

	if p.is(token.CATCH) {
		catchLoc := p.loc()
		p.next()

		var parameter *ast.Identifier

		if p.is(token.LEFT_PARENTHESIS) {
			p.consumeExpected(token.LEFT_PARENTHESIS)

			if !p.is(token.IDENTIFIER) {
				p.unexpectedToken()

				return nil
			} else {
				parameter = p.parseIdentifier()
				p.consumeExpected(token.RIGHT_PARENTHESIS)
			}
		}

		node.Catch = &ast.CatchStatement{
			Loc:       catchLoc,
			Parameter: parameter,
			Body:      p.parseBlockStatement(),
		}
	}

	if p.is(token.FINALLY) {
		p.next()
		node.Finally = p.parseBlockStatement()
	}

	if node.Catch == nil && node.Finally == nil {
		p.error(node.Loc, "Missing catch or finally after try")

		return nil
	}

	return node
}
