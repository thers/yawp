package generator

import (
	"yawp/parser/ast"
)

func (g *Generator) VariableBinding(b *ast.VariableBinding) *ast.VariableBinding {
	g.PatternBinder(b.Binder)

	if b.Initializer != nil {
		g.rune('=')
		g.Expression(b.Initializer)
	}

	return b
}

func (g *Generator) VariableStatement(stmt *ast.VariableStatement) ast.IStmt {
	g.str(stmt.Kind.String())
	g.rune(' ')

	for index, binding := range stmt.List {
		if index > 0 {
			g.rune(',')
		}

		g.VariableBinding(binding)
	}

	g.semicolon()

	return stmt
}
