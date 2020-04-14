package ast

import "yawp/parser/file"

type (
	ArrayLiteral struct {
		Start file.Loc
		End   file.Loc
		Value []Expression
	}

	ArraySpread struct {
		Expression
	}
)

func (*ArrayLiteral) _expressionNode() {}
func (*ArraySpread) _expressionNode()  {}

func (a *ArrayLiteral) StartAt() file.Loc { return a.Start }
func (a *ArraySpread) StartAt() file.Loc  { return a.Expression.StartAt() }

func (a *ArrayLiteral) EndAt() file.Loc { return a.End }
func (a *ArraySpread) EndAt() file.Loc  { return a.Expression.EndAt() }
