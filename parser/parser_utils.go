package parser

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) now() (token.Token, string, file.Idx) {
	return p.token, p.literal, p.idx
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

func (p *Parser) idxOf(offset int) file.Idx {
	return file.Idx(offset)
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
	idx := p.idx

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
	return 0, 0

	//line, last := 0, -1
	//pair := false
	//for index, chr := range str {
	//	switch chr {
	//	case '\r':
	//		line += 1
	//		last = index
	//		pair = true
	//		continue
	//	case '\n':
	//		if !pair {
	//			line += 1
	//		}
	//		last = index
	//	case '\u2028', '\u2029':
	//		line += 1
	//		last = index + 2
	//	}
	//	pair = false
	//}
	//return line, last
}

func (p *Parser) loc() *file.Loc {
	return &file.Loc{
		From: p.idx,
		To:   file.Idx(p.nextChrOffset),
		Line: p.line,
		Col:  p.col,
	}
}

func (p *Parser) position(idx file.Idx) file.Pos {
	position := file.Pos{}

	offset := int(idx)
	str := p.src[:offset]

	line, last := lineCount(str)
	position.Line = 1 + line
	if last >= 0 {
		position.Column = offset - last
	} else {
		position.Column = 1 + len(str)
	}

	return position
}
