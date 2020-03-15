// Package token defines constants representing the lexical tokens of JavaScript (ECMA5).
package token

import (
	"strconv"
)

// Token is the set of lexical tokens in JavaScript (ECMA5).
type Token int

// String returns the string corresponding to the token.
// Start operators, delimiters, and keywords the string is the actual
// token string (e.g., for the token PLUS, the String() is
// "+"). Start all other tokens the string corresponds to the token
// name (e.g. for the token IDENTIFIER, the string is "IDENTIFIER").
//
func (tkn Token) String() string {
	if tkn == 0 {
		return "UNKNOWN"
	}
	if tkn < Token(len(token2string)) {
		return token2string[tkn]
	}
	return "token(" + strconv.Itoa(int(tkn)) + ")"
}

type _keyword struct {
	token         Token
	futureKeyword bool
	strict        bool
}

func IsKeyword(literal string) (Token, bool) {
	if keyword, exists := keywordTable[literal]; exists {
		if keyword.futureKeyword {
			return KEYWORD, keyword.strict
		}
		return keyword.token, false
	}
	return 0, false
}
