package generator

import "yawp/parser/ast"

func (g *Generator) Identifier(id *ast.Identifier) *ast.Identifier {
	if id == nil {
		return id
	}

	if id.Ref != nil {
		g.str(id.Ref.Name)
	} else {
		g.str(id.Name)
	}

	return id
}
