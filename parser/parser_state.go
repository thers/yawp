package parser

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type ParserStateSnapshot struct {
	idx               file.Idx
	token             token.Token
	errors            ErrorList
	offset            int
	chrOffset         int
	nextChr           rune
	literal           string
	parsedStr         string
	insertSemicolon   bool
	implicitSemicolon bool
}

func (p *Parser) captureState() *ParserStateSnapshot {
	return &ParserStateSnapshot{
		idx:               p.idx,
		token:             p.token,
		errors:            p.errors,
		offset:            p.nextChrOffset,
		chrOffset:         p.chrOffset,
		nextChr:           p.chr,
		literal:           p.literal,
		parsedStr:         p.parsedStr,
		insertSemicolon:   p.insertSemicolon,
		implicitSemicolon: p.implicitSemicolon,
	}
}

func (p *Parser) rewindStateTo(state *ParserStateSnapshot) {
	p.idx = state.idx
	p.token = state.token
	p.errors = state.errors
	p.nextChrOffset = state.offset
	p.chrOffset = state.chrOffset
	p.chr = state.nextChr
	p.literal = state.literal
	p.parsedStr = state.parsedStr
	p.insertSemicolon = state.insertSemicolon
	p.implicitSemicolon = state.implicitSemicolon
}

func (p *Parser) statesNeedsRewinding(ps *ParserStateSnapshot) bool {
	return len(ps.errors) != len(p.errors)
}

