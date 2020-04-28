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

func (g *Generator) ArrayLiteral(al *ast.ArrayLiteral) *ast.ArrayLiteral {
	g.rune('[')
	defer g.rune(']')

	for index, item := range al.List {
		al.List[index] = g.Expression(item)
	}

	return al
}

func (g *Generator) ObjectLiteral(o *ast.ObjectLiteral) *ast.ObjectLiteral {
	g.rune('{')
	defer g.rune('}')

	for _, prop := range o.Properties {
		switch p := prop.(type) {
		case *ast.ObjectPropertyValue:
			g.ObjectPropertyName(p.PropertyName)
			g.rune(':')
			g.Expression(p.Value)
		}
	}

	return o
}

func (g *Generator) ObjectPropertyName(opn ast.ObjectPropertyName) ast.ObjectPropertyName {
	switch o := opn.(type) {
	case *ast.Identifier:
		return g.Identifier(o)
	case *ast.ComputedName:
		return g.ComputedName(o)

	default:
		panic("Unknown object property name type")
	}
}

func (g *Generator) ComputedName(cn *ast.ComputedName) *ast.ComputedName {
	g.rune('[')
	g.Expression(cn.Expression)
	g.rune(']')

	return cn
}

func (g *Generator) BinaryExpression(b *ast.BinaryExpression) *ast.BinaryExpression {
	if g.wrapExpression {
		g.rune('(')
		defer g.rune(')')
	}

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

func (g *Generator) BracketExpression(be *ast.BracketExpression) *ast.BracketExpression {
	wrapExpression := g.wrapExpression
	g.wrapExpression = true
	g.Expression(be.Left)
	g.wrapExpression = wrapExpression

	g.rune('[')
	g.Expression(be.Member)
	g.rune(']')

	return be
}
