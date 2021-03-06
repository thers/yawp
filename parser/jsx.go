package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func areElementNamesEqual(a, b *ast.JSXElementName) bool {
	return a.StringName == b.StringName
}

func (p *Parser) parseJSXElementName() *ast.JSXElementName {
	var left ast.IExpr

	rootIdentifier := p.parseIdentifier()
	stringName := rootIdentifier.Name
	left = rootIdentifier

	if p.is(token.COLON) {
		// x:tag
		p.consumeExpected(token.COLON)

		name := p.parseIdentifierIncludingKeywords()
		loc := rootIdentifier.Loc.Copy()
		loc.To += file.Idx(len(name.Name))

		return &ast.JSXElementName{
			Expression: &ast.JSXNamespacedName{
				ExprNode:  p.exprNodeAt(loc),
				Namespace: rootIdentifier.Name,
				Name:      name.Name,
			},
			StringName: rootIdentifier.Name + ":" + name.Name,
		}
	}

	p.symbol(rootIdentifier, ast.SRead, ast.SRUnknown)

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
			left = &ast.MemberExpression{
				Left:  left,
				Right: member,
				Kind:  ast.MKObject,
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

		var exp ast.IExpr

		if !p.is(token.RIGHT_BRACE) {
			exp = p.parseAssignmentExpression()
		}

		p.jsxTextParseFrom = int(p.tokenOffset)
		p.consumeExpected(token.RIGHT_BRACE)

		return &ast.JSXChildExpression{
			IExpr: exp,
		}
	}

	// parsing text
	loc := p.loc()
	loc.From = file.Idx(p.jsxTextParseFrom)

	for p.chr != '<' && p.chr != '{' && p.chr != -1 {
		p.read()
	}

	loc.To = file.Idx(p.chrOffset)
	text := p.src[p.jsxTextParseFrom:p.chrOffset]

	p.jsxTextParseFrom = p.chrOffset

	p.next()

	return &ast.JSXText{
		Node: p.nodeAt(loc),
		Text: text,
	}
}

func (p *Parser) parseJSXElementAttributes() []ast.JSXAttribute {
	attrs := make([]ast.JSXAttribute, 0)

	// until we meet /> or > or EOF
	for !p.is(token.EOF) && !p.is(token.JSX_TAG_SELF_CLOSE) && !p.is(token.GREATER) {
		if p.is(token.IDENTIFIER) {
			attribute := &ast.JSXNamedAttribute{
				Name: p.parseIdentifier(),
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
					ExprNode: p.exprNodeAt(attribute.Name.GetLoc()),
					Literal:  "true",
				}
			}

			attrs = append(attrs, attribute)
		} else if p.is(token.LEFT_BRACE) {
			// attributes spreading
			p.consumeExpected(token.LEFT_BRACE)
			start := p.tokenOffset
			p.consumeExpected(token.DOTDOTDOT)

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
	loc := p.loc()

	p.consumeExpected(token.JSX_FRAGMENT_START)
	children := make([]ast.JSXChild, 0)

	for p.until(token.JSX_FRAGMENT_END) {
		children = append(children, p.parseJSXChild())
	}

	loc.End(p.consumeExpected(token.JSX_FRAGMENT_END))
	// it's 3 chars wide token
	p.jsxTextParseFrom = int(loc.To) + 2

	return &ast.JSXFragment{
		ExprNode: p.exprNodeAt(loc),
		Children: children,
	}
}

func (p *Parser) parseJSXElement() *ast.JSXElement {
	loc := p.loc()
	p.consumeExpected(token.LESS)

	elm := &ast.JSXElement{
		ExprNode: p.exprNodeAt(loc),
		Name:     p.parseJSXElementName(),
		Children: make([]ast.JSXChild, 0),
	}

	// if we're looking for generic arrow fn
	if p.is(token.COMMA) && p.genericTypeParametersMode {
		p.unexpectedToken()
		return nil
	}

	elm.Attributes = p.parseJSXElementAttributes()

	// self closing element />
	if p.is(token.JSX_TAG_SELF_CLOSE) {
		elm.Loc.End(p.consumeExpected(token.JSX_TAG_SELF_CLOSE))
		// it's a 2 chars wide token
		p.jsxTextParseFrom = int(elm.Loc.To) + 1

		return elm
	}

	// end of element >
	p.jsxTextParseFrom = int(p.tokenOffset)
	p.consumeExpected(token.GREATER)

	// until </
	for p.until(token.JSX_TAG_CLOSE) {
		elm.Children = append(elm.Children, p.parseJSXChild())
	}

	p.consumeExpected(token.JSX_TAG_CLOSE)
	closeElementNameLoc := p.loc()
	closeElementName := p.parseJSXElementName()

	if !areElementNamesEqual(elm.Name, closeElementName) {
		p.error(closeElementNameLoc, "Closing JSX element tag must be identical to the opening one")
		p.next()
		return nil
	}

	elm.Loc.End(p.consumeExpected(token.GREATER))
	p.jsxTextParseFrom = int(elm.Loc.To)

	return elm
}

func (p *Parser) maybeParseJSXElement() ast.IExpr {
	defer func() { _ = recover() }()

	p.genericTypeParametersMode = true
	jsx := p.parseJSXElement()
	p.genericTypeParametersMode = false

	return jsx
}

func (p *Parser) parseJSXElementOrGenericArrowFunction() ast.IExpr {
	snapshot := p.snapshot()

	// first try to parse as jsx
	jsx := p.maybeParseJSXElement()

	if jsx != nil {
		return jsx
	}

	// now we can safely assume we're in arrow function
	p.toSnapshot(snapshot)

	return p.parseParametrizedArrowFunction()
}
