package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseBindingDefaultValue() ast.IExpr {
	if p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return p.parseAssignmentExpression()
	}

	return nil
}

func (p *Parser) parseBindingMemberExpressionOrIdentifier() ast.PatternBinder {
	binding := p.parseMemberExpressionOrIdentifier()

	switch bnd := binding.(type) {
	case *ast.Identifier:
		return &ast.IdentifierBinder{Id: bnd}
	case *ast.MemberExpression:
		return &ast.ExpressionBinder{Expression: bnd}
	default:
		panic("Impossible case")
	}
}

func (p *Parser) parseBinder() ast.PatternBinder {
	switch p.token {
	case token.LEFT_BRACKET:
		return p.parseArrayBinding()
	case token.LEFT_BRACE:
		return p.parseObjectBinding()
	}

	if p.allowPatternBindingLeftHandSideExpressions {
		return p.parseBindingMemberExpressionOrIdentifier()
	} else {
		if p.is(token.IDENTIFIER) {
			return &ast.IdentifierBinder{
				Id: p.symbol(p.parseIdentifier(), ast.SymbolWrite.Add(p.symbolFlags), ast.SRUnknown),
			}
		}
	}

	p.unexpectedToken()

	return nil
}

func (p *Parser) parseObjectBinding() *ast.ObjectBinding {
	loc := p.loc()
	p.consumeExpected(token.LEFT_BRACE)

	pattern := &ast.ObjectBinding{
		ExprNode: p.exprNodeAt(loc),
		List:     make([]ast.PatternBinder, 0),
	}

	boundProperties := make([]ast.ObjectPropertyName, 0)

	for p.until(token.RIGHT_BRACE) {
		// { ...a }
		if p.is(token.DOTDOTDOT) {
			p.consumeExpected(token.DOTDOTDOT)

			pattern.List = append(pattern.List, &ast.ObjectRestBinder{
				Binder:         p.parseBinder(),
				OmitProperties: boundProperties,
			})
			break
		}

		propertyLoc := p.loc()
		property := &ast.ObjectPropertyBinder{
			PropertyName: p.parseObjectPropertyName(),
		}

		if property.PropertyName == nil {
			p.unexpectedToken()

			return nil
		}

		boundProperties = append(boundProperties, property.PropertyName)

		if p.is(token.COLON) {
			p.consumeExpected(token.COLON)

			property.Binder = p.parseBinder()
		} else {
			switch propertyId := property.PropertyName.(type) {
			case *ast.Identifier:
				if _, isKeyword := token.IsKeyword(propertyId.Name); isKeyword {
					p.unexpectedTokenAt(propertyLoc)

					return nil
				}

				property.Binder = &ast.IdentifierBinder{
					Id: p.symbol(propertyId.Copy(), ast.SymbolWrite.Add(p.symbolFlags), ast.SRUnknown),
				}
			default:
				p.error(propertyLoc, "Can not use computed property name without pattern binding")

				return nil
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
		ExprNode: p.exprNodeAt(loc),
		List:     make([]ast.PatternBinder, 0),
	}

	itemIndex := 0

	for p.until(token.RIGHT_BRACKET) {
		// [...a]
		if p.is(token.DOTDOTDOT) {
			p.consumeExpected(token.DOTDOTDOT)

			pattern.List = append(pattern.List, &ast.ArrayRestBinder{
				Binder:    p.parseBinder(),
				FromIndex: itemIndex,
			})
			break
		}

		if p.is(token.COMMA) {
			itemIndex++
			p.next()
			continue
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

func (p *Parser) parseArrayBindingAllowLHS() *ast.ArrayBinding {
	old := p.allowPatternBindingLeftHandSideExpressions
	p.allowPatternBindingLeftHandSideExpressions = true

	bnd := p.parseArrayBinding()
	p.allowPatternBindingLeftHandSideExpressions = old

	return bnd
}

func (p *Parser) parseObjectBindingAllowLHS() *ast.ObjectBinding {
	old := p.allowPatternBindingLeftHandSideExpressions
	p.allowPatternBindingLeftHandSideExpressions = true

	bnd := p.parseObjectBinding()
	p.allowPatternBindingLeftHandSideExpressions = old

	return bnd
}
