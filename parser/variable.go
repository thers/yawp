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
	idx := p.idx
	kind := p.token

	if !p.isVariableStatementStart() {
		p.unexpectedToken()
	}

	p.next()
	list := p.parseVariableDeclarationList(idx, kind)
	p.semicolon()

	return &ast.VariableStatement{
		Kind: kind,
		Var:  idx,
		List: list,
	}
}

func (p *Parser) parseVariableDeclaration(declarationList *[]*ast.VariableExpression, kind token.Token) ast.Expression {

	if !p.is(token.IDENTIFIER) {
		idx := p.consumeExpected(token.IDENTIFIER)
		p.nextStatement()
		return &ast.BadExpression{From: idx, To: p.idx}
	}

	literal := p.literal
	idx := p.idx
	p.next()
	node := &ast.VariableExpression{
		Kind:  kind,
		Name:  literal,
		Start: idx,
	}

	if declarationList != nil {
		*declarationList = append(*declarationList, node)
	}

	if p.is(token.ASSIGN) {
		p.next()
		node.Initializer = p.parseAssignmentExpression()
	}

	return node
}

func (p *Parser) parseVariableDeclarationList(var_ file.Idx, kind token.Token) []ast.Expression {

	var declarationList []*ast.VariableExpression // Avoid bad expressions
	var list []ast.Expression

	for {
		list = append(list, p.parseVariableDeclaration(
			&declarationList,
			kind,
		))
		if !p.is(token.COMMA) {
			break
		}
		p.next()
	}

	p.scope.declare(&ast.VariableDeclaration{
		Kind: kind,
		Var:  var_,
		List: declarationList,
	})

	return list
}

//func (p *Parser) parseArrayAssignmentPatternOrArrayLiteral() ast.Expression {
//
//}
