package parser

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type ParserSnapshot struct {
	wasNewLine        bool
	tokenOffset       file.Idx
	token             token.Token
	line              int
	col               int
	nextChrOffset     int
	chrOffset         int
	chr               rune
	literal           string
	parsedStr         string
	insertSemicolon   bool
	implicitSemicolon bool
}

func (p *Parser) snapshot() *ParserSnapshot {
	return &ParserSnapshot{
		wasNewLine:        p.advanceLine,
		line:              p.line,
		col:               p.chrCol,
		tokenOffset:       p.tokenOffset,
		token:             p.token,
		nextChrOffset:     p.nextChrOffset,
		chrOffset:         p.chrOffset,
		chr:               p.chr,
		literal:           p.literal,
		parsedStr:         p.parsedSrc,
		insertSemicolon:   p.insertSemicolon,
		implicitSemicolon: p.implicitSemicolon,
	}
}

func (p *Parser) toSnapshot(state *ParserSnapshot) {
	p.advanceLine = state.wasNewLine
	p.line = state.line
	p.chrCol = state.col
	p.tokenOffset = state.tokenOffset
	p.token = state.token
	p.nextChrOffset = state.nextChrOffset
	p.chrOffset = state.chrOffset
	p.chr = state.chr
	p.literal = state.literal
	p.parsedSrc = state.parsedStr
	p.insertSemicolon = state.insertSemicolon
	p.implicitSemicolon = state.implicitSemicolon
}
