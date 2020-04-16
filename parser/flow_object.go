package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseFlowNamedObjectPropertyRemainder(
	loc *file.Loc,
	covariant,
	contravariant bool,
	name string,
) *ast.FlowNamedObjectProperty {
	prop := &ast.FlowNamedObjectProperty{
		Loc:           loc,
		Name:          name,
		Optional:      false,
		Covariant:     covariant,
		Contravariant: contravariant,
	}

	if p.is(token.QUESTION_MARK) {
		p.next()
		prop.Optional = true
	}

	prop.Value = p.parseFlowTypeAnnotation()

	return prop
}

func (p *Parser) parseFlowTypeVariance() (covariant, contravariant bool) {
	if p.isAny(token.MINUS, token.PLUS) {
		covariant = p.token == token.PLUS
		contravariant = p.token == token.MINUS

		p.next()
	}

	return
}

func (p *Parser) parseFlowNamedObjectProperty() *ast.FlowNamedObjectProperty {
	loc := p.loc()

	covariant, contravariant := p.parseFlowTypeVariance()
	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		p.unexpectedToken()
		p.next()

		return nil
	}

	return p.parseFlowNamedObjectPropertyRemainder(
		loc,
		covariant,
		contravariant,
		identifier.Name,
	)
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
		loc := p.loc()

		if p.isIdentifierOrKeyword() || p.isAny(token.PLUS, token.MINUS) {
			prop := p.parseFlowNamedObjectProperty()

			if prop == nil {
				goto Next
			}

			props = append(props, prop)
		} else if p.is(token.DOTDOTDOT) {
			p.next()

			// inexact object specifier
			if p.is(token.COMMA) || p.is(terminator) {
				if exact {
					// explicit inexact disallowed in exact
					p.error(loc, "Explicit inexact syntax is not allowed in exact object type")
					break
				}

				props = append(props, &ast.FlowInexactSpecifierProperty{
					Loc: loc,
				})
				break
			} else {
				props = append(props, &ast.FlowSpreadObjectProperty{
					Loc:      loc,
					FlowType: p.parseFlowType(),
				})
			}
		} else if p.is(token.LEFT_BRACKET) {
			p.next()

			prop := &ast.FlowIndexerObjectProperty{
				Loc:     loc,
				KeyName: "",
				KeyType: nil,
				Value:   nil,
			}

			if p.isIdentifierOrKeyword() {
				isKeyword := p.tokenIsKeyword
				identifier := p.parseIdentifierIncludingKeywords()

				if identifier == nil {
					p.unexpectedToken()
					p.next()
					goto Next
				}

				if p.is(token.COLON) {
					if isKeyword {
						// keywords aren't legal identifiers
						p.error(identifier.GetLoc(), "Cannot use keyword as type identifier")
						p.next()
						goto Next
					}

					prop.KeyName = identifier.Name
					prop.KeyType = p.parseFlowTypeAnnotation()
				} else if p.is(token.RIGHT_BRACKET) {
					prop.KeyType = &ast.FlowIdentifier{
						Loc:  identifier.Loc,
						Name: identifier.Name,
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
		if !p.is(terminator) {
			if p.is(token.SEMICOLON) {
				p.next()
			} else {
				p.consumeExpected(token.COMMA)
			}
		}
	}

	return props
}

func (p *Parser) parseFlowInexactObjectType() *ast.FlowInexactObject {
	loc := p.loc()
	p.consumeExpected(token.LEFT_BRACE)

	properties := p.parseFlowObjectProperties(false)
	p.consumePossible(token.COMMA)

	loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return &ast.FlowInexactObject{
		Loc:        loc,
		Properties: properties,
	}
}

func (p *Parser) parseFlowExactObjectType() *ast.FlowExactObject {
	loc := p.loc()
	p.consumeExpected(token.TYPE_EXACT_OBJECT_START)

	properties := p.parseFlowObjectProperties(true)
	p.consumePossible(token.COMMA)

	loc.End(p.consumeExpected(token.TYPE_EXACT_OBJECT_END))

	return &ast.FlowExactObject{
		Loc:        loc,
		Properties: properties,
	}
}
