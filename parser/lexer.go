package parser

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"

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
	return 16 // Larger than isAny legal digit value
}

func isDigit(chr rune, base int) bool {
	return digitValue(chr) < base
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

	p.tokenIsKeyword = false
	p.implicitSemicolon = false
	p.newLineBeforeCurrentToken = false

	for {
		p.skipWhiteSpace()

		p.tokenCol = p.chrCol
		idx = file.Idx(p.chrOffset)
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

				p.tokenIsKeyword = tkn > 0

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

				// some special keywords can only be used as keywords in known context
				case
					token.OF,
					token.AS,
					token.FROM:
						p.insertSemicolon = true
						tkn = token.IDENTIFIER
						return

				case token.TYPE_TYPE:
					p.insertSemicolon = true
					if !p.scope.inModuleRoot() {
						tkn = token.IDENTIFIER
					}
					return

				// following are keywords only when in type
				case
					token.TYPE_STRING,
					token.TYPE_NUMBER,
					token.TYPE_MIXED,
					token.TYPE_BOOLEAN,
					token.TYPE_ANY,
					token.TYPE_OPAQUE:
						if !p.scope.inType {
							tkn = token.IDENTIFIER
						}
						return

				case token.AWAIT:
					p.insertSemicolon = true
					if !p.scope.allowAwait {
						tkn = token.IDENTIFIER
					}
					return

				case token.YIELD:
					p.insertSemicolon = true
					if !p.scope.allowYield {
						tkn = token.IDENTIFIER
					}
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
				if p.chr == '|' {
					p.read()

					tkn = token.TYPE_EXACT_OBJECT_START
				} else {
					tkn = token.LEFT_BRACE
				}
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
				// shifts are not allowed inside of type
				if p.scope.inType {
					tkn = token.GREATER
				} else {
					tkn = p.switchAssignment6(
						token.GREATER,
						token.GREATER_OR_EQUAL,
						'>',
						token.SHIFT_RIGHT,
						token.SHIFT_RIGHT_ASSIGN,
						'>',
						token.UNSIGNED_SHIFT_RIGHT,
						token.UNSIGNED_SHIFT_RIGHT_ASSIGN,
					)
				}
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
				if p.chr == '}' {
					p.read()

					tkn = token.TYPE_EXACT_OBJECT_END
				} else {
					tkn = p.switchAssignment3(token.OR, token.OR_ASSIGN, '|', token.LOGICAL_OR)
				}
			case '~':
				tkn = token.BITWISE_NOT
			case '?':
				if p.chr == '.' {
					// excluding foo?.2:1
					if !isDigit(p.peekChr(), 10) {
						p.read()
						tkn = token.OPTIONAL_CHAINING
					} else {
						tkn = token.QUESTION_MARK
					}
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
				p.errorUnexpected(p.loc(), chr)
				tkn = token.ILLEGAL
			}
		}
		p.insertSemicolon = insertSemicolon
		return
	}
}

func (p *Parser) chrAt(index int) _chr {
	value, width := utf8.DecodeRuneInString(p.src[index:])
	return _chr{
		value: value,
		width: width,
	}
}

func (p *Parser) _peek() rune {
	if p.nextChrOffset+1 < p.length {
		return rune(p.src[p.nextChrOffset+1])
	}
	return -1
}

func (p *Parser) peekChr() rune {
	if p.nextChrOffset >= p.length {
		return -1
	}

	return rune(p.src[p.nextChrOffset])
}

func (p *Parser) read() {
	p.parsedSrc = p.src[:p.nextChrOffset]

	if p.advanceLine {
		p.advanceLine = false
		p.line++

		// shouldn't ever be the case as it's for \n rune
		//p.chrCol = 0
	} else {
		p.chrCol += p.nextChrOffset - p.chrOffset
	}

	if p.nextChrOffset < p.length {
		p.chrOffset = p.nextChrOffset
		chr, width := rune(p.src[p.nextChrOffset]), 1

		if chr >= utf8.RuneSelf { // !ASCII
			chr, width = utf8.DecodeRuneInString(p.src[p.nextChrOffset:])
			if chr == utf8.RuneError && width == 1 {
				p.error(p.loc(), "Invalid UTF-8 character")
			}
		}

		p.nextChrOffset += width
		p.chr = chr

		switch chr {
		case '\n', '\u2028', '\u2029':
			p.chrCol = 0
			p.advanceLine = true
			p.newLineBeforeCurrentToken = true
		}
	} else {
		p.chrOffset = p.length
		p.chr = -1 // EOF
	}

	return
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

	p.errorUnexpected(p.loc(), p.chr)
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
	quote := rune(p.src[offset])

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

	return p.src[offset:p.chrOffset], nil

newline:
	p.scanNewline()
	err := "String not terminated"
	if quote == '/' {
		err = "Invalid regular expression: missing /"
		p.error(p.loc(), err)
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

func parseNumberLiteralIntoNumber(literal string) (value interface{}, err error) {
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
				return token.ILLEGAL, p.src[offset:p.chrOffset]
			}
			p.scanNumberRemainder(16)

			if p.chrOffset-offset <= 2 {
				// Only "0x" or "0X"
				p.error(p.loc(), "Illegal hexadecimal number")
			}

			goto hexadecimal
		} else if p.chr == 'b' || p.chr == 'B' {
			// Binary
			p.read()
			if isDigit(p.chr, 2) {
				p.read()
			} else {
				return token.ILLEGAL, p.src[offset:p.chrOffset]
			}
			p.scanNumberRemainder(2)

			if p.chrOffset-offset <= 2 {
				// Only "0b" or "0B"
				p.error(p.loc(), "Illegal binary number")
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
				return token.ILLEGAL, p.src[offset:p.chrOffset]
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
			return token.ILLEGAL, p.src[offset:p.chrOffset]
		}
	}

binary:
hexadecimal:
octal:
	if isIdentifierStart(p.chr) || isDecimalDigit(p.chr) {
		return token.ILLEGAL, p.src[offset:p.chrOffset]
	}

	return tkn, p.src[offset:p.chrOffset]
}
