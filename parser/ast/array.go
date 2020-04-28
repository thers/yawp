package ast

import "yawp/parser/file"

type (
	ArrayLiteral struct {
		Loc  *file.Loc
		List []Expression
	}

	ArraySpread struct {
		Expression
	}
)

func (*ArrayLiteral) _expressionNode() {}
func (*ArraySpread) _expressionNode()  {}

func (a *ArrayLiteral) GetLoc() *file.Loc { return a.Loc }
