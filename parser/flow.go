package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeAnnotation() ast.FlowType {
	p.consumeExpected(token.COLON)

	start := p.idx

	switch p.token {
	case token.TYPE_BOOLEAN:
		p.next()

		return &ast.FlowBooleanType{
			Start: start,
		}
	case token.BOOLEAN:
		kind := p.literal
		p.next()

		if kind == "true" {
			return &ast.FlowTrueType{
				Start: start,
			}
		} else {
			return &ast.FlowFalseType{
				Start: start,
			}
		}
	case token.STRING:
		str := p.literal[1 : len(p.literal) - 1]
		p.next()

		return &ast.FlowStringType{
			Start: start,
			String: str,
		}
	case token.NUMBER:
		number, err := parseNumberLiteral(p.literal)
		p.next()

		if err != nil {
			p.error(start, err.Error())
		} else {
			return &ast.FlowNumberType{
				Start:  start,
				Number: number,
			}
		}
	}

	return nil
}
