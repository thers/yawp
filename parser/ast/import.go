package ast

import "yawp/parser/file"

type (
	ImportClause struct {
		Start            file.Idx
		Namespace        bool
		ModuleIdentifier *Identifier
		LocalIdentifier  *Identifier
	}

	ImportDeclaration struct {
		Start              file.Idx
		Imports            []*ImportClause
		From               string
		End                file.Idx
		HasNamespaceClause bool
		HasDefaultClause   bool
		HasNamedClause     bool
	}

	ImportCall struct {
		Start file.Idx
		End   file.Idx
		Expression
	}
)

func (*ImportClause) _expressionNode() {}
func (*ImportCall) _expressionNode()   {}

func (*ImportDeclaration) _statementNode() {}

func (i *ImportClause) StartAt() file.Idx      { return i.Start }
func (i *ImportCall) StartAt() file.Idx        { return i.Start }
func (i *ImportDeclaration) StartAt() file.Idx { return i.Start }

func (i *ImportClause) EndAt() file.Idx      { return i.LocalIdentifier.EndAt() }
func (i *ImportCall) EndAt() file.Idx        { return i.End }
func (i *ImportDeclaration) EndAt() file.Idx { return i.End }
