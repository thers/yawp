package ast

import "yawp/parser/file"

type (
	CoalesceExpression struct {
		Head       Expression
		Consequent Expression
	}
)

func (*CoalesceExpression) _expressionNode() {}

func (c *CoalesceExpression) StartAt() file.Idx { return c.Head.StartAt() }

func (c *CoalesceExpression) EndAt() file.Idx { return c.Head.EndAt() }