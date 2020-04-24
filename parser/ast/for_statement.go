package ast

import "yawp/parser/file"

type (
	ForInStatement struct {
		Loc    *file.Loc
		Into   Expression
		Source Expression
		Body   Statement
	}

	ForOfStatement struct {
		Loc      *file.Loc
		Binder   Expression
		Iterator Expression
		Body     Statement
	}

	ForStatement struct {
		Loc         *file.Loc
		Initializer *VariableStatement
		Update      Expression
		Test        Expression
		Body        Statement
	}
)

func (*ForInStatement) _statementNode() {}
func (*ForOfStatement) _statementNode() {}
func (*ForStatement) _statementNode()   {}

func (f *ForStatement) GetLoc() *file.Loc   { return f.Loc }
func (f *ForInStatement) GetLoc() *file.Loc { return f.Loc }
func (f *ForOfStatement) GetLoc() *file.Loc { return f.Loc }
