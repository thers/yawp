package generator

import "yawp/parser/ast"

func (g *Generator) Identifier(id *ast.Identifier) *ast.Identifier {
	if id == nil {
		return id
	}

	if id.LegacyRef != nil {
		if id.LegacyRef.Type == ast.SRBuiltin && !id.LegacyRef.Mangled {
			id.LegacyRef.Mangled = true
			id.LegacyRef.Name = g.ids.Next()
		}

		g.str(id.LegacyRef.Name)
	} else {
		g.str(id.Name)
	}

	return id
}
