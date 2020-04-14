package ast

import "yawp/parser/file"

type (
	ForInStatement struct {
		Start  file.Loc
		Into   Expression
		Source Expression
		Body   Statement
	}

	ForOfStatement struct {
		Start    file.Loc
		Binder   Expression
		Iterator Expression
		Body     Statement
	}

	ForStatement struct {
		Start       file.Loc
		Initializer Expression
		Update      Expression
		Test        Expression
		Body        Statement
	}
)

func (*ForInStatement) _statementNode() {}
func (*ForOfStatement) _statementNode() {}
func (*ForStatement) _statementNode()   {}

func (f *ForInStatement) StartAt() file.Loc { return f.Start }
func (f *ForOfStatement) StartAt() file.Loc { return f.Start }
func (f *ForStatement) StartAt() file.Loc   { return f.Start }

func (f *ForInStatement) EndAt() file.Loc { return f.Body.EndAt() }
func (f *ForOfStatement) EndAt() file.Loc { return f.Body.EndAt() }
func (f *ForStatement) EndAt() file.Loc   { return f.Body.EndAt() }
