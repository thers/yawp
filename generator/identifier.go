package generator

import "yawp/parser/ast"

func (g *Generator) Identifier(id *ast.Identifier) *ast.Identifier {
	if id == nil {
		return id
	}

	if id.Ref != nil {
		if id.Ref.Kind == ast.RBuiltin && !id.Ref.Mangled {
			id.Ref.Mangled = true
			id.Ref.Name = g.ids.Next()
		}

		g.str(id.Ref.Name)
	} else {
		g.str(id.Name)
	}

	return id
}
