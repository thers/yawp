package parser

import (
	"fmt"
	"yawp/parser/file"
	"yawp/parser/token"
)

const (
	err_UnexpectedToken      = "Unexpected token %v"
	err_UnexpectedEndOfInput = "Unexpected end of input"
)


type SyntaxError struct {
	Loc *file.Loc
	error error
}

func (se SyntaxError) Error() error {
	return fmt.Errorf("%d:%d %s", se.Loc.Line, se.Loc.Col, se.error)
}


func (p *Parser) error(loc *file.Loc, msg string, msgValues ...interface{}) {
	message := fmt.Errorf(msg, msgValues...)

	panic(&SyntaxError{
		Loc: loc,
		error: message,
	})
}

func (p *Parser) errorUnexpected(loc *file.Loc, chr rune) {
	if chr == -1 {
		p.error(loc, err_UnexpectedEndOfInput)
	} else {
		p.error(loc, err_UnexpectedToken, token.ILLEGAL)
	}
}

func (p *Parser) errorUnexpectedToken(tkn token.Token) {
	p.errorUnexpectedTokenAt(tkn, p.loc())
}

func (p *Parser) errorUnexpectedTokenAt(tkn token.Token, at *file.Loc) {
	value := tkn.String()

	switch tkn {
	case token.EOF:
		p.error(p.loc(), err_UnexpectedEndOfInput)
	case token.BOOLEAN, token.NULL:
		p.error(at, err_UnexpectedToken, value)
	case token.IDENTIFIER:
		p.error(at, "Unexpected identifier")
	case token.KEYWORD:
		p.error(at, "Unexpected keyword")
	case token.NUMBER:
		p.error(at, "Unexpected number")
	case token.STRING:
		p.error(at, "Unexpected string")
	}

	p.error(at, err_UnexpectedToken, value)
}

func (p *Parser) unexpectedToken() {
	p.errorUnexpectedToken(p.token)
}

func (p *Parser) unexpectedTokenAt(at *file.Loc) {
	p.errorUnexpectedTokenAt(p.token, at)
}
