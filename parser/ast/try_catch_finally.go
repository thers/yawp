package ast

import "yawp/parser/file"

type (
	TryStatement struct {
		Loc     *file.Loc
		Body    Statement
		Catch   Statement
		Finally Statement
	}

	CatchStatement struct {
		Loc       *file.Loc
		Parameter PatternBinder
		Body      Statement
	}
)

func (*TryStatement) _statementNode()   {}
func (*CatchStatement) _statementNode() {}
