package ast

import "yawp/parser/file"

type (
	CoalesceExpression struct {
		Head       Expression
		Consequent Expression
	}
)

func (*CoalesceExpression) _expressionNode() {}

func (c *CoalesceExpression) GetLoc() *file.Loc {
	return c.Head.GetLoc().Add(c.Consequent.GetLoc())
}
