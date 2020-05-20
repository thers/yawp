package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) isObjectPropertyNameStart() bool {
	if p.isAny(token.LEFT_BRACKET, token.STRING, token.NUMBER) {
		return true
	}

	return p.isIdentifierOrKeyword()
}

func (p *Parser) parseObjectPropertyName() ast.ObjectPropertyName {
	switch p.token {
	case token.LEFT_BRACKET:
		return p.parseObjectPropertyComputedName()
	case token.STRING:
		return p.parseString()
	case token.NUMBER:
		return p.parseNumber()
	default:
		id := p.parseIdentifierIncludingKeywords()

		if id == nil {
			return nil
		}

		return id
	}

	return nil
}

func (p *Parser) parseObjectPropertyComputedName() ast.ObjectPropertyName {
	loc := p.loc()
	p.consumeExpected(token.LEFT_BRACKET)

	propertyName := &ast.ComputedName{
		ExprNode:   p.exprNodeAt(loc),
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
	loc *file.Loc,
	async bool,
	generator bool,
	propertyName ast.ObjectPropertyName,
) *ast.ObjectPropertyValue {
	parameterList := p.parseFunctionParameterList()
	functionLiteral := &ast.FunctionLiteral{
		Node:       p.nodeAt(loc),
		Async:      async,
		Generator:  generator,
		Parameters: parameterList,
	}

	p.parseFunctionNodeBody(functionLiteral)

	return &ast.ObjectPropertyValue{
		PropertyName: propertyName,
		Value:        functionLiteral,
	}
}

func (p *Parser) parseObjectPropertyValue(loc *file.Loc, propertyName ast.ObjectPropertyName) *ast.ObjectPropertyValue {
	// Object function shorthand
	if p.is(token.LEFT_PARENTHESIS) {
		parameterList := p.parseFunctionParameterList()
		functionLiteral := &ast.FunctionLiteral{
			Node:       p.nodeAt(loc),
			Parameters: parameterList,
		}

		p.parseFunctionNodeBody(functionLiteral)

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
	loc := p.loc()

	// handle spreads first as the easiest variant
	if p.is(token.DOTDOTDOT) {
		p.next()

		return &ast.ObjectSpread{
			ExprNode:   p.exprNodeAt(loc),
			Expression: p.parseAssignmentExpression(),
		}
	}

	// generator method is also an easy variant
	if p.is(token.MULTIPLY) {
		p.next()

		return p.parseObjectPropertyMethodShorthand(loc, false, true, p.parseObjectPropertyName())
	}

	propertyName := p.parseObjectPropertyName()

	// easy variants when it's string/number/computedName
	if propertyName != nil {
		switch id := propertyName.(type) {
		// could be:
		// 'foo'(){} / 'foo': ...
		// 0(){}     / 0: ...
		// [...](){} / [...]: ...
		case *ast.StringLiteral, *ast.NumberLiteral, *ast.ComputedName:
			switch p.token {
			case token.LEFT_PARENTHESIS:
				return p.parseObjectPropertyMethodShorthand(loc, false, false, propertyName)
			case token.COLON:
				p.next()

				return &ast.ObjectPropertyValue{
					PropertyName: propertyName,
					Value:        p.parseAssignmentExpression(),
				}

			default:
				p.unexpectedToken()

				return nil
			}

		// could be just a shorthand for foo: foo
		// or a fucking rabbit hole of possibilities
		case *ast.Identifier:
			if p.isAny(token.COMMA, token.RIGHT_BRACE) || p.implicitSemicolon {
				return p.parseObjectPropertyFromShorthand(propertyName)
			}

			literal := id.Name
			possibleAsync := literal == "async"

			if literal == "get" || literal == "set" {
				loc = p.loc()
				accessor := id.Name

				// we have parsed set or get by now
				// if next is valid property identifier then it's an accessor
				if p.isObjectPropertyNameStart() {
					propertyName = p.parseObjectPropertyName()
					parameterList := p.parseFunctionParameterList()

					functionLiteral := &ast.FunctionLiteral{
						Node:       p.nodeAt(loc),
						Parameters: parameterList,
					}
					p.parseFunctionNodeBody(functionLiteral)

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

			loc = p.loc()

			// now that we have first two tokens we can start to disambiguate
			if possibleAsync {
				// `async blah` can only end with `() {}`, so a method shorthand
				if p.isIdentifierOrKeyword() {
					propertyName = p.parseIdentifierIncludingKeywords()

					return p.parseObjectPropertyMethodShorthand(loc, true, false, propertyName)
				}

				// `async []() {}`
				if p.is(token.LEFT_BRACKET) {
					propertyName = p.parseObjectPropertyComputedName()

					return p.parseObjectPropertyMethodShorthand(loc, true, false, propertyName)
				}
			}

			property := p.parseObjectPropertyValue(loc, propertyName)

			if property == nil {
				return p.parseObjectPropertyFromShorthand(propertyName)
			}

			return property
		}
	}

	p.unexpectedToken()
	p.next()

	return nil
}

func (p *Parser) parseObjectLiteral() *ast.ObjectLiteral {
	var value []ast.ObjectProperty

	loc := p.loc()

	p.consumeExpected(token.LEFT_BRACE)
	for p.until(token.RIGHT_BRACE) {
		property := p.parseObjectProperty()

		p.consumePossible(token.COMMA)

		value = append(value, property)
	}
	loc.End(p.consumeExpected(token.RIGHT_BRACE))

	return &ast.ObjectLiteral{
		ExprNode:   p.exprNodeAt(loc),
		Properties: value,
	}
}

func (p *Parser) maybeParseObjectBinding() (*ast.ObjectBinding, bool) {
	wasLeftHandSideAllowed := p.allowPatternBindingLeftHandSideExpressions
	p.allowPatternBindingLeftHandSideExpressions = true

	defer func() {
		_ = recover()

		p.allowPatternBindingLeftHandSideExpressions = wasLeftHandSideAllowed
	}()

	return p.parseObjectBinding(), true
}

func (p *Parser) parseObjectPatternAssignmentOrLiteral() ast.IExpr {
	snapshot := p.snapshot()
	restoreSymbolFlags := p.useSymbolFlags(ast.SWrite.Add(p.symbolFlags))

	objectBinding, success := p.maybeParseObjectBinding()

	restoreSymbolFlags()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.AssignmentExpression{
			Operator: token.ASSIGN,
			Left:     objectBinding,
			Right:    p.parseAssignmentExpression(),
		}
	}

	p.toSnapshot(snapshot)

	return p.parseObjectLiteral()
}
