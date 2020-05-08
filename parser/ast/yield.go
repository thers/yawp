package ast

import "yawp/parser/file"

type (
	YieldExpression struct {
		Loc      *file.Loc
		Delegate bool
		Argument Expression
	}

	YieldStatement struct {
		Expression *YieldExpression
	}
)

func (*YieldExpression) _expressionNode() {}

func (*YieldStatement) _statementNode() {}

func (y *YieldExpression) GetLoc() *file.Loc { return y.Loc }
func (y *YieldStatement) GetLoc() *file.Loc  { return y.Expression.Loc }
