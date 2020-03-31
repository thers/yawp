package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func areElementNamesEqual(a, b *ast.JSXElementName) bool {
	return a.StringName == b.StringName
}

func (p *Parser) parseString() *ast.StringLiteral {
	value, err := parseStringLiteral(p.literal[1 : len(p.literal)-1])
	if err != nil {
		p.error(p.idx, err.Error())
	}

	p.next()

	return &ast.StringLiteral{
		Start:   p.idx,
		Literal: p.literal,
		Value:   value,
	}
}

func (p *Parser) parseJSXElementName() *ast.JSXElementName {
	var left ast.Expression

	rootIdentifier := p.parseIdentifier()
	stringName := rootIdentifier.Name
	left = rootIdentifier


	if p.is(token.COLON) {
		// x:tag
		p.consumeExpected(token.COLON)

		name := p.parseIdentifierIncludingKeywords()

		return &ast.JSXElementName{
			Expression: &ast.JSXNamespacedName{
				Start:     rootIdentifier.Start,
				Namespace: rootIdentifier.Name,
				Name:      name.Name,
			},
			StringName: rootIdentifier.Name + ":" + name.Name,
		}
	}

	for {
		if p.is(token.PERIOD) {
			p.consumeExpected(token.PERIOD)

			member := p.parseIdentifierIncludingKeywords()

			if member == nil {
				p.unexpectedToken()
				p.next()
				return nil
			}

			stringName += "." + member.Name

			// left.member
			left = &ast.DotExpression{
				Left:       left,
				Identifier: member,
			}

			continue
		}

		return &ast.JSXElementName{
			Expression: left,
			StringName: stringName,
		}
	}
}

func (p *Parser) parseJSXChild() ast.JSXChild {
	switch p.token {
	case token.LESS:
		return p.parseJSXElement()
	case token.JSX_FRAGMENT_START:
		return p.parseJSXFragment()
	case token.LEFT_BRACE:
		p.consumeExpected(token.LEFT_BRACE)
		exp := p.parseAssignmentExpression()
		p.jsxTextParseFrom = int(p.consumeExpected(token.RIGHT_BRACE))

		return &ast.JSXChildExpression{
			Expression: exp,
		}
	}

	// parsing text
	start := p.jsxTextParseFrom
	text := p.str[start : p.chrOffset]
	for p.chr != '<' && p.chr != '{' && p.chr != -1 {
		text += string(p.chr)
		p.read()
	}

	p.jsxTextParseFrom = p.chrOffset

	p.next()

	return &ast.JSXText{
		Start: file.Idx(start),
		End:   p.idx,
		Text:  text,
	}
}

func (p *Parser) parseJSXElementAttributes() []ast.JSXAttribute {
	attrs := make([]ast.JSXAttribute, 0)

	// until we meet /> or > or EOF
	for !p.is(token.EOF) && !p.is(token.JSX_TAG_SELF_CLOSE) && !p.is(token.GREATER) {
		if p.is(token.IDENTIFIER) {
			attribute := &ast.JSXNamedAttribute{
				Name:  p.parseIdentifier(),
			}

			// attribute with initializer
			if p.is(token.ASSIGN) {
				p.consumeExpected(token.ASSIGN)

				if p.is(token.LEFT_BRACE) {
					p.consumeExpected(token.LEFT_BRACE)
					attribute.Value = p.parseAssignmentExpression()
					p.consumeExpected(token.RIGHT_BRACE)
				} else if p.is(token.STRING) {
					attribute.Value = p.parseString()
				} else {
					p.unexpectedToken()
					p.next()
				}
			} else {
				attribute.Value = &ast.BooleanLiteral{
					Start:   attribute.Name.Start,
					Literal: "true",
					Value:   true,
				}
			}

			attrs = append(attrs, attribute)
		} else if p.is(token.LEFT_BRACE) {
			// attributes spreading
			p.consumeExpected(token.LEFT_BRACE)
			start := p.consumeExpected(token.DOTDOTDOT)

			attrs = append(attrs, &ast.JSXSpreadAttribute{
				Start:      start,
				Expression: p.parseAssignmentExpression(),
			})

			p.consumeExpected(token.RIGHT_BRACE)
		} else {
			p.unexpectedToken()
			p.next()
		}
	}

	return attrs
}

func (p *Parser) parseJSXFragment() *ast.JSXFragment {
	start := p.consumeExpected(token.JSX_FRAGMENT_START)
	children := make([]ast.JSXChild, 0)

	for p.until(token.JSX_FRAGMENT_END) {
		children = append(children, p.parseJSXChild())
	}

	end := p.consumeExpected(token.JSX_FRAGMENT_END)
	// it's 3 chars wide token
	p.jsxTextParseFrom = int(end) + 2

	return &ast.JSXFragment{
		Start:    start,
		End:      end,
		Children: children,
	}
}

func (p *Parser) parseJSXElement() *ast.JSXElement {
	elm := &ast.JSXElement{
		Start:      p.consumeExpected(token.LESS),
		Name:       p.parseJSXElementName(),
		Children:   make([]ast.JSXChild, 0),
	}

	elm.Attributes = p.parseJSXElementAttributes()

	// self closing element />
	if p.is(token.JSX_TAG_SELF_CLOSE) {
		elm.End = p.consumeExpected(token.JSX_TAG_SELF_CLOSE)
		// it's a 2 chars wide token
		p.jsxTextParseFrom = int(elm.End) + 1

		return elm
	}

	// end of element >
	p.jsxTextParseFrom = int(p.consumeExpected(token.GREATER))

	// until </
	for p.until(token.JSX_TAG_CLOSE) {
		elm.Children = append(elm.Children, p.parseJSXChild())
	}

	p.consumeExpected(token.JSX_TAG_CLOSE)
	closeElementNamePos := p.idx
	closeElementName := p.parseJSXElementName()

	if !areElementNamesEqual(elm.Name, closeElementName) {
		p.error(closeElementNamePos, "Closing JSX element tag must be identical to the opening one")
		p.next()
		return nil
	}

	elm.End = p.consumeExpected(token.GREATER)
	p.jsxTextParseFrom = int(elm.End)

	return elm
}

func (p *Parser) parseJSXElementOrGenericArrowFunction() ast.Expression {
	partialState := p.getPartialState()
	errorsCount := len(p.errors)

	// first try to parse as flow type parameters
	typeParameters := p.parseFlowTypeParameters()

	if len(p.errors) == errorsCount && p.is(token.LEFT_PARENTHESIS) {
		parameters := p.parseFunctionParameterList()
		p.consumeExpected(token.ARROW)

		return &ast.ArrowFunctionExpression{
			Start:          parameters.Opening,
			Async:          true,
			TypeParameters: typeParameters,
			Parameters:     parameters.List,
			Body:           p.parseArrowFunctionBody(),
		}
	}

	// this isn't arrow function, so continue with jsx element
	p.restorePartialState(partialState)

	return p.parseJSXElement()
}
