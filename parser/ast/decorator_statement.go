package ast

import "yawp/parser/file"

type (
	LegacyDecoratorSubject interface {
		EndAt() file.Loc
		_legacyDecoratorSubject()
	}

	LegacyDecoratorStatement struct {
		Decorators []Expression
		Subject    LegacyDecoratorSubject
	}
)

func (c *ClassStatement) _legacyDecoratorSubject()       {}
func (c *ClassFieldStatement) _legacyDecoratorSubject()  {}
func (c *ClassMethodStatement) _legacyDecoratorSubject() {}

func (l *LegacyDecoratorStatement) _statementNode() {}

func (l *LegacyDecoratorStatement) StartAt() file.Loc { return l.Decorators[0].StartAt() }

func (l *LegacyDecoratorStatement) EndAt() file.Loc { return l.Subject.EndAt() }
