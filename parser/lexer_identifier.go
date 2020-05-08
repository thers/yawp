package parser

import (
	"unicode"
)

// These functions are copies of
// https://github.com/evanw/esbuild/blob/766e48876293dfe6bb92cce906b1d4db4811b792/internal/lexer/lexer.go#L489

func isIdentifierStart(chr rune) bool {
	switch chr {
	case '_', '$',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return true
	}

	// All ASCII identifier start code points are listed above
	if chr < 0x7F {
		return false
	}

	return unicode.Is(idStart, chr)
}

func isIdentifierPart(chr rune) bool {
	switch chr {
	case '_', '$', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return true
	}

	// All ASCII identifier start code points are listed above
	if chr < 0x7F {
		return false
	}

	// ZWNJ and ZWJ are allowed in identifiers
	if chr == 0x200C || chr == 0x200D {
		return true
	}

	return unicode.Is(idContinue, chr)
}

func (p *Parser) scanIdentifier() (string, error) {
	offset := p.chrOffset

	for isIdentifierPart(p.chr) {
		p.read()
	}

	return p.src[offset:p.chrOffset], nil
}
