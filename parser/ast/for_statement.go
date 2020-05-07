package ast

import "yawp/parser/file"

type (
	ForInStatement struct {
		Loc   *file.Loc
		Left  Statement
		Right Expression
		Body  Statement
	}

	ForOfStatement struct {
		Loc   *file.Loc
		Left  Statement
		Right Expression
		Body  Statement
	}

	ForStatement struct {
		Loc         *file.Loc
		Initializer Statement
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
