package ast

import "yawp/parser/file"

type (
	NewTargetExpression struct {
		Start file.Loc
		End   file.Loc
	}
)

func (*NewTargetExpression) _expressionNode() {}

func (n *NewTargetExpression) StartAt() file.Loc { return n.Start }

func (n *NewTargetExpression) EndAt() file.Loc { return n.End }
