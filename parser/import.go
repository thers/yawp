package parser

import (
	"yawp/parser/ast"
	"yawp/parser/importKind"
	"yawp/parser/token"
)

func (p *Parser) parseImportDefaultClause(stmt *ast.ImportDeclaration) {
	identifier := p.parseIdentifier()
	moduleIdentifier := &ast.Identifier{
		Start: identifier.Start,
		Name:  "default",
	}

	exp := &ast.ImportClause{
		Start:            p.loc,
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

		p.allowNext(token.AS)
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

		p.consumePossible(token.COMMA)
	}

	p.consumeExpected(token.RIGHT_BRACE)

	stmt.HasNamedClause = true
}

func (p *Parser) parseImportNamespaceClause(stmt *ast.ImportDeclaration) {
	exp := &ast.ImportClause{
		Start:            p.loc,
		Namespace:        true,
		ModuleIdentifier: nil,
		LocalIdentifier:  nil,
	}

	p.consumeExpected(token.MULTIPLY)
	p.allowNext(token.AS)
	p.consumeExpected(token.AS)

	exp.LocalIdentifier = p.parseIdentifier()

	stmt.Imports = append(stmt.Imports, exp)
	stmt.HasNamespaceClause = true
}

func (p *Parser) parseImportFromClause(stmt *ast.ImportDeclaration) {
	if p.is(token.STRING) {
		stmt.From = p.literal
		p.next()
		stmt.End = p.loc
	} else {
		p.unexpectedToken()
	}
}

func (p *Parser) parseImportDeclaration() *ast.ImportDeclaration {
	start := p.loc

	p.consumeExpected(token.IMPORT)
	stmt := &ast.ImportDeclaration{
		Start:   start,
		Imports: make([]*ast.ImportClause, 0),
		End:     0,
	}

	p.allowNext(token.TYPE_TYPE)
	if p.is(token.TYPE_TYPE) {
		stmt.Kind = importKind.TYPE
		p.next()
	} else if p.is(token.TYPEOF) {
		stmt.Kind = importKind.TYPEOF
		p.next()
	} else {
		stmt.Kind = importKind.VALUE
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
			p.error(p.loc, "Can not use multiple import clauses when first one isn't default")
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

	p.allowNext(token.FROM)
	p.consumeExpected(token.FROM)

	p.parseImportFromClause(stmt)

	return stmt
}

func (p *Parser) parseImportCall() ast.Expression {
	start := p.consumeExpected(token.IMPORT)
	p.consumeExpected(token.LEFT_PARENTHESIS)

	call := &ast.ImportCall{
		Start:      start,
		Expression: p.parseAssignmentExpression(),
		End:        p.consumeExpected(token.RIGHT_PARENTHESIS),
	}

	return call
}
