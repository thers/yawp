package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseTryStatement() ast.Statement {

	node := &ast.TryStatement{
		Try:  p.consumeExpected(token.TRY),
		Body: p.parseBlockStatement(),
	}

	if p.is(token.CATCH) {
		catch := p.loc
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
			Catch:     catch,
			Parameter: parameter,
			Body:      p.parseBlockStatement(),
		}
	}

	if p.is(token.FINALLY) {
		p.next()
		node.Finally = p.parseBlockStatement()
	}

	if node.Catch == nil && node.Finally == nil {
		p.error(node.Try, "Missing catch or finally after try")

		return nil
	}

	return node
}

