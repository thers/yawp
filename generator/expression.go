package generator

import "yawp/parser/ast"

func (g *Generator) expression(aexp ast.Expression) {
	if aexp == nil {
		return
	}

	switch exp := aexp.(type) {
	case *ast.Identifier:
		g.identifier(exp)
	case *ast.BooleanLiteral:
		g.boolean(exp)
	case *ast.StringLiteral:
		g.string(exp)
	case *ast.NumberLiteral:
		g.number(exp)
	case *ast.BinaryExpression:
		g.binary(exp)
	case *ast.AssignExpression:
		g.assign(exp)
	case *ast.CallExpression:
		g.call(exp)
	case *ast.DotExpression:
		g.dot(exp)
	default:
		g.str("'unknown expression';")
	}
}

func (g *Generator) boolean(b *ast.BooleanLiteral) {
	if !g.opt.Minify {
		g.str(b.Literal)
		return
	}

	if b.Literal == ast.LBooleanFalse {
		g.str("!1")
	} else {
		g.str("!0")
	}
}

func (g *Generator) string(s *ast.StringLiteral) {
	g.str(s.Literal)
}

func (g *Generator) number(s *ast.NumberLiteral) {
	g.str(s.Literal)
}

func (g *Generator) binary(b *ast.BinaryExpression) {
	g.expression(b.Left)
	g.str(b.Operator.String())
	g.expression(b.Right)
}

func (g *Generator) assign(a *ast.AssignExpression) {
	g.expression(a.Left)
	g.str(a.Operator.String())
	g.expression(a.Right)
}

func (g *Generator) call(c *ast.CallExpression) {
	g.expression(c.Callee)
	g.rune('(')

	for index, arg := range c.ArgumentList {
		if index > 0 {
			g.rune(',')
		}

		g.expression(arg)
	}

	g.rune(')')
}

func (g *Generator) dot(d *ast.DotExpression) {
	g.expression(d.Left)
	g.rune('.')
	g.identifier(d.Identifier)
}
