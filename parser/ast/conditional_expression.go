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

func (c *ConditionalExpression) StartAt() file.Idx { return c.Test.StartAt() }

func (c *ConditionalExpression) EndAt() file.Idx { return c.Test.EndAt() }
