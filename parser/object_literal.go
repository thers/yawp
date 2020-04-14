package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseObjectPropertyComputedName() ast.ObjectPropertyName {
	p.consumeExpected(token.LEFT_BRACKET)

	propertyName := &ast.ComputedName{
		Expression: p.parseAssignmentExpression(),
	}

	p.consumeExpected(token.RIGHT_BRACKET)

	return propertyName
}

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

func (p *Parser) parseObjectPropertyMethodShorthand(
	start file.Loc,
	async bool,
	propertyName ast.ObjectPropertyName,
) *ast.ObjectPropertyValue {
	parameterList := p.parseFunctionParameterList()
	functionLiteral := &ast.FunctionLiteral{
		Start:      start,
		Async:      async,
		Parameters: parameterList,
	}

	p.parseFunctionBlock(functionLiteral)

	return &ast.ObjectPropertyValue{
		PropertyName: propertyName,
		Value:        functionLiteral,
	}
}

func (p *Parser) parseObjectPropertyValue(start file.Loc, propertyName ast.ObjectPropertyName) *ast.ObjectPropertyValue {
	// Object function shorthand
	if p.is(token.LEFT_PARENTHESIS) {
		parameterList := p.parseFunctionParameterList()
		functionLiteral := &ast.FunctionLiteral{
			Start:      start,
			Parameters: parameterList,
		}

		p.parseFunctionBlock(functionLiteral)

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

/*
Modern JS standard have brought to the light ugly creatures
called contextual keywords.

These rare species are very special kinds of snowflakes.
Normal keyword is a reserved word,
that means that you can not use it as a variable or property name.

But there was `class: foo` all the time, you may have wonder now.
Yes, but it is a cake for parser to support it.

But there is one keyword, that is still ambiguous: async

consider these valid examples:

`foo = { async }` - a shorthand for `{ async: async }` where second async is some var
`foo = { async: 0 }` - just using async as a property name, pretty much like with class, ok
`foo = { async async() {} }` - well, lol, async method with name async
`foo = { async await() {} }` - not really related, but still funny
`foo = { set async(v) {} }` - not really different from `async async() {}`
`foo = { async set(v) {} }` - ffs

*/
func (p *Parser) parseObjectProperty() ast.ObjectProperty {
	var propertyName ast.ObjectPropertyName

	// start with easy variants
	switch p.token {
	case token.NUMBER:
		name, start := p.literal, p.loc
		_, err := parseNumberLiteral(p.literal)
		p.next()

		if err != nil {
			name = ""
			p.error(p.loc, err.Error())
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

	case token.STRING:
		start := p.loc
		name, err := parseStringLiteral(p.literal[1 : len(p.literal)-1])
		p.next()

		if err != nil {
			p.error(p.loc, err.Error())
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

	case token.LEFT_BRACKET:
		// computed property name
		start := p.loc
		propertyName = p.parseObjectPropertyComputedName()

		property := p.parseObjectPropertyValue(start, propertyName)
		if property != nil {
			return property
		}

	case token.DOTDOTDOT:
		p.consumeExpected(token.DOTDOTDOT)

		return &ast.ObjectSpread{
			Expression: p.parseAssignmentExpression(),
		}
	}

	// now that easy variants are all gone
	// let the show of identifiers and keywords begin
	if p.isIdentifierOrKeyword() {
		possibleAsync := p.is(token.ASYNC)
		shouldConsumeNext := true

		if p.literal == "get" || p.literal == "set" {
			start, accessor := p.loc, p.literal

			p.next()
			shouldConsumeNext = false

			// we have parsed set or get by now
			// if next is valid property identifier then it's an accessor
			if p.isIdentifierOrKeyword() {
				propertyName = p.parseIdentifierIncludingKeywords()
				parameterList := p.parseFunctionParameterList()

				functionLiteral := &ast.FunctionLiteral{
					Start:      start,
					Parameters: parameterList,
				}
				p.parseFunctionBlock(functionLiteral)

				if accessor == "set" {
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
		}

		start := p.loc
		propertyName = p.currentIdentifier()

		if shouldConsumeNext {
			p.next()
		}
		
		// now that we have first two tokens we can start to disambiguate
		if possibleAsync {
			// `async blah` can only end with `() {}`, so a method shorthand
			if p.isIdentifierOrKeyword() {
				propertyName = p.parseIdentifierIncludingKeywords()

				return p.parseObjectPropertyMethodShorthand(start, true, propertyName)
			}

			// `async []() {}`
			if p.is(token.LEFT_BRACKET) {
				propertyName = p.parseObjectPropertyComputedName()

				return p.parseObjectPropertyMethodShorthand(start, true, propertyName)
			}
		}

		property := p.parseObjectPropertyValue(start, propertyName)

		if property == nil {
			return p.parseObjectPropertyFromShorthand(propertyName)
		}

		return property
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
	defer func() {
		err := recover()
		if err != nil {
			return
		}
	}()

	return p.parseObjectBinding(), true
}

func (p *Parser) parseObjectLiteralOrObjectPatternBinding() ast.Expression {
	start := p.loc
	partialState := p.captureState()

	objectBinding, success := p.maybeParseObjectBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.VariableBinding{
			Start:       start,
			Binder:      objectBinding,
			Initializer: p.parseAssignmentExpression(),
		}
	}

	p.rewindStateTo(partialState)

	return p.parseObjectLiteral()
}
