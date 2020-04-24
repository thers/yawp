package generator

import (
	"fmt"
	"yawp/options"
	"yawp/parser/ast"
)

func (g *Generator) ObjectBinding(ob *ast.ObjectBinding) *ast.ObjectBinding {
	if !g.options.Minify && g.options.Target == options.ES2020 {
		g.src(ob.GetLoc())

		return ob
	}

	g.str(fmt.Sprintf(
		"var %s=%s;",
		g.ids.Next(),
		g.ids.Next(),
	))

	g.rune('{')
	defer g.rune('}')



	return g.DefaultVisitor.ObjectBinding(ob)
}
