package parser

import (
	"fmt"
	"sort"

	"yawp/parser/file"
	"yawp/parser/token"
)

const (
	err_UnexpectedToken      = "Unexpected token %v"
	err_UnexpectedEndOfInput = "Unexpected end of input"
	err_UnexpectedEscape     = "Unexpected escape"
)

//    UnexpectedNumber:  'Unexpected number',
//    UnexpectedString:  'Unexpected string',
//    UnexpectedIdentifier:  'Unexpected identifier',
//    UnexpectedReserved:  'Unexpected reserved word',
//    NewlineAfterThrow:  'Illegal newline after throw',
//    InvalidRegExp: 'Invalid regular expression',
//    UnterminatedRegExp:  'Invalid regular expression: missing /',
//    InvalidLHSInAssignment:  'Invalid left-hand side in assignment',
//    InvalidLHSInForIn:  'Invalid left-hand side in for-in',
//    MultipleDefaultsInSwitch: 'More than one default clause in switch statement',
//    NoCatchOrFinally:  'Missing catch or finally after try',
//    UnknownLabel: 'Undefined label \'%0\'',
//    Redeclaration: '%0 \'%1\' has already been declared',
//    IllegalContinue: 'Illegal continue statement',
//    IllegalBreak: 'Illegal break statement',
//    IllegalReturn: 'Illegal return statement',
//    StrictModeWith:  'Strict mode code may not include a with statement',
//    StrictCatchVariable:  'Catch variable may not be eval or arguments in strict mode',
//    StrictVarName:  'Variable name may not be eval or arguments in strict mode',
//    StrictParamName:  'Parameter name eval or arguments is not allowed in strict mode',
//    StrictParamDupe: 'Strict mode function may not have duplicate parameter names',
//    StrictFunctionName:  'Function name may not be eval or arguments in strict mode',
//    StrictOctalLiteral:  'Octal literals are not allowed in strict mode.',
//    StrictDelete:  'Delete of an unqualified identifier in strict mode.',
//    StrictDuplicateProperty:  'Duplicate data property in object literal not allowed in strict mode',
//    AccessorDataProperty:  'Object literal may not have data and accessor property with the same name',
//    AccessorGetSet:  'Object literal may not have multiple get/set accessors with the same name',
//    StrictLHSAssignment:  'Assignment to eval or arguments is not allowed in strict mode',
//    StrictLHSPostfix:  'Postfix increment/decrement may not have eval or arguments operand in strict mode',
//    StrictLHSPrefix:  'Prefix increment/decrement may not have eval or arguments operand in strict mode',
//    StrictReservedWord:  'Use of future reserved word in strict mode'

// A SyntaxError is a description of an ECMAScript syntax error.

// An Error represents a parsing error. It includes the position where the error occurred and a message/description.
type Error struct {
	Position file.Position
	Message  string
}

// FIXME Should this be "SyntaxError"?

func (self Error) Error() string {
	filename := self.Position.Filename
	if filename == "" {
		filename = "(anonymous)"
	}
	return fmt.Sprintf("%s: Line %d:%d %s",
		filename,
		self.Position.Line,
		self.Position.Column,
		self.Message,
	)
}

func (p *Parser) error(place interface{}, msg string, msgValues ...interface{}) *Error {
	idx := file.Idx(0)
	switch place := place.(type) {
	case int:
		idx = p.idxOf(place)
	case file.Idx:
		if place == 0 {
			idx = p.idxOf(p.chrOffset)
		} else {
			idx = place
		}
	default:
		panic(fmt.Errorf("error(%T, ...)", place))
	}

	position := p.position(idx)
	msg = fmt.Sprintf(msg, msgValues...)
	p.errors.Add(position, msg)
	return p.errors[len(p.errors)-1]
}

func (p *Parser) errorUnexpected(idx file.Idx, chr rune) error {
	if chr == -1 {
		return p.error(idx, err_UnexpectedEndOfInput)
	}
	return p.error(idx, err_UnexpectedToken, token.ILLEGAL)
}

func (p *Parser) errorUnexpectedToken(tkn token.Token) error {
	switch tkn {
	case token.EOF:
		return p.error(file.Idx(0), err_UnexpectedEndOfInput)
	}
	value := tkn.String()
	switch tkn {
	case token.BOOLEAN, token.NULL:
		value = p.literal
	case token.IDENTIFIER:
		return p.error(p.idx, "Unexpected identifier")
	case token.KEYWORD:
		// TODO Might be a future reserved word
		return p.error(p.idx, "Unexpected reserved word")
	case token.NUMBER:
		return p.error(p.idx, "Unexpected number")
	case token.STRING:
		return p.error(p.idx, "Unexpected string")
	}
	return p.error(p.idx, err_UnexpectedToken, value)
}

// ErrorList is a list of *Errors.
//
type ErrorList []*Error

// Add adds an Error with given position and message to an ErrorList.
func (self *ErrorList) Add(position file.Position, msg string) {
	*self = append(*self, &Error{position, msg})
}

// Reset resets an ErrorList to no errors.
func (self *ErrorList) Reset() { *self = (*self)[0:0] }

func (self ErrorList) Len() int      { return len(self) }
func (self ErrorList) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self ErrorList) Less(i, j int) bool {
	x := &self[i].Position
	y := &self[j].Position
	if x.Filename < y.Filename {
		return true
	}
	if x.Filename == y.Filename {
		if x.Line < y.Line {
			return true
		}
		if x.Line == y.Line {
			return x.Column < y.Column
		}
	}
	return false
}

func (self ErrorList) Sort() {
	sort.Sort(self)
}

// Error implements the Error interface.
func (self ErrorList) Error() string {
	switch len(self) {
	case 0:
		return "no errors"
	case 1:
		return self[0].Error()
	}

	str := ""

	for index, err := range self {
		str += fmt.Sprintf("\n%d: %s", index, err.Error())
	}

	return str
}

// Err returns an error equivalent to this ErrorList.
// If the list is empty, Err returns nil.
func (self ErrorList) Err() error {
	if len(self) == 0 {
		return nil
	}
	return self
}
