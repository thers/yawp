package parser

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) now() (token.Token, string, file.Idx) {
	return p.token, p.literal, p.tokenOffset
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

func (p *Parser) is(value token.Token) bool {
	return p.token == value
}

func (p *Parser) allowToken(value token.Token) {
	if p.tokenIsKeyword {
		tkn, _ := token.IsKeyword(p.literal)

		if value == tkn {
			p.token = value
		}
	}
}

func (p *Parser) isAllowed(value token.Token) bool {
	p.allowToken(value)

	return p.is(value)
}

func (p *Parser) isAny(values ...token.Token) bool {
	for _, value := range values {
		if value == p.token {
			return true
		}
	}

	return false
}

func (p *Parser) isIdentifierOrKeyword() bool {
	return p.is(token.IDENTIFIER) || p.tokenIsKeyword
}

func (p *Parser) until(value token.Token) bool {
	return p.token != value && p.token != token.EOF
}

func (p *Parser) consumeExpected(value token.Token) file.Idx {
	idx := file.Idx(p.chrOffset)
	if p.token != value {
		p.unexpectedToken()
	}
	p.next()
	return idx
}

func (p *Parser) consumePossible(value token.Token) file.Idx {
	idx := p.tokenOffset

	if p.token == value {
		p.next()
	}

	return idx
}

func (p *Parser) shouldBe(value token.Token) {
	if p.token != value {
		p.errorUnexpectedToken(p.token)
	}
}

func (p *Parser) loc() *file.Loc {
	return &file.Loc{
		From: p.tokenOffset,
		To:   file.Idx(p.nextChrOffset),
		Line: p.line,
		Col:  p.tokenCol,
	}
}
