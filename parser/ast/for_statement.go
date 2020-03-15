package ast

import "yawp/parser/file"

type (
	ForInStatement struct {
		Start  file.Idx
		Into   Expression
		Source Expression
		Body   Statement
	}

	ForOfStatement struct {
		Start    file.Idx
		Binder   Expression
		Iterator Expression
		Body     Statement
	}

	ForStatement struct {
		Start       file.Idx
		Initializer Expression
		Update      Expression
		Test        Expression
		Body        Statement
	}
)

func (*ForInStatement) _statementNode() {}
func (*ForOfStatement) _statementNode() {}
func (*ForStatement) _statementNode()   {}

func (f *ForInStatement) StartAt() file.Idx { return f.Start }
func (f *ForOfStatement) StartAt() file.Idx { return f.Start }
func (f *ForStatement) StartAt() file.Idx   { return f.Start }

func (f *ForInStatement) EndAt() file.Idx { return f.Body.EndAt() }
func (f *ForOfStatement) EndAt() file.Idx { return f.Body.EndAt() }
func (f *ForStatement) EndAt() file.Idx   { return f.Body.EndAt() }
