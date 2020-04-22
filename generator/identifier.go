package generator

import "yawp/parser/ast"

func (g *Generator) identifier(id *ast.Identifier) {
	g.str(id.Name)
}

func (g *Generator) ref(ref *ast.Ref) {
	if !ref.IsMangled && g.opt.Minify {
		ref.IsMangled = true
		ref.Name = g.refScope.NextMangledId()
	}

	g.str(ref.Name)
}

func (g *Generator) refOrExpression(node ast.Expression) {
	if refId, ok := node.(*ast.Identifier); ok {
		g.ref(g.refScope.UseRef(refId.Name))
	} else {
		g.expression(node)
	}
}
