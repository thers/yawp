package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseFlowInterfaceMethodParts() ([]interface{}, ast.FlowType) {
	p.consumeExpected(token.LEFT_PARENTHESIS)
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	params := make([]interface{}, 0)
	retType := p.parseFlowTypeAnnotation()

	return params, retType
}

func (p *Parser) parseFlowInterfaceBodyStatement() ast.FlowInterfaceBodyStatement {
	isVariance := p.isAny(token.MINUS, token.PLUS)
	isTypeParameters := p.is(token.LESS)

	if isVariance || p.isIdentifierOrKeyword() {
		return p.parseFlowNamedObjectProperty()
	}

	if isTypeParameters || p.is(token.LEFT_PARENTHESIS) {
		method := &ast.FlowInterfaceMethod{
			Start: p.idx,
		}

		if isTypeParameters {
			method.TypeParameters = p.parseFlowTypeParameters()
		}

		method.Parameters, method.ReturnType = p.parseFlowInterfaceMethodParts()

		return method
	}

	return nil
}

func (p *Parser) parseFlowInterfaceStatement() *ast.FlowInterfaceStatement {
	start := p.consumeExpected(token.INTERFACE)
	name := p.parseFlowTypeIdentifier()
	var typeParameters []*ast.FlowTypeParameter

	if p.is(token.LESS) {
		typeParameters = p.parseFlowTypeParameters()
	}

	p.consumeExpected(token.LEFT_BRACE)

	stmts := make([]ast.FlowInterfaceBodyStatement, 0)

	for p.until(token.RIGHT_BRACE) {
		stmts = append(stmts, p.parseFlowInterfaceBodyStatement())

		p.implicitSemicolon = true

		if p.is(token.COMMA) {
			p.next()
		} else {
			p.optionalSemicolon()
		}
	}

	end := p.consumeExpected(token.RIGHT_BRACE)

	return &ast.FlowInterfaceStatement{
		Start:          start,
		End:            end,
		Name:           name,
		TypeParameters: typeParameters,
		Body:           stmts,
	}
}
