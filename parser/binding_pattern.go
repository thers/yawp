package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseBindingDefaultValue() ast.Expression {
	if p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return p.parseAssignmentExpression()
	}

	return nil
}

func (p *Parser) parseBinder() ast.PatternBinder {
	switch p.token {
	case token.IDENTIFIER:
		return &ast.IdentifierBinder{
			Name: p.parseIdentifier(),
		}
	case token.LEFT_BRACKET:
		return p.parseArrayBinding()
	case token.LEFT_BRACE:
		return p.parseObjectBinding()
	default:
		p.unexpectedToken()
		p.next()
	}

	return nil
}

func (p *Parser) parseObjectBinding() *ast.ObjectBinding {
	start := p.consumeExpected(token.LEFT_BRACE)

	pattern := &ast.ObjectBinding{
		Start: start,
		List:  make([]*ast.ObjectPropertyBinder, 0),
	}

	for p.until(token.RIGHT_BRACE) {
		property := &ast.ObjectPropertyBinder{
			Property:     nil,
			PropertyName: nil,
			DefaultValue: nil,
		}

		propertyIdx := p.idx
		propertyName := p.parseIdentifierIncludingKeywords()

		if propertyName == nil {
			p.unexpectedToken()
			p.next()
			continue
		}

		property.PropertyName = propertyName

		if p.is(token.COLON) {
			p.consumeExpected(token.COLON)

			property.Property = p.parseBinder()
		} else {
			_, isKeyword := token.IsKeyword(propertyName.Name)

			if isKeyword {
				p.error(propertyIdx, "Unexpected token")
				p.next()
				continue
			} else {
				property.Property = &ast.IdentifierBinder{
					Name: property.PropertyName,
				}
			}
		}

		property.DefaultValue = p.parseBindingDefaultValue()
		p.consumePossible(token.COMMA)

		pattern.List = append(pattern.List, property)
	}

	pattern.End = p.consumeExpected(token.RIGHT_BRACKET)

	return pattern
}

func (p *Parser) parseArrayBinding() *ast.ArrayBinding {
	start := p.consumeExpected(token.LEFT_BRACKET)

	pattern := &ast.ArrayBinding{
		Start: start,
		List:  make([]*ast.ArrayItemBinder, 0),
	}

	for p.until(token.RIGHT_BRACKET) {
		item := &ast.ArrayItemBinder{
			Item:         p.parseBinder(),
			DefaultValue: p.parseBindingDefaultValue(),
		}

		p.consumePossible(token.COMMA)

		pattern.List = append(pattern.List, item)
	}

	pattern.End = p.consumeExpected(token.RIGHT_BRACKET)

	return pattern
}
