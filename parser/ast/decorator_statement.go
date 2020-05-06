package ast

import "yawp/parser/file"

type (
	LegacyDecoratorStatement struct {
		Decorators []Expression
		Subject    Statement
	}
)

func (l *LegacyDecoratorStatement) _statementNode() {}

func (l *LegacyDecoratorStatement) GetLoc() *file.Loc { return l.Decorators[0].GetLoc() }
