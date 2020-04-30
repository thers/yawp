package generator

import "yawp/parser/ast"

func (g *Generator) Js(js *ast.Js) *ast.Js {
	g.str(js.Code)

	return js
}
