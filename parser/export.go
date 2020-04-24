package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseExportVarClause() *ast.ExportVarClause {
	return &ast.ExportVarClause{
		Declaration: p.parseVariableStatement(),
	}
}

func (p *Parser) parseExportNamespaceFromClause() *ast.ExportNamespaceFromClause {
	p.consumeExpected(token.MULTIPLY)

	clause := &ast.ExportNamespaceFromClause{
		ModuleIdentifier: nil,
		From:             "",
	}

	if p.is(token.AS) {
		p.consumeExpected(token.AS)
		clause.ModuleIdentifier = p.parseIdentifier()
	}

	p.allowToken(token.FROM)
	p.consumeExpected(token.FROM)

	if p.is(token.STRING) {
		clause.From = p.literal
	} else {
		p.unexpectedToken()
		p.next()
	}

	return clause
}

func (p *Parser) parseExportNamedClause() *ast.ExportNamedClause {
	clause := &ast.ExportNamedClause{
		Exports: make([]*ast.NamedExportClause, 0),
	}

	p.consumeExpected(token.LEFT_BRACE)
	p.allowToken(token.AS)

	for !p.is(token.RIGHT_BRACE) {
		localIdentifier := p.parseIdentifier()
		moduleIdentifier := localIdentifier

		if p.is(token.AS) {
			p.consumeExpected(token.AS)
			moduleIdentifier = p.parseIdentifier()
		}

		clause.Exports = append(clause.Exports, &ast.NamedExportClause{
			ModuleIdentifier: moduleIdentifier,
			LocalIdentifier:  localIdentifier,
		})

		p.consumePossible(token.COMMA)
	}

	p.consumeExpected(token.RIGHT_BRACE)

	return clause
}

func (p *Parser) parseExportNamedMaybeFromClause() ast.ExportClause {
	exportNamedClause := p.parseExportNamedClause()

	p.allowToken(token.FROM)
	if p.is(token.FROM) {
		p.next()

		var from string

		if p.is(token.STRING) {
			from = p.literal
		}

		return &ast.ExportNamedFromClause{
			Exports: exportNamedClause.Exports,
			From:    from,
		}
	}

	return exportNamedClause
}

func (p *Parser) parseExportDefaultClause() *ast.ExportDefaultClause {
	p.consumeExpected(token.DEFAULT)

	clause := &ast.ExportDefaultClause{
		Declaration: p.parseAssignmentExpression(),
	}

	return clause
}

func (p *Parser) parseExportDeclaration() *ast.ExportDeclaration {
	loc := p.loc()
	p.consumeExpected(token.EXPORT)
	p.allowToken(token.TYPE_TYPE)

	declaration := &ast.ExportDeclaration{
		Loc: loc,
	}

	switch p.token {
	case token.CONST, token.LET, token.VAR:
		declaration.Clause = p.parseExportVarClause()
	case token.MULTIPLY:
		declaration.Clause = p.parseExportNamespaceFromClause()
	case token.LEFT_BRACE:
		declaration.Clause = p.parseExportNamedMaybeFromClause()
	case token.DEFAULT:
		declaration.Clause = p.parseExportDefaultClause()
	case token.TYPE_TYPE:
		declaration.Clause = p.parseFlowTypeStatement()
	case token.INTERFACE:
		declaration.Clause = p.parseFlowInterfaceStatement()
	}

	return declaration
}
