package generator

import "yawp/parser/ast"

func (g *Generator) identifier(id *ast.Identifier) {
	g.str(id.Name)
}
