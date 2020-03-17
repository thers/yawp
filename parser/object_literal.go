package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseObjectPropertyFromShorthand(propertyName ast.ObjectPropertyName) ast.ObjectProperty {
	if propertyStringName, ok := propertyName.(*ast.Identifier); ok {
		if !matchIdentifier.MatchString(propertyStringName.Name) {
			return nil
		}

		if p.is(token.COMMA) || p.is(token.RIGHT_BRACE) {
			p.consumePossible(token.COMMA)

			return &ast.ObjectPropertyValue{
				PropertyName: propertyName,
				Value:        propertyStringName,
			}
		}
	}

	p.unexpectedToken()
	return nil
}

func (p *Parser) parseObjectPropertyValue(start file.Idx, propertyName ast.ObjectPropertyName) *ast.ObjectPropertyValue {
	// Object function shorthand
	if p.is(token.LEFT_PARENTHESIS) {
		parameterList := p.parseFunctionParameterList()
		functionLiteral := &ast.FunctionLiteral{
			Start:      start,
			Parameters: parameterList,
		}

		p.parseFunctionBlock(functionLiteral)
		functionLiteral.Source = p.slice(start, functionLiteral.Body.EndAt())

		return &ast.ObjectPropertyValue{
			PropertyName: propertyName,
			Value:        functionLiteral,
		}
	}

	if p.is(token.COLON) {
		p.consumeExpected(token.COLON)

		return &ast.ObjectPropertyValue{
			PropertyName: propertyName,
			Value:        p.parseAssignmentExpression(),
		}
	}

	return nil
}

func (p *Parser) parseObjectProperty() ast.ObjectProperty {
	var propertyName ast.ObjectPropertyName

	if p.is(token.IDENTIFIER) {
		shouldConsumeNext := true

		// Maybe setter or getter
		if p.literal == "get" || p.literal == "set" {
			start, name := p.idx, p.literal
			p.next()

			// Setter or getter
			if p.is(token.IDENTIFIER) {
				propertyName = &ast.Identifier{
					Start: start,
					Name:  name,
				}

				p.next()
				parameterList := p.parseFunctionParameterList()

				functionLiteral := &ast.FunctionLiteral{
					Start:      start,
					Parameters: parameterList,
				}
				p.parseFunctionBlock(functionLiteral)

				if name == "set" {
					return &ast.ObjectPropertySetter{
						PropertyName: propertyName,
						Setter:       functionLiteral,
					}
				} else {
					return &ast.ObjectPropertyGetter{
						PropertyName: propertyName,
						Getter:       functionLiteral,
					}
				}
			}

			shouldConsumeNext = false
		}

		name, start := p.literal, p.idx

		propertyName = &ast.Identifier{
			Start: start,
			Name:  name,
		}

		if shouldConsumeNext {
			p.next()
		}

		property := p.parseObjectPropertyValue(start, propertyName)

		if property == nil {
			return p.parseObjectPropertyFromShorthand(propertyName)
		}

		return property
	} else if p.is(token.NUMBER) {
		name, start := p.literal, p.idx
		_, err := parseNumberLiteral(p.literal)

		if err != nil {
			name = ""
			p.error(p.idx, err.Error())
		}

		propertyName = &ast.Identifier{
			Start: start,
			Name:  name,
		}

		p.consumeExpected(token.COLON)

		return &ast.ObjectPropertyValue{
			PropertyName: propertyName,
			Value:        p.parseAssignmentExpression(),
		}
	} else if p.is(token.STRING) {
		start := p.idx
		name, err := parseStringLiteral(p.literal[1 : len(p.literal)-1])
		if err != nil {
			p.error(p.idx, err.Error())
		}

		propertyName = &ast.Identifier{
			Start: start,
			Name:  name,
		}

		p.consumeExpected(token.COLON)

		return &ast.ObjectPropertyValue{
			PropertyName: propertyName,
			Value:        p.parseAssignmentExpression(),
		}
	} else if p.is(token.LEFT_BRACKET) {
		// computed property name
		start := p.idx
		p.consumeExpected(token.LEFT_BRACKET)

		propertyName = &ast.ComputedName{
			Expression: p.parseAssignmentExpression(),
		}

		p.consumeExpected(token.RIGHT_BRACKET)

		property := p.parseObjectPropertyValue(start, propertyName)

		if property != nil {
			return property
		}
	} else if p.is(token.DOTDOTDOT) {
		p.consumeExpected(token.DOTDOTDOT)

		return &ast.ObjectSpread{
			Expression: p.parseAssignmentExpression(),
		}
	}

	p.unexpectedToken()
	p.next()

	return nil
}

func (p *Parser) parseObjectLiteral() *ast.ObjectLiteral {
	var value []ast.ObjectProperty

	leftBrace := p.consumeExpected(token.LEFT_BRACE)
	for p.until(token.RIGHT_BRACE) {
		property := p.parseObjectProperty()

		p.consumePossible(token.COMMA)

		value = append(value, property)
	}
	rightBrace := p.consumeExpected(token.RIGHT_BRACE)

	return &ast.ObjectLiteral{
		LeftBrace:  leftBrace,
		RightBrace: rightBrace,
		Properties: value,
	}
}

func (p *Parser) maybeParseObjectBinding() (*ast.ObjectBinding, bool) {
	wasPatternBindingMode := p.patternBindingMode
	p.patternBindingMode = true

	defer func() {
		p.patternBindingMode = wasPatternBindingMode
	}()

	return p.parseObjectBinding(), p.patternBindingMode
}

func (p *Parser) parseObjectLiteralOrObjectPatternBinding() ast.Expression {
	start := p.idx
	partialState := p.getPartialState()

	objectBinding, success := p.maybeParseObjectBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.VariableBinding{
			Start:       start,
			Binder:      objectBinding,
			Initializer: p.parseAssignmentExpression(),
		}
	}

	p.restorePartialState(partialState)

	return p.parseObjectLiteral()
}
