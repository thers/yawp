package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseExportVarClause() *ast.ExportVarClause {
	var kind string

	switch p.token {
	case token.CONST:
		kind = "const"
	case token.LET:
		kind = "let"
	case token.VAR:
		kind = "var"
	}

	p.next()
	identifier := p.parseIdentifier()
	p.consumeExpected(token.ASSIGN)

	return &ast.ExportVarClause{
		Identifier:  identifier,
		Kind:        kind,
		Initializer: p.parseAssignmentExpression(),
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

func (p *Parser) parseExportDeclaration() *ast.ExportDeclaration {
	start := p.consumeExpected(token.EXPORT)
	declaration := &ast.ExportDeclaration{
		Start: start,
	}

	switch p.token {
	case token.CONST, token.LET, token.VAR:
		declaration.Clause = p.parseExportVarClause()
	case token.MULTIPLY:
		declaration.Clause = p.parseExportNamespaceFromClause()
	case token.LEFT_BRACE:
		declaration.Clause = p.parseExportNamedMaybeFromClause()
	}

	return declaration
}
