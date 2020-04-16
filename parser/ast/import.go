package ast

import (
	"yawp/parser/file"
	"yawp/parser/importKind"
)

type (
	ImportClause struct {
		Loc              *file.Loc
		Namespace        bool
		ModuleIdentifier *Identifier
		LocalIdentifier  *Identifier
	}

	ImportDeclaration struct {
		Loc                *file.Loc
		Kind               importKind.ImportKind
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

func (*ImportDeclaration) _statementNode() {}

func (i *ImportClause) GetLoc() *file.Loc      { return i.Loc }
func (i *ImportCall) GetLoc() *file.Loc        { return i.Loc }
func (i *ImportDeclaration) GetLoc() *file.Loc { return i.Loc }
