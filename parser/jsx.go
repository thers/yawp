package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func areElementNamesEqual(a, b *ast.JSXElementName) bool {
	return a.StringName == b.StringName
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

func (p *Parser) parseJSXNode() ast.JSXNode {
	switch p.token {
	case token.LESS:
		return p.parseJSXElement()
	case token.JSX_FRAGMENT_START:
		return p.parseJSXFragment()
	}

	return nil
}

func (p *Parser) parseJSXFragment() *ast.JSXFragment {
	start := p.consumeExpected(token.JSX_FRAGMENT_START)
	children := make([]ast.JSXNode, 0)

	for p.until(token.JSX_FRAGMENT_END) {
		children = append(children, p.parseJSXNode())
	}

	end := p.consumeExpected(token.JSX_FRAGMENT_END)

	return &ast.JSXFragment{
		Start:    start,
		End:      end,
		Children: children,
	}
}

func (p *Parser) parseJSXElement() *ast.JSXElement {
	elm := &ast.JSXElement{
		Start:      p.consumeExpected(token.LESS),
		End:        0,
		Name:       p.parseJSXElementName(),
		Attributes: nil,
		Children:   make([]ast.JSXNode, 0),
	}

	// TODO: Attributes parsing

	if p.is(token.JSX_TAG_SELF_CLOSE) {
		elm.End = p.consumeExpected(token.JSX_TAG_SELF_CLOSE)

		return elm
	}

	p.consumeExpected(token.GREATER)

	for p.until(token.JSX_TAG_CLOSE) {
		elm.Children = append(elm.Children, p.parseJSXNode())
	}

	p.consumeExpected(token.JSX_TAG_CLOSE)
	closeElementNamePos := p.idx
	closeElementName := p.parseJSXElementName()

	if !areElementNamesEqual(elm.Name, closeElementName) {
		p.error(closeElementNamePos, "Closing JSX element tag must be identical to the opening one")
		p.next()
		return nil
	}

	p.consumeExpected(token.GREATER)

	return elm
}

func (p *Parser) parseJSX() ast.Expression {
	if p.is(token.JSX_FRAGMENT_START) {
		return p.parseJSXFragment()
	} else if p.is(token.LESS) {
		return p.parseJSXElement()
	}

	return &ast.BadExpression{
		From: p.idx,
		To:   0,
	}
}
