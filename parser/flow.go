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

func (p *Parser) parseFlowObjectProperties(exact bool) []ast.FlowObjectProperty {
	var terminator token.Token

	props := make([]ast.FlowObjectProperty, 0)

	if exact {
		terminator = token.TYPE_EXACT_OBJECT_END
	} else {
		terminator = token.RIGHT_BRACE
	}

	for p.until(terminator) {
		start := p.idx

		if p.isIdentifierOrKeyword() {
			identifier := p.parseIdentifierIncludingKeywords()

			if identifier == nil {
				p.unexpectedToken()
				p.next()

				goto Next
			}

			prop := &ast.FlowNamedObjectProperty{
				Start:    start,
				Optional: false,
				Name:     identifier.Name,
				Value:    nil,
			}

			if p.is(token.QUESTION_MARK) {
				p.next()
				prop.Optional = true
			}

			prop.Value = p.parseFlowTypeAnnotation()

			props = append(props, prop)
		} else if p.is(token.DOTDOTDOT) {
			idx := p.idx
			p.next()

			// inexact object specifier
			if p.is(token.COMMA) || p.is(terminator) {
				if exact {
					// explicit inexact disallowed in exact
					p.error(idx, "Explicit inexact syntax is not allowed in exact object type")
					break
				}

				props = append(props, &ast.FlowInexactSpecifierProperty{
					Start: start,
				})
				break
			} else {
				props = append(props, &ast.FlowSpreadObjectProperty{
					Start:    start,
					FlowType: p.parseFlowType(),
				})
			}
		} else if p.is(token.LEFT_BRACKET) {
			p.next()

			prop := &ast.FlowIndexerObjectProperty{
				Start:   start,
				KeyName: "",
				KeyType: nil,
				Value:   nil,
			}

			if p.isIdentifierOrKeyword() {
				isKeyword := p.isKeyword
				identifier := p.parseIdentifierIncludingKeywords()

				if identifier == nil {
					p.unexpectedToken()
					p.next()
					goto Next
				}

				if p.is(token.COLON) {
					prop.KeyName = identifier.Name
					prop.KeyType = p.parseFlowTypeAnnotation()
				} else if p.is(token.RIGHT_BRACKET) {
					if isKeyword {
						// keywords aren't legal identifiers
						p.error(identifier.Start, "Cannot use keyword as type identifier")
						p.next()
						goto Next
					}

					prop.KeyType = &ast.FlowIdentifier{
						Start: identifier.Start,
						Name:  identifier.Name,
					}
				} else {
					p.unexpectedToken()
					p.next()
				}
			} else if p.is(token.RIGHT_BRACKET) {
				p.unexpectedToken()
			} else {
				prop.KeyType = p.parseFlowType()
			}

			p.consumeExpected(token.RIGHT_BRACKET)

			prop.Value = p.parseFlowTypeAnnotation()

			props = append(props, prop)
		} else {
			p.unexpectedToken()
			p.next()
		}

	Next:
		p.consumePossible(token.COMMA)
	}

	return props
}

func (p *Parser) parseFlowInexactObjectType() *ast.FlowInexactObject {
	start := p.consumeExpected(token.LEFT_BRACE)

	properties := p.parseFlowObjectProperties(false)
	p.consumePossible(token.COMMA)

	end := p.consumeExpected(token.RIGHT_BRACE)

	return &ast.FlowInexactObject{
		Start:      start,
		End:        end,
		Properties: properties,
	}
}

func (p *Parser) parseFlowExactObjectType() *ast.FlowExactObject {
	start := p.consumeExpected(token.TYPE_EXACT_OBJECT_START)

	properties := p.parseFlowObjectProperties(true)
	p.consumePossible(token.COMMA)

	end := p.consumeExpected(token.TYPE_EXACT_OBJECT_END)

	return &ast.FlowExactObject{
		Start:      start,
		End:        end,
		Properties: properties,
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
	}

	return nil
}

func (p *Parser) parseFlowTypeAnnotation() ast.FlowType {
	p.consumeExpected(token.COLON)

	return p.parseFlowType()
}
