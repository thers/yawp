package ast

import "yawp/parser/file"

type (
	NewTargetExpression struct {
		Start file.Idx
		End   file.Idx
	}
)

func (*NewTargetExpression) _expressionNode() {}

func (n *NewTargetExpression) StartAt() file.Idx { return n.Start }

func (n *NewTargetExpression) EndAt() file.Idx { return n.End }
