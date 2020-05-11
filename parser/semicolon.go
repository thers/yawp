package parser

import "yawp/parser/token"

func (p *Parser) consumePossibleSemicolon() bool {
	if p.is(token.SEMICOLON) {
		p.next()
		return true
	}

	if p.implicitSemicolon {
		return true
	}

	if p.isAny(token.RIGHT_PARENTHESIS, token.RIGHT_BRACE, token.RIGHT_BRACKET) {
		return true
	}

	return false
}

func (p *Parser) optionalSemicolon() {
	if p.is(token.SEMICOLON) {
		p.next()
		return
	}

	if p.implicitSemicolon {
		p.implicitSemicolon = false
		return
	}

	if !p.is(token.EOF) && !p.is(token.RIGHT_BRACE) {
		p.consumeExpected(token.SEMICOLON)
	}
}

func (p *Parser) semicolon() {
	if !p.is(token.RIGHT_PARENTHESIS) && !p.is(token.RIGHT_BRACE) {
		if p.implicitSemicolon {
			p.implicitSemicolon = false
			return
		}

		p.consumeExpected(token.SEMICOLON)
	}
}
