package ast

import (
	"yawp/parser/file"
)

type (
	ImportClause struct {
		Loc              *file.Loc
		Namespace        bool
		ModuleIdentifier *Identifier
		LocalIdentifier  *Identifier
	}

	ImportStatement struct {
		Loc                *file.Loc
		Kind               ImportKind
		Imports            []*ImportClause
		From               string
		HasNamespaceClause bool
		HasDefaultClause   bool
		HasNamedClause     bool
	}

	ImportCall struct {
		Loc *file.Loc
		Expression
	}
)

func (*ImportClause) _expressionNode() {}
func (*ImportCall) _expressionNode()   {}

func (*ImportStatement) _statementNode() {}

func (i *ImportClause) GetLoc() *file.Loc      { return i.Loc }
func (i *ImportCall) GetLoc() *file.Loc      { return i.Loc }
func (i *ImportStatement) GetLoc() *file.Loc { return i.Loc }
