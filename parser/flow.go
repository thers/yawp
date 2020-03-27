package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeIdentifier() *ast.FlowIdentifier {
	identifier := p.parseIdentifier()

	return &ast.FlowIdentifier{
		Start: identifier.Start,
		Name:  identifier.Name,
	}
}

func (p *Parser) parseFlowType() ast.FlowType {
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

		return &ast.FlowStringLiteralType{
			Start: start,
			String: str,
		}
	case token.NUMBER:
		number, err := parseNumberLiteral(p.literal)
		p.next()

		if err != nil {
			p.error(start, err.Error())
		} else {
			return &ast.FlowNumberLiteralType{
				Start:  start,
				Number: number,
			}
		}
	case token.IDENTIFIER:
		return p.parseFlowTypeIdentifier()
	case token.TYPEOF:
		start := p.consumeExpected(token.TYPEOF)

		return &ast.FlowTypeOfType{
			Start:      start,
			Identifier: p.parseFlowTypeIdentifier(),
		}
	case token.TYPE_STRING:
		start := p.consumeExpected(token.TYPE_STRING)

		return &ast.FlowStringType{
			Start:start,
		}
	case token.TYPE_NUMBER:
		start := p.consumeExpected(token.TYPE_NUMBER)

		return &ast.FlowNumberType{
			Start:start,
		}
	case token.QUESTION_MARK:
		p.consumeExpected(token.QUESTION_MARK)

		return &ast.FlowOptionalType{
			FlowType: p.parseFlowType(),
		}
	}

	return nil
}

func (p *Parser) parseFlowTypeAnnotation() ast.FlowType {
	p.consumeExpected(token.COLON)

	return p.parseFlowType()
}
