package parser

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"unicode/utf16"
	"yawp/parser/file"
	"yawp/parser/token"
)

type _chr struct {
	value rune
	width int
}

var matchIdentifier = regexp.MustCompile(`^[$_\p{L}][$_\p{L}\d}]*$`)

func isDecimalDigit(chr rune) bool {
	return '0' <= chr && chr <= '9'
}

func digitValue(chr rune) int {
	switch {
	case '0' <= chr && chr <= '9':
		return int(chr - '0')
	case 'a' <= chr && chr <= 'f':
		return int(chr - 'a' + 10)
	case 'A' <= chr && chr <= 'F':
		return int(chr - 'A' + 10)
	}
	return 16 // Larger than any legal digit value
}

func isDigit(chr rune, base int) bool {
	return digitValue(chr) < base
}

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
	parse := false
	for isIdentifierPart(p.chr) {
		if p.chr == '\\' {
			distance := p.chrOffset - offset
			p.read()
			if p.chr != 'u' {
				return "", fmt.Errorf("Invalid identifier escape character: %c (%s)", p.chr, string(p.chr))
			}
			parse = true
			var value rune
			for j := 0; j < 4; j++ {
				p.read()
				decimal, ok := hex2decimal(byte(p.chr))
				if !ok {
					return "", fmt.Errorf("Invalid identifier escape character: %c (%s)", p.chr, string(p.chr))
				}
				value = value<<4 | decimal
			}
			if value == '\\' {
				return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
			} else if distance == 0 {
				if !isIdentifierStart(value) {
					return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
				}
			} else if distance > 0 {
				if !isIdentifierPart(value) {
					return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
				}
			}
		}
		p.read()
	}
	literal := string(p.str[offset:p.chrOffset])
	if parse {
		return parseStringLiteral(literal)
	}
	return literal, nil
}

// 7.2
func isLineWhiteSpace(chr rune) bool {
	switch chr {
	case '\u0009', '\u000b', '\u000c', '\u0020', '\u00a0', '\ufeff':
		return true
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return false
	case '\u0085':
		return false
	}
	return unicode.IsSpace(chr)
}

// 7.3
func isLineTerminator(chr rune) bool {
	switch chr {
	case '\n', '\r', '\u2028', '\u2029':
		return true
	}
	return false
}

