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
			Id: p.parseIdentifier(),
		}
	case token.LEFT_BRACKET:
		return p.parseArrayBinding()
	case token.LEFT_BRACE:
		return p.parseObjectBinding()
	default:
		p.unexpectedToken()

		return nil
	}

	return nil
}

func (p *Parser) parseObjectBinding() *ast.ObjectBinding {
	loc := p.loc()
	p.consumeExpected(token.LEFT_BRACE)

	pattern := &ast.ObjectBinding{
		Loc:  loc,
		List: make([]ast.PatternBinder, 0),
	}

	boundProperties := make([]*ast.Identifier, 0)

	for p.until(token.RIGHT_BRACE) {
		// { ...a }
		if p.is(token.DOTDOTDOT) {
			p.consumeExpected(token.DOTDOTDOT)

			pattern.List = append(pattern.List, &ast.ObjectRestBinder{
				Id:             p.parseIdentifier(),
				OmitProperties: boundProperties,
			})
			break
		}

		property := &ast.ObjectPropertyBinder{
			Binder:       nil,
			Id:           nil,
			DefaultValue: nil,
		}

		propertyLoc := p.loc()
		propertyName := p.parseIdentifierIncludingKeywords()

		if propertyName == nil {
			p.unexpectedToken()

			return nil
		}

		property.Id = propertyName
		boundProperties = append(boundProperties, property.Id)

		if p.is(token.COLON) {
			p.consumeExpected(token.COLON)

			property.Binder = p.parseBinder()
		} else {
			_, isKeyword := token.IsKeyword(propertyName.Name)

			if isKeyword {
				p.unexpectedTokenAt(propertyLoc)

				return nil
			} else {
				property.Binder = &ast.IdentifierBinder{
					Id: property.Id,
				}
			}
		}

		property.DefaultValue = p.parseBindingDefaultValue()
		p.consumePossible(token.COMMA)

		pattern.List = append(pattern.List, property)
	}

	pattern.Loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return pattern
}

func (p *Parser) parseArrayBinding() *ast.ArrayBinding {
	loc := p.loc()
	p.consumeExpected(token.LEFT_BRACKET)

	pattern := &ast.ArrayBinding{
		Loc:  loc,
		List: make([]ast.PatternBinder, 0),
	}

	itemIndex := 0

	for p.until(token.RIGHT_BRACKET) {
		// [...a]
		if p.is(token.DOTDOTDOT) {
			p.consumeExpected(token.DOTDOTDOT)

			pattern.List = append(pattern.List, &ast.ArrayRestBinder{
				Id:        p.parseIdentifier(),
				FromIndex: itemIndex,
			})
			break
		}

		item := &ast.ArrayItemBinder{
			Binder:       p.parseBinder(),
			Index:        itemIndex,
			DefaultValue: p.parseBindingDefaultValue(),
		}

		itemIndex++
		p.consumePossible(token.COMMA)

		pattern.List = append(pattern.List, item)
	}

	pattern.Loc.End(p.consumeExpected(token.RIGHT_BRACKET))

	return pattern
}
