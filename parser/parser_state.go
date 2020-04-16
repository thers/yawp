package parser

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type ParserSnapshot struct {
	idx               file.Idx
	token             token.Token
	offset            int
	chrOffset         int
	nextChr           rune
	literal           string
	parsedStr         string
	insertSemicolon   bool
	implicitSemicolon bool
}

func (p *Parser) snapshot() *ParserSnapshot {
	return &ParserSnapshot{
		idx:               p.idx,
		token:             p.token,
		offset:            p.nextChrOffset,
		chrOffset:         p.chrOffset,
		nextChr:           p.chr,
		literal:           p.literal,
		parsedStr:         p.parsedStr,
		insertSemicolon:   p.insertSemicolon,
		implicitSemicolon: p.implicitSemicolon,
	}
}

func (p *Parser) toSnapshot(state *ParserSnapshot) {
	p.idx = state.idx
	p.token = state.token
	p.nextChrOffset = state.offset
	p.chrOffset = state.chrOffset
	p.chr = state.nextChr
	p.literal = state.literal
	p.parsedStr = state.parsedStr
	p.insertSemicolon = state.insertSemicolon
	p.implicitSemicolon = state.implicitSemicolon
}
