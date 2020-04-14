package parser

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) now() (token.Token, string, file.Loc) {
	return p.token, p.literal, p.loc
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

func (p *Parser) locOf(offset int) file.Loc {
	return file.Loc(p.base + offset)
}

func (p *Parser) is(value token.Token) bool {
	return p.token == value
}

func (p *Parser) allowNext(value token.Token) {
	if p.tokenIsKeyword {
		tkn, _ := token.IsKeyword(p.literal)

		if value == tkn {
			p.token = value
		}
	}
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

func (p *Parser) consumeExpected(value token.Token) file.Loc {
	idx := p.loc
	if p.token != value {
		p.unexpectedToken()
	}
	p.next()
	return idx
}

func (p *Parser) consumePossible(value token.Token) file.Loc {
	idx := p.loc

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

func lineCount(str string) (int, int) {
	line, last := 0, -1
	pair := false
	for index, chr := range str {
		switch chr {
		case '\r':
			line += 1
			last = index
			pair = true
			continue
		case '\n':
			if !pair {
				line += 1
			}
			last = index
		case '\u2028', '\u2029':
			line += 1
			last = index + 2
		}
		pair = false
	}
	return line, last
}

func (p *Parser) position(idx file.Loc) file.Position {
	position := file.Position{}
	offset := int(idx) - p.base
	str := p.src[:offset]
	position.Filename = p.file.Name()
	line, last := lineCount(str)
	position.Line = 1 + line
	if last >= 0 {
		position.Column = offset - last
	} else {
		position.Column = 1 + len(str)
	}

	return position
}
