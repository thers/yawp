package parser

import "yawp/parser/token"

func (p *Parser) switchAssignment2(tkn0, tkn1 token.Token) token.Token {
	if p.chr == '=' {
		p.read()
		return tkn1
	}
	return tkn0
}

func (p *Parser) switchAssignment3(tkn0, tkn1 token.Token, chr2 rune, tkn2 token.Token) token.Token {
	if p.chr == '=' {
		p.read()
		return tkn1
	}
	if p.chr == chr2 {
		p.read()
		return tkn2
	}
	return tkn0
}

func (p *Parser) switchAssignment4(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token) token.Token {
	if p.chr == '=' {
		p.read()
		return tkn1
	}
	if p.chr == chr2 {
		p.read()
		if p.chr == '=' {
			p.read()
			return tkn3
		}
		return tkn2
	}
	return tkn0
}

func (p *Parser) switchAssignment6(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token, chr3 rune, tkn4, tkn5 token.Token) token.Token {
	if p.chr == '=' {
		p.read()
		return tkn1
	}
	if p.chr == chr2 {
		p.read()
		if p.chr == '=' {
			p.read()
			return tkn3
		}
		if p.chr == chr3 {
			p.read()
			if p.chr == '=' {
				p.read()
				return tkn5
			}
			return tkn4
		}
		return tkn2
	}
	return tkn0
}

