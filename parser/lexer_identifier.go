package parser

import (
	"unicode"
	"unicode/utf8"
)

func isIdentifierStart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		chr >= utf8.RuneSelf && unicode.IsLetter(chr)
}

func isIdentifierPart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		'0' <= chr && chr <= '9' ||
		chr >= utf8.RuneSelf && (unicode.IsLetter(chr) || unicode.IsDigit(chr))
}

func (p *Parser) scanIdentifier() (string, error) {
	offset := p.chrOffset
	//parse := false

	for isIdentifierPart(p.chr) {
		//if p.chr == '\\' {
		//	distance := p.chrOffset - offset
		//	p.read()
		//	if p.chr != 'u' {
		//		return "", fmt.Errorf("Invalid identifier escape character: %c (%s)", p.chr, string(p.chr))
		//	}
		//	parse = true
		//	var value rune
		//	for j := 0; j < 4; j++ {
		//		p.read()
		//		decimal, ok := hex2decimal(byte(p.chr))
		//		if !ok {
		//			return "", fmt.Errorf("Invalid identifier escape character: %c (%s)", p.chr, string(p.chr))
		//		}
		//		value = value<<4 | decimal
		//	}
		//	if value == '\\' {
		//		return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
		//	} else if distance == 0 {
		//		if !isIdentifierStart(value) {
		//			return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
		//		}
		//	} else if distance > 0 {
		//		if !isIdentifierPart(value) {
		//			return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
		//		}
		//	}
		//}
		p.read()
	}

	return p.src[offset:p.chrOffset], nil
	//if parse {
	//	return parseStringLiteral(literal)
	//}
	//return literal, nil
}
