package generator

import (
	"yawp/parser/ast"
)

func (g *Generator) ObjectBinding(ob *ast.ObjectBinding) *ast.ObjectBinding {
	if !g.options.Minify {
		g.src(ob.GetLoc())

		return ob
	}

	g.rune('{')
	defer g.rune('}')

	return g.DefaultVisitor.ObjectBinding(ob)
}
