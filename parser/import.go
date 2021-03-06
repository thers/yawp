package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseImportDefaultClause(stmt *ast.ImportStatement) {
	identifier := p.symbol(p.parseIdentifier(), ast.SRead, ast.SRImport)
	moduleIdentifier := &ast.Identifier{
		ExprNode: p.exprNodeAt(identifier.Loc),
		Name:     "default",
	}

	exp := &ast.ImportClause{
		ExprNode:         p.exprNode(),
		Namespace:        false,
		ModuleIdentifier: moduleIdentifier,
		LocalIdentifier:  identifier,
	}

	stmt.Imports = append(stmt.Imports, exp)
	stmt.HasDefaultClause = true
}

func (p *Parser) parseImportNamedClause(stmt *ast.ImportStatement) {
	p.consumeExpected(token.LEFT_BRACE)

	for !p.is(token.RIGHT_BRACE) {
		if !p.is(token.IDENTIFIER) {
			p.unexpectedToken()
			return
		}

		localIdentifier := p.symbol(p.parseIdentifier(), ast.SDeclaration, ast.SRImport)
		moduleIdentifier := localIdentifier

		localIdentifier.Symbol.RefType = ast.SRImport

		p.allowToken(token.AS)
		if p.is(token.AS) {
			p.consumeExpected(token.AS)
			localIdentifier = p.symbol(p.parseIdentifier(), ast.SDeclaration, ast.SRImport)
		}

		stmt.Imports = append(stmt.Imports, &ast.ImportClause{
			ExprNode:         p.exprNodeAt(moduleIdentifier.Loc),
			Namespace:        false,
			ModuleIdentifier: moduleIdentifier,
			LocalIdentifier:  localIdentifier,
		})

		p.consumePossible(token.COMMA)
	}

	p.consumeExpected(token.RIGHT_BRACE)

	stmt.HasNamedClause = true
}

func (p *Parser) parseImportNamespaceClause(stmt *ast.ImportStatement) {
	exp := &ast.ImportClause{
		ExprNode:         p.exprNode(),
		Namespace:        true,
		ModuleIdentifier: nil,
		LocalIdentifier:  nil,
	}

	p.consumeExpected(token.MULTIPLY)
	p.allowToken(token.AS)
	p.consumeExpected(token.AS)

	if !p.is(token.IDENTIFIER) {
		p.unexpectedToken()
		return
	}

	exp.LocalIdentifier = p.symbol(p.parseIdentifier(), ast.SDeclaration, ast.SRImport)

	stmt.Imports = append(stmt.Imports, exp)
	stmt.HasNamespaceClause = true
}

func (p *Parser) parseImportFromClause(stmt *ast.ImportStatement) {
	if p.is(token.STRING) {
		stmt.From = p.literal
		p.next()
		stmt.Loc.End(p.tokenOffset)
	} else {
		p.unexpectedToken()
	}
}

func (p *Parser) parseImportDeclaration() *ast.ImportStatement {
	loc := p.loc()

	p.consumeExpected(token.IMPORT)
	stmt := &ast.ImportStatement{
		StmtNode: p.stmtNodeAt(loc),
		Imports:  make([]*ast.ImportClause, 0),
	}

	p.allowToken(token.TYPE_TYPE)
	if p.is(token.TYPE_TYPE) {
		stmt.Kind = ast.IKType
		p.next()
	} else if p.is(token.TYPEOF) {
		stmt.Kind = ast.IKTypeOf
		p.next()
	} else {
		stmt.Kind = ast.IKValue
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
			p.error(p.loc(), "Can not use multiple import clauses when first one isn't default")
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

	p.allowToken(token.FROM)
	p.consumeExpected(token.FROM)

	p.parseImportFromClause(stmt)

	return stmt
}

func (p *Parser) parseImportCall() ast.IExpr {
	loc := p.loc()

	p.consumeExpected(token.IMPORT)
	p.consumeExpected(token.LEFT_PARENTHESIS)

	loc = loc.End(p.consumeExpected(token.RIGHT_PARENTHESIS))

	call := &ast.ImportCall{
		ExprNode:   p.exprNodeAt(loc),
		Expression: p.parseAssignmentExpression(),
	}

	return call
}
