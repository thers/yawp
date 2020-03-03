package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseImportDefaultClause(stmt *ast.ImportDeclaration) {
	identifier := p.parseIdentifier()
	moduleIdentifier := &ast.Identifier{
		Start: identifier.Start,
		Name:  "default",
	}

	exp := &ast.ImportClause{
		Start:            p.idx,
		Namespace:        false,
		ModuleIdentifier: moduleIdentifier,
		LocalIdentifier:  identifier,
	}

	stmt.Imports = append(stmt.Imports, exp)
	stmt.HasDefaultClause = true
}

func (p *Parser) parseImportNamedClause(stmt *ast.ImportDeclaration) {
	p.consumeExpected(token.LEFT_BRACE)

	for !p.is(token.RIGHT_BRACE) {
		localIdentifier := p.parseIdentifier()
		moduleIdentifier := localIdentifier

		if p.is(token.AS) {
			p.consumeExpected(token.AS)
			localIdentifier = p.parseIdentifier()
		}

		stmt.Imports = append(stmt.Imports, &ast.ImportClause{
			Start:            moduleIdentifier.Start,
			Namespace:        false,
			ModuleIdentifier: moduleIdentifier,
			LocalIdentifier:  localIdentifier,
		})

		if p.is(token.COMMA) {
			p.next()
		}
	}

	p.consumeExpected(token.RIGHT_BRACE)

	stmt.HasNamedClause = true
}

func (p *Parser) parseImportNamespaceClause(stmt *ast.ImportDeclaration) {
	exp := &ast.ImportClause{
		Start:            p.idx,
		Namespace:        true,
		ModuleIdentifier: nil,
		LocalIdentifier:  nil,
	}

	p.consumeExpected(token.MULTIPLY)
	p.consumeExpected(token.AS)

	exp.LocalIdentifier = p.parseIdentifier()

	stmt.Imports = append(stmt.Imports, exp)
	stmt.HasNamespaceClause = true
}

func (p *Parser) parseImportFromClause(stmt *ast.ImportDeclaration) {
	if p.is(token.STRING) {
		stmt.From = p.literal
		p.next()
		stmt.End = p.idx
	} else {
		p.unexpectedToken()
	}
}

func (p *Parser) parseImportDeclaration() *ast.ImportDeclaration {
	start := p.idx

	p.consumeExpected(token.IMPORT)
	stmt := &ast.ImportDeclaration{
		Start:   start,
		Imports: make([]*ast.ImportClause, 0),
		End:     0,
	}

	switch p.token {
	case token.MULTIPLY:
		// import * as identifier
		p.parseImportNamespaceClause(stmt)
	case token.IDENTIFIER:
		// import defaultExport
		p.parseImportDefaultClause(stmt)
	case token.LEFT_BRACE:
		// import { namedExport, another as anotherExport }
		p.parseImportNamedClause(stmt)
	case token.STRING:
		// import "module"
		p.parseImportFromClause(stmt)
		return stmt
	}

	if p.is(token.COMMA) {
		if !stmt.HasDefaultClause {
			p.error(p.idx, "Can not use multiple import clauses when first one isn't default")
		} else {
			p.next()

			switch p.token {
			case token.MULTIPLY:
				// import default, * as identifier
				p.parseImportNamespaceClause(stmt)
			case token.LEFT_BRACE:
				// import default, { namedExport, another as anotherExport }
				p.parseImportNamedClause(stmt)
			default:
				p.unexpectedToken()
				p.next()
			}
		}
	}

	p.consumeExpected(token.FROM)

	p.parseImportFromClause(stmt)

	return stmt
}
