package generator

import (
	"yawp/options"
	"yawp/parser/ast"
)

func (g *Generator) ClassStatement(c *ast.ClassStatement) ast.IStmt {
	g.ClassExpression(c.Expression)

	return c
}

func (g *Generator) ClassExpression(c *ast.ClassExpression) *ast.ClassExpression {
	if g.options.Target == options.ES5 {
		return g.classES5(c)
	}

	g.str("class ")
	g.str(c.Name.Name)

	if c.SuperClass != nil {
		g.str(" extends ")
		g.Expression(c.SuperClass)
	}

	g.rune('{')
	g.rune('}')

	return c
}

func (g *Generator) classES5(c *ast.ClassExpression) *ast.ClassExpression {
	g.str("'es5 class expression'")

	return c
}
