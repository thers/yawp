package ast

import "yawp/parser/file"

type (
	LegacyDecoratorSubject interface {
		EndAt() file.Idx
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

func (l *LegacyDecoratorStatement) StartAt() file.Idx { return l.Decorators[0].StartAt() }

func (l *LegacyDecoratorStatement) EndAt() file.Idx { return l.Subject.EndAt() }
