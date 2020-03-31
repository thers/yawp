package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeIdentifier() *ast.FlowIdentifier {
	identifier := p.parseIdentifier()

	return &ast.FlowIdentifier{
		Start: identifier.Start,
		Name:  identifier.Name,
	}
}

func (p *Parser) parseFlowTypeIdentifierIncludingKeywords() *ast.FlowIdentifier {
	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		return nil
	}

	return &ast.FlowIdentifier{
		Start: identifier.Start,
		Name:  identifier.Name,
	}
}

func (p *Parser) parseFlowType() ast.FlowType {
	start := p.idx

	switch p.token {
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
	case token.TYPE_BOOLEAN, token.TYPE_ANY, token.TYPE_STRING, token.TYPE_NUMBER, token.VOID, token.NULL, token.TYPE_MIXED:
		kind := p.token
		end := start + file.Idx(len(p.literal))
		p.next()

		return &ast.FlowPrimitiveType{
			Start: start,
			End:   end,
			Kind:  kind,
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
	case token.TYPEOF:
		p.next()

		return &ast.FlowTypeOfType{
			Start:      start,
			Identifier: p.parseFlowTypeIdentifier(),
		}
	case token.QUESTION_MARK:
		p.next()

		return &ast.FlowOptionalType{
			FlowType: p.parseFlowType(),
		}
	case token.MULTIPLY:
		p.next()

		return &ast.FlowExistentialType{
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
