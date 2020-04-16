package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) isVariableStatementStart() bool {
	return p.is(token.VAR) || p.is(token.CONST) || p.is(token.LET)
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	loc := p.loc()
	kind := p.token

	if !p.isVariableStatementStart() {
		p.unexpectedToken()
	}

	p.next()
	list := p.parseVariableDeclarationList(loc, kind)

	p.optionalSemicolon()

	return &ast.VariableStatement{
		Loc:  loc,
		Kind: kind,
		List: list,
	}
}

func (p *Parser) parseVariableDeclaration(declarationList *[]*ast.VariableExpression, kind token.Token) ast.Expression {
	if p.is(token.LEFT_BRACKET) || p.is(token.LEFT_BRACE) {
		loc := p.loc()

		var binder ast.PatternBinder

		if p.is(token.LEFT_BRACKET) {
			binder = p.parseArrayBinding()
		} else {
			binder = p.parseObjectBinding()
		}

		bnd := &ast.VariableBinding{
			Loc:    loc,
			Binder: binder,
		}

		if p.is(token.ASSIGN) {
			p.consumeExpected(token.ASSIGN)

			bnd.Initializer = p.parseAssignmentExpression()
		}

		return bnd
	}

	if !p.is(token.IDENTIFIER) {
		p.unexpectedToken()

		return nil
	}

	loc := p.loc()
	literal := p.literal

	p.next()
	node := &ast.VariableExpression{
		Loc:  loc,
		Kind: kind,
		Name: literal,
	}

	if declarationList != nil {
		*declarationList = append(*declarationList, node)
	}

	// feat(type)
	if p.is(token.COLON) {
		node.FlowType = p.parseFlowTypeAnnotation()
	}

	if p.is(token.ASSIGN) {
		p.next()
		node.Initializer = p.parseAssignmentExpression()
	}

	return node
}

func (p *Parser) parseVariableDeclarationList(loc *file.Loc, kind token.Token) []ast.Expression {

	var declarationList []*ast.VariableExpression // Avoid bad expressions
	var list []ast.Expression

	for {
		list = append(list, p.parseVariableDeclaration(
			&declarationList,
			kind,
		))

		if p.is(token.COMMA) {
			p.consumeExpected(token.COMMA)
		} else {
			break
		}
	}

	p.scope.declare(&ast.VariableDeclaration{
		Loc:  loc,
		Kind: kind,
		List: declarationList,
	})

	return list
}
