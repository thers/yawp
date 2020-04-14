package ast

import "yawp/parser/file"

type (
	YieldExpression struct {
		Start      file.Loc
		Delegate   bool
		Expression Expression
	}

	YieldStatement struct {
		Expression *YieldExpression
	}
)

func (*YieldExpression) _expressionNode() {}

func (*YieldStatement) _statementNode() {}

func (y *YieldExpression) StartAt() file.Loc { return y.Start }
func (y *YieldStatement) StartAt() file.Loc  { return y.Expression.StartAt() }

func (y *YieldExpression) EndAt() file.Loc { return y.Expression.EndAt() }
func (y *YieldStatement) EndAt() file.Loc  { return y.Expression.EndAt() }
