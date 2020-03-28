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

func (p *Parser) parseFlowInterfaceMethod() *ast.FlowInterfaceMethod {
	method := &ast.FlowInterfaceMethod{
		Start: p.idx,
	}

	if p.is(token.LESS) {
		method.TypeParameters = p.parseFlowTypeParameters()
	}

	method.Parameters, method.ReturnType = p.parseFlowInterfaceMethodParts()

	return method
}

func (p *Parser) parseFlowInterfaceBodyStatement() ast.FlowInterfaceBodyStatement {
	start := p.idx
	covariant, contravariant := p.parseFlowTypeVariance()
	isTypeParameters := p.token == token.LESS

	if isTypeParameters || p.is(token.LEFT_PARENTHESIS) {
		// only [[call]] can start right with type parameters
		return p.parseFlowInterfaceMethod()
	}

	if p.isIdentifierOrKeyword() {
		identifier := p.parseFlowTypeIdentifierIncludingKeywords()

		if identifier == nil {
			p.unexpectedToken()
			p.next()
			return nil
		}

		if p.isAny(token.LEFT_PARENTHESIS, token.LESS) {
			// it's a method
			method := p.parseFlowInterfaceMethod()
			method.Start = identifier.Start
			method.Name = identifier

			return method
		}

		return p.parseFlowNamedObjectPropertyRemainder(
			start,
			covariant,
			contravariant,
			identifier.Name,
		)
	}

	p.unexpectedToken()
	p.next()

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
