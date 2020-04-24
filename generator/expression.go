package generator

import "yawp/parser/ast"

func (g *Generator) BooleanLiteral(b *ast.BooleanLiteral) *ast.BooleanLiteral {
	if !g.options.Minify {
		g.str(b.Literal)
		return b
	}

	if b.Literal == ast.LBooleanFalse {
		g.str("!1")
	} else {
		g.str("!0")
	}

	return b
}

func (g *Generator) StringLiteral(s *ast.StringLiteral) *ast.StringLiteral {
	g.str(s.Literal)

	return s
}

func (g *Generator) NumberLiteral(s *ast.NumberLiteral) *ast.NumberLiteral {
	g.str(s.Literal)

	return s
}

func (g *Generator) BinaryExpression(b *ast.BinaryExpression) *ast.BinaryExpression {
	g.Expression(b.Left)
	g.str(b.Operator.String())
	g.Expression(b.Right)

	return b
}

func (g *Generator) AssignExpression(a *ast.AssignExpression) *ast.AssignExpression {
	g.Expression(a.Left)
	g.str(a.Operator.String())
	g.Expression(a.Right)

	return a
}

func (g *Generator) CallExpression(c *ast.CallExpression) *ast.CallExpression {
	g.Expression(c.Callee)
	g.rune('(')

	for index, arg := range c.ArgumentList {
		if index > 0 {
			g.rune(',')
		}

		g.Expression(arg)
	}

	g.rune(')')

	return c
}

func (g *Generator) DotExpression(d *ast.DotExpression) *ast.DotExpression {
	g.Expression(d.Left)

	g.rune('.')
	g.Identifier(d.Identifier)

	return d
}
