package generator

import (
	"yawp/options"
	"yawp/parser/ast"
)

func (g *Generator) classStatement(c *ast.ClassStatement) {
	g.class(c.Expression)
}

func (g *Generator) class(c *ast.ClassExpression) {
	if g.opt.Target == options.ES5 {
		g.classES5(c)
		return
	}

	g.str("class ")
	g.str(c.Name.Name)

	if c.SuperClass != nil {
		g.str(" extends ")
		g.expression(c.SuperClass)
	}

	g.rune('{')
	g.rune('}')
}

func (g *Generator) classES5(_ *ast.ClassExpression) {
	g.str("'es5 class expression'")
}
