package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func cmpElementNames(a, b ast.Expression) bool {
	if aIdentifier, ok := a.(*ast.Identifier); ok {
		if bIdentifier, ok := b.(*ast.Identifier); ok {
			return aIdentifier.Name == bIdentifier.Name
		}
	}

	if aDot, ok := a.(*ast.DotExpression); ok {
		if bDot, ok := b.(*ast.DotExpression); ok {
			return cmpElementNames(aDot.Left, bDot.Left) && aDot.Identifier.Name == bDot.Identifier.Name
		}
	}

	return false
}

func (p *Parser) parseJSXElementTag() ast.Expression {
	var tag ast.Expression

	left := p.parseIdentifier()

	if p.is(token.PERIOD) {
		p.consumeExpected(token.PERIOD)
		member := p.parseIdentifierIncludingKeywords()

		if member == nil {
			p.unexpectedToken()
			p.next()
			return nil
		}

		tag = &ast.DotExpression{
			Left:       left,
			Identifier: member,
		}
	} else {
		tag = left
	}

	return tag
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
		Tag:        p.parseJSXElementTag(),
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
	closeTagPos := p.idx
	closeTag := p.parseJSXElementTag()

	if !cmpElementNames(elm.Tag, closeTag) {
		p.error(closeTagPos, "Closing JSX element tag must be identical to the opening one")
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
