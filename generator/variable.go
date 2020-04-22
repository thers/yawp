package generator

import "yawp/parser/ast"

func (g *Generator) variableExpression(exp *ast.VariableExpression) {
	g.ref(g.refScope.SetRef(exp.Kind, exp.Name))

	if exp.Initializer != nil {
		g.rune('=')
		g.expression(exp.Initializer)
	}
}

func (g *Generator) variableBinding(bnd *ast.VariableBinding) {

}

func (g *Generator) variableStatement(stmt *ast.VariableStatement) {
	g.str(stmt.Kind.String())
	g.rune(' ')

	for index, decl := range stmt.List {
		if index > 0 {
			g.rune(',')
		}

		switch exp := decl.(type) {
		case *ast.VariableExpression:
			g.variableExpression(exp)
		case *ast.VariableBinding:
			g.variableBinding(exp)
		}
	}

	g.semicolon()
}
