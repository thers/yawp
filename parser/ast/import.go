package ast

import (
	"yawp/parser/file"
	"yawp/parser/importKind"
)

type (
	ImportClause struct {
		Start            file.Loc
		Namespace        bool
		ModuleIdentifier *Identifier
		LocalIdentifier  *Identifier
	}

	ImportDeclaration struct {
		Start              file.Loc
		Kind               importKind.ImportKind
		Imports            []*ImportClause
		From               string
		End                file.Loc
		HasNamespaceClause bool
		HasDefaultClause   bool
		HasNamedClause     bool
	}

	ImportCall struct {
		Start file.Loc
		End   file.Loc
		Expression
	}
)

func (*ImportClause) _expressionNode() {}
func (*ImportCall) _expressionNode()   {}

func (*ImportDeclaration) _statementNode() {}

func (i *ImportClause) StartAt() file.Loc      { return i.Start }
func (i *ImportCall) StartAt() file.Loc        { return i.Start }
func (i *ImportDeclaration) StartAt() file.Loc { return i.Start }

func (i *ImportClause) EndAt() file.Loc      { return i.LocalIdentifier.EndAt() }
func (i *ImportCall) EndAt() file.Loc        { return i.End }
func (i *ImportDeclaration) EndAt() file.Loc { return i.End }