func (p *Parser) scan() (tkn token.Token, literal string, idx file.Idx) {

	p.implicitSemicolon = false

	for {
		p.skipWhiteSpace()

		idx = p.idxOf(p.chrOffset)
		insertSemicolon := false

		switch chr := p.chr; {
		case isIdentifierStart(chr):
			var err error
			literal, err = p.scanIdentifier()
			if err != nil {
				tkn = token.ILLEGAL
				break
			}
			if len(literal) > 1 {
				// Keywords are longer than 1 character, avoid lookup otherwise
				var strict bool
				tkn, strict = token.IsKeyword(literal)

				switch tkn {

				case 0: // Not a keyword
					if literal == "true" || literal == "false" {
						p.insertSemicolon = true
						tkn = token.BOOLEAN
						return
					} else if literal == "null" {
						p.insertSemicolon = true
						tkn = token.NULL
						return
					}

				case token.KEYWORD:
					tkn = token.KEYWORD
					if strict {
						// TODO If strict and in strict mode, then this is not a break
						break
					}
					return

				case
					token.THIS,
					token.BREAK,
					token.THROW, // A newline after a throw is not allowed, but we need to detect it
					token.RETURN,
					token.CONTINUE,
					token.DEBUGGER:
					p.insertSemicolon = true
					return

				default:
					return

				}
			}
			p.insertSemicolon = true
			tkn = token.IDENTIFIER
			return
		case '0' <= chr && chr <= '9':
			p.insertSemicolon = true
			tkn, literal = p.scanNumericLiteral(false)
			return
		default:
			p.read()
			switch chr {
			case -1:
				if p.insertSemicolon {
					p.insertSemicolon = false
					p.implicitSemicolon = true
				}
				tkn = token.EOF
			case '\r', '\n', '\u2028', '\u2029':
				p.insertSemicolon = false
				p.implicitSemicolon = true
				continue
			case '@':
				tkn = token.AT
			case '#':
				tkn = token.HASH
			case ':':
				tkn = token.COLON
			case '.':
				if digitValue(p.chr) < 10 {
					insertSemicolon = true
					tkn, literal = p.scanNumericLiteral(true)
				} else {
					if p.chr == '.' {
						p.read()
						if p.chr == '.' {
							p.read()
							tkn = token.DOTDOTDOT
						}
					} else {
						tkn = token.PERIOD
					}
				}
			case ',':
				tkn = token.COMMA
			case ';':
				tkn = token.SEMICOLON
			case '(':
				tkn = token.LEFT_PARENTHESIS
			case ')':
				tkn = token.RIGHT_PARENTHESIS
				insertSemicolon = true
			case '[':
				tkn = token.LEFT_BRACKET
			case ']':
				tkn = token.RIGHT_BRACKET
				insertSemicolon = true
			case '{':
				tkn = token.LEFT_BRACE
			case '}':
				tkn = token.RIGHT_BRACE
				insertSemicolon = true
			case '+':
				tkn = p.switchAssignment3(token.PLUS, token.ADD_ASSIGN, '+', token.INCREMENT)
				if tkn == token.INCREMENT {
					insertSemicolon = true
				}
			case '-':
				tkn = p.switchAssignment3(token.MINUS, token.SUBTRACT_ASSIGN, '-', token.DECREMENT)
				if tkn == token.DECREMENT {
					insertSemicolon = true
				}
			case '*':
				if p.chr == '*' {
					p.read()
					tkn = token.EXPONENTIATION
				} else {
					tkn = p.switchAssignment2(token.MULTIPLY, token.MULTIPLY_ASSIGN)
				}
			case '/':
				if p.chr == '/' {
					p.skipSingleLineComment()
					continue
				} else if p.chr == '*' {
					p.skipMultiLineComment()
					continue
				} else if p.chr == '>' {
					p.read()
					tkn = token.JSX_TAG_SELF_CLOSE
				} else {
					// Could be division, could be RegExp literal
					tkn = p.switchAssignment2(token.SLASH, token.QUOTIENT_ASSIGN)
					insertSemicolon = true
				}
			case '%':
				tkn = p.switchAssignment2(token.REMAINDER, token.REMAINDER_ASSIGN)
			case '^':
				tkn = p.switchAssignment2(token.EXCLUSIVE_OR, token.EXCLUSIVE_OR_ASSIGN)
			case '<':
				if p.chr == '>' {
					p.read()
					tkn = token.JSX_FRAGMENT_START
				} else if p.chr == '/' {
					p.read()

					if p.chr == '>' {
						p.read()
						tkn = token.JSX_FRAGMENT_END
					} else {
						tkn = token.JSX_TAG_CLOSE
					}
				} else {
					tkn = p.switchAssignment4(token.LESS, token.LESS_OR_EQUAL, '<', token.SHIFT_LEFT, token.SHIFT_LEFT_ASSIGN)
				}
			case '>':
				tkn = p.switchAssignment6(token.GREATER, token.GREATER_OR_EQUAL, '>', token.SHIFT_RIGHT, token.SHIFT_RIGHT_ASSIGN, '>', token.UNSIGNED_SHIFT_RIGHT, token.UNSIGNED_SHIFT_RIGHT_ASSIGN)
			case '=':
				if p.chr == '>' {
					p.read()
					tkn = token.ARROW
				} else {
					tkn = p.switchAssignment2(token.ASSIGN, token.EQUAL)
					if tkn == token.EQUAL && p.chr == '=' {
						p.read()
						tkn = token.STRICT_EQUAL
					}
				}
			case '!':
				tkn = p.switchAssignment2(token.NOT, token.NOT_EQUAL)
				if tkn == token.NOT_EQUAL && p.chr == '=' {
					p.read()
					tkn = token.STRICT_NOT_EQUAL
				}
			case '&':
				if p.chr == '^' {
					p.read()
					tkn = p.switchAssignment2(token.AND_NOT, token.AND_NOT_ASSIGN)
				} else {
					tkn = p.switchAssignment3(token.AND, token.AND_ASSIGN, '&', token.LOGICAL_AND)
				}
			case '|':
				tkn = p.switchAssignment3(token.OR, token.OR_ASSIGN, '|', token.LOGICAL_OR)
			case '~':
				tkn = token.BITWISE_NOT
			case '?':
				if p.chr == '.' {
					p.read()
					tkn = token.OPTIONAL_CHAINING
				} else if p.chr == '?' {
					p.read()
					tkn = token.NULLISH_COALESCING
				} else {
					tkn = token.QUESTION_MARK
				}
			case '`':
				tkn = token.TEMPLATE_QUOTE
			case '"', '\'':
				insertSemicolon = true
				tkn = token.STRING
				var err error
				literal, err = p.scanString(p.chrOffset - 1)
				if err != nil {
					tkn = token.ILLEGAL
				}
			default:
				p.errorUnexpected(idx, chr)
				tkn = token.ILLEGAL
			}
		}
		p.insertSemicolon = insertSemicolon
		return
	}
}

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

