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
)

func (*ImportClause) _expressionNode()     {}
func (*ImportDeclaration) _statementNode() {}

func (self *ImportClause) StartAt() file.Idx      { return self.Start }
func (self *ImportDeclaration) StartAt() file.Idx { return self.Start }

func (self *ImportClause) EndAt() file.Idx      { return self.LocalIdentifier.EndAt() }
func (self *ImportDeclaration) EndAt() file.Idx { return self.End }
