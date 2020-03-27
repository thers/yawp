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
		str := p.literal[1 : len(p.literal)-1]
		p.next()

		return &ast.FlowStringLiteralType{
			Start:  start,
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
	case token.TYPE_ANY:
		p.next()

		return &ast.FlowAnyType{
			Start: start,
		}
	case token.TYPEOF:
		p.next()

		return &ast.FlowTypeOfType{
			Start:      start,
			Identifier: p.parseFlowTypeIdentifier(),
		}
	case token.TYPE_STRING:
		p.next()

		return &ast.FlowStringType{
			Start: start,
		}
	case token.TYPE_NUMBER:
		p.next()

		return &ast.FlowNumberType{
			Start: start,
		}
	case token.QUESTION_MARK:
		p.next()

		return &ast.FlowOptionalType{
			FlowType: p.parseFlowType(),
		}
	case token.VOID:
		p.next()

		return &ast.FlowVoidType{
			Start: start,
		}
	case token.NULL:
		p.next()

		return &ast.FlowNullType{
			Start: start,
		}
	case token.LEFT_BRACE:
		return p.parseFlowInexactObjectType()
	case token.TYPE_EXACT_OBJECT_START:
		return p.parseFlowExactObjectType()
	case token.LEFT_BRACKET:
		return p.parseFlowTupleType()
	}

	return nil
}

func (p *Parser) parseFlowTypeAnnotation() ast.FlowType {
	p.consumeExpected(token.COLON)

	return p.parseFlowType()
}