func (p *Parser) chrAt(index int) _chr {
	value, width := utf8.DecodeRuneInString(p.str[index:])
	return _chr{
		value: value,
		width: width,
	}
}

func (p *Parser) _peek() rune {
	if p.offset+1 < p.length {
		return rune(p.str[p.offset+1])
	}
	return -1
}

type ParserPartialState struct {
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

func (p *Parser) getPartialState() *ParserPartialState {
	return &ParserPartialState{
		idx:               p.idx,
		token:             p.token,
		errors:            p.errors,
		offset:            p.offset,
		chrOffset:         p.chrOffset,
		nextChr:           p.chr,
		literal:           p.literal,
		parsedStr:         p.parsedStr,
		insertSemicolon:   p.insertSemicolon,
		implicitSemicolon: p.implicitSemicolon,
	}
}

func (p *Parser) restorePartialState(state *ParserPartialState) {
	p.idx = state.idx
	p.token = state.token
	p.errors = state.errors
	p.offset = state.offset
	p.chrOffset = state.chrOffset
	p.chr = state.nextChr
	p.literal = state.literal
	p.parsedStr = state.parsedStr
	p.insertSemicolon = state.insertSemicolon
	p.implicitSemicolon = state.implicitSemicolon
}

func (p *Parser) read() {
	p.parsedStr += string(p.chr)

	if p.offset < p.length {
		p.chrOffset = p.offset
		chr, width := rune(p.str[p.offset]), 1
		if chr >= utf8.RuneSelf { // !ASCII
			chr, width = utf8.DecodeRuneInString(p.str[p.offset:])
			if chr == utf8.RuneError && width == 1 {
				p.error(p.chrOffset, "Invalid UTF-8 character")
			}
		}
		p.offset += width
		p.chr = chr
	} else {
		p.chrOffset = p.length
		p.chr = -1 // EOF
	}
}

// This is here since the functions are so similar
func (self *_RegExp_parser) read() {
	if self.offset < self.length {
		self.chrOffset = self.offset
		chr, width := rune(self.str[self.offset]), 1
		if chr >= utf8.RuneSelf { // !ASCII
			chr, width = utf8.DecodeRuneInString(self.str[self.offset:])
			if chr == utf8.RuneError && width == 1 {
				self.error(self.chrOffset, "Invalid UTF-8 character")
			}
		}
		self.offset += width
		self.chr = chr
	} else {
		self.chrOffset = self.length
		self.chr = -1 // EOF
	}
}

func (p *Parser) skipSingleLineComment() {
	for p.chr != -1 {
		p.read()
		if isLineTerminator(p.chr) {
			return
		}
	}
}

func (p *Parser) skipMultiLineComment() {
	p.read()
	for p.chr >= 0 {
		chr := p.chr
		p.read()
		if chr == '*' && p.chr == '/' {
			p.read()
			return
		}
	}

	p.errorUnexpected(0, p.chr)
}

func (p *Parser) skipWhiteSpace() {
	for {
		switch p.chr {
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			p.read()
			continue
		case '\r':
			if p._peek() == '\n' {
				p.read()
			}
			fallthrough
		case '\u2028', '\u2029', '\n':
			if p.insertSemicolon {
				return
			}
			p.read()
			continue
		}
		if p.chr >= utf8.RuneSelf {
			if unicode.IsSpace(p.chr) {
				p.read()
				continue
			}
		}
		break
	}
}

func (p *Parser) skipLineWhiteSpace() {
	for isLineWhiteSpace(p.chr) {
		p.read()
	}
}

func (p *Parser) scanNumberRemainder(base int) {
	for digitValue(p.chr) < base {
		p.read()
	}
}

func (p *Parser) scanEscape(quote rune) {

	var length, base uint32
	switch p.chr {
	//case '0', '1', '2', '3', '4', '5', '6', '7':
	//    Octal:
	//    length, base, limit = 3, 8, 255
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'', '0':
		p.read()
		return
	case '\r', '\n', '\u2028', '\u2029':
		p.scanNewline()
		return
	case 'x':
		p.read()
		length, base = 2, 16
	case 'u':
		p.read()
		length, base = 4, 16
	default:
		p.read() // Always make progress
		return
	}

	var value uint32
	for ; length > 0 && p.chr != quote && p.chr >= 0; length-- {
		digit := uint32(digitValue(p.chr))
		if digit >= base {
			break
		}
		value = value*base + digit
		p.read()
	}
}

func (p *Parser) scanString(offset int) (string, error) {
	// " '  /
	quote := rune(p.str[offset])

	for p.chr != quote {
		chr := p.chr
		if isLineTerminator(chr) || chr < 0 {
			goto newline
		}
		p.read()
		if chr == '\\' {
			if quote == '/' {
				if isLineTerminator(chr) || p.chr < 0 {
					goto newline
				}
				p.read()
			} else {
				p.scanEscape(quote)
			}
		} else if chr == '[' && quote == '/' {
			// Allow a slash (/) in a bracket character class ([...])
			// TODO Fix this, this is hacky...
			quote = -1
		} else if chr == ']' && quote == -1 {
			quote = '/'
		}
	}

	// " '  /
	p.read()

	return p.str[offset:p.chrOffset], nil

newline:
	p.scanNewline()
	err := "String not terminated"
	if quote == '/' {
		err = "Invalid regular expression: missing /"
		p.error(p.idxOf(offset), err)
	}
	return "", errors.New(err)
}

func (p *Parser) scanNewline() {
	if p.chr == '\r' {
		p.read()
		if p.chr != '\n' {
			return
		}
	}
	p.read()
}

func hex2decimal(chr byte) (value rune, ok bool) {
	{
		chr := rune(chr)
		switch {
		case '0' <= chr && chr <= '9':
			return chr - '0', true
		case 'a' <= chr && chr <= 'f':
			return chr - 'a' + 10, true
		case 'A' <= chr && chr <= 'F':
			return chr - 'A' + 10, true
		}
		return
	}
}

func parseNumberLiteral(literal string) (value interface{}, err error) {
	// TODO Is Uint okay? What about -MAX_UINT
	value, err = strconv.ParseInt(literal, 0, 64)
	if err == nil {
		return
	}

	parseIntErr := err // Save this first error, just in case

	value, err = strconv.ParseFloat(literal, 64)
	if err == nil {
		return
	} else if err.(*strconv.NumError).Err == strconv.ErrRange {
		// Infinity, etc.
		return value, nil
	}

	err = parseIntErr

	if err.(*strconv.NumError).Err == strconv.ErrRange {
		if len(literal) > 2 && literal[0] == '0' && (literal[1] == 'X' || literal[1] == 'x') {
			// Could just be a very large number (e.g. 0x8000000000000000)
			var value float64
			literal = literal[2:]
			for _, chr := range literal {
				digit := digitValue(chr)
				if digit >= 16 {
					goto error
				}
				value = value*16 + float64(digit)
			}
			return value, nil
		}
	}

error:
	return nil, errors.New("Illegal numeric literal")
}

func parseStringLiteral(literal string) (string, error) {
	// Best case scenario...
	if literal == "" {
		return "", nil
	}

	// Slightly less-best case scenario...
	if !strings.ContainsRune(literal, '\\') {
		return literal, nil
	}

	str := literal
	buffer := bytes.NewBuffer(make([]byte, 0, 3*len(literal)/2))
	var surrogate rune
S:
	for len(str) > 0 {
		switch chr := str[0]; {
		// We do not explicitly handle the case of the quote
		// value, which can be: " ' /
		// This assumes we're already passed a partially well-formed literal
		case chr >= utf8.RuneSelf:
			chr, size := utf8.DecodeRuneInString(str)
			buffer.WriteRune(chr)
			str = str[size:]
			continue
		case chr != '\\':
			buffer.WriteByte(chr)
			str = str[1:]
			continue
		}

		if len(str) <= 1 {
			panic("len(str) <= 1")
		}
		chr := str[1]
		var value rune
		if chr >= utf8.RuneSelf {
			str = str[1:]
			var size int
			value, size = utf8.DecodeRuneInString(str)
			str = str[size:] // \ + <character>
		} else {
			str = str[2:] // \<character>
			switch chr {
			case 'b':
				value = '\b'
			case 'f':
				value = '\f'
			case 'n':
				value = '\n'
			case 'r':
				value = '\r'
			case 't':
				value = '\t'
			case 'v':
				value = '\v'
			case 'x', 'u':
				size := 0
				switch chr {
				case 'x':
					size = 2
				case 'u':
					size = 4
				}
				if len(str) < size {
					return "", fmt.Errorf("invalid escape: \\%s: len(%q) != %d", string(chr), str, size)
				}
				for j := 0; j < size; j++ {
					decimal, ok := hex2decimal(str[j])
					if !ok {
						return "", fmt.Errorf("invalid escape: \\%s: %q", string(chr), str[:size])
					}
					value = value<<4 | decimal
				}
				str = str[size:]
				if chr == 'x' {
					break
				}
				if value > utf8.MaxRune {
					panic("value > utf8.MaxRune")
				}
			case '0':
				if len(str) == 0 || '0' > str[0] || str[0] > '7' {
					value = 0
					break
				}
				fallthrough
			case '1', '2', '3', '4', '5', '6', '7':
				// TODO strict
				value = rune(chr) - '0'
				j := 0
				for ; j < 2; j++ {
					if len(str) < j+1 {
						break
					}
					chr := str[j]
					if '0' > chr || chr > '7' {
						break
					}
					decimal := rune(str[j]) - '0'
					value = (value << 3) | decimal
				}
				str = str[j:]
			case '\\':
				value = '\\'
			case '\'', '"':
				value = rune(chr)
			case '\r':
				if len(str) > 0 {
					if str[0] == '\n' {
						str = str[1:]
					}
				}
				fallthrough
			case '\n':
				continue
			default:
				value = rune(chr)
			}
			if surrogate != 0 {
				value = utf16.DecodeRune(surrogate, value)
				surrogate = 0
			} else {
				if utf16.IsSurrogate(value) {
					surrogate = value
					continue S
				}
			}
		}
		buffer.WriteRune(value)
	}

	return buffer.String(), nil
}

func (p *Parser) scanNumericLiteral(decimalPoint bool) (token.Token, string) {

	offset := p.chrOffset
	tkn := token.NUMBER

	if decimalPoint {
		offset--
		p.scanNumberRemainder(10)
		goto exponent
	}

	if p.chr == '0' {
		offset := p.chrOffset
		p.read()
		if p.chr == 'x' || p.chr == 'X' {
			// Hexadecimal
			p.read()
			if isDigit(p.chr, 16) {
				p.read()
			} else {
				return token.ILLEGAL, p.str[offset:p.chrOffset]
			}
			p.scanNumberRemainder(16)

			if p.chrOffset-offset <= 2 {
				// Only "0x" or "0X"
				p.error(0, "Illegal hexadecimal number")
			}

			goto hexadecimal
		} else if p.chr == 'b' || p.chr == 'B' {
			// Binary
			p.read()
			if isDigit(p.chr, 2) {
				p.read()
			} else {
				return token.ILLEGAL, p.str[offset:p.chrOffset]
			}
			p.scanNumberRemainder(2)

			if p.chrOffset-offset <= 2 {
				// Only "0b" or "0B"
				p.error(0, "Illegal binary number")
			}

			goto binary
		} else if p.chr == '.' {
			// Float
			goto float
		} else {
			// Octal, Float
			if p.chr == 'e' || p.chr == 'E' {
				goto exponent
			}
			p.scanNumberRemainder(8)
			if p.chr == '8' || p.chr == '9' {
				return token.ILLEGAL, p.str[offset:p.chrOffset]
			}
			goto octal
		}
	}

	p.scanNumberRemainder(10)

float:
	if p.chr == '.' {
		p.read()
		p.scanNumberRemainder(10)
	}

exponent:
	if p.chr == 'e' || p.chr == 'E' {
		p.read()
		if p.chr == '-' || p.chr == '+' {
			p.read()
		}
		if isDecimalDigit(p.chr) {
			p.read()
			p.scanNumberRemainder(10)
		} else {
			return token.ILLEGAL, p.str[offset:p.chrOffset]
		}
	}

binary:
hexadecimal:
octal:
	if isIdentifierStart(p.chr) || isDecimalDigit(p.chr) {
		return token.ILLEGAL, p.str[offset:p.chrOffset]
	}

	return tkn, p.str[offset:p.chrOffset]
}
