package generator

import "yawp/parser/ast"

func (g *Generator) FunctionLiteral(fl *ast.FunctionLiteral) *ast.FunctionLiteral {
	g.str("function")

	if fl.Id != nil {
		g.rune(' ')
		g.Identifier(fl.Id)
	}

	g.rune('(')

	g.FunctionParameters(fl.Parameters)

	g.str("){")

	g.Statement(fl.Body)

	g.rune('}')

	return fl
}
