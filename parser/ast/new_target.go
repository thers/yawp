package ast

import "yawp/parser/file"

type (
	NewTargetExpression struct {
		Loc *file.Loc
	}
)

func (*NewTargetExpression) _expressionNode() {}

func (n *NewTargetExpression) GetLoc() *file.Loc { return n.Loc }
