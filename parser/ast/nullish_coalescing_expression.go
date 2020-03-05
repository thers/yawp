package ast

import "yawp/parser/file"

type (
	NullishCoalescingExpression struct {
		Test       Expression
		Consequent Expression
	}
)

func (*NullishCoalescingExpression) _expressionNode() {}

func (c *NullishCoalescingExpression) StartAt() file.Idx { return c.Test.StartAt() }

func (c *NullishCoalescingExpression) EndAt() file.Idx { return c.Test.EndAt() }
