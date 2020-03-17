package ast

import "yawp/parser/file"

type (
	ArrayLiteral struct {
		Start file.Idx
		End   file.Idx
		Value []Expression
	}

	ArraySpread struct {
		Expression
	}
)

func (*ArrayLiteral) _expressionNode() {}
func (*ArraySpread) _expressionNode()  {}

func (a *ArrayLiteral) StartAt() file.Idx { return a.Start }
func (a *ArraySpread) StartAt() file.Idx { return a.Expression.StartAt() }

func (a *ArrayLiteral) EndAt() file.Idx { return a.End }
func (a *ArraySpread) EndAt() file.Idx { return a.Expression.EndAt() }
