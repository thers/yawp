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
		catch := p.idx
		p.next()

		var parameter *ast.Identifier

		if p.is(token.LEFT_PARENTHESIS) {
			p.consumeExpected(token.LEFT_PARENTHESIS)

			if !p.is(token.IDENTIFIER) {
				p.consumeExpected(token.IDENTIFIER)
				p.nextStatement()
				return &ast.BadStatement{From: catch, To: p.idx}
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
		return &ast.BadStatement{From: node.Try, To: node.Body.EndAt()}
	}

	return node
}

