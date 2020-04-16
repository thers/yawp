package ast

import "yawp/parser/file"

type (
	ConditionalExpression struct {
		Test       Expression
		Consequent Expression
		Alternate  Expression
	}
)

func (*ConditionalExpression) _expressionNode() {}

func (c *ConditionalExpression) GetLoc() *file.Loc {
	loc := c.Test.GetLoc().Add(c.Consequent.GetLoc())

	if c.Alternate != nil {
		loc = loc.Add(c.Alternate.GetLoc())
	}

	return loc
}
