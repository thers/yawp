package ast

import "yawp/parser/file"

type (
	YieldExpression struct {
		Start      file.Idx
		Delegate   bool
		Expression Expression
	}

	YieldStatement struct {
		Expression *YieldExpression
	}
)

func (*YieldExpression) _expressionNode() {}

func (*YieldStatement) _statementNode() {}

func (y *YieldExpression) StartAt() file.Idx { return y.Start }
func (y *YieldStatement) StartAt() file.Idx  { return y.Expression.StartAt() }

func (y *YieldExpression) EndAt() file.Idx { return y.Expression.EndAt() }
func (y *YieldStatement) EndAt() file.Idx  { return y.Expression.EndAt() }
