package ast

import "yawp/parser/file"

type (
	ExportClause interface {
		_exportClauseNode()
	}

	ExportDeclaration struct {
		Start  file.Idx
		Clause ExportClause
	}

	// export * from 'module'
	ExportNamespaceFromClause struct {
		ModuleIdentifier *Identifier
		From             string
	}

	// a as b
	NamedExportClause struct {
		ModuleIdentifier *Identifier
		LocalIdentifier  *Identifier
	}

	// export { a as b } from 'module'
	ExportNamedFromClause struct {
		Exports []*NamedExportClause
		From    string
	}

	// export { a as b }
	ExportNamedClause struct {
		Exports []*NamedExportClause
	}

	// export var/const/let aa = ''
	ExportVarClause struct {
		Declaration *VariableStatement
	}

	// export function fn() {}
	ExportFunctionClause struct {
		FunctionLiteral *FunctionLiteral
	}

	// export class Test {}
	ExportClassClause struct {
		ClassExpression *ClassExpression
	}

	// export default
	ExportDefaultClause struct {
		Declaration Expression
	}
)

func (*ExportDeclaration) _statementNode()              {}
func (s *ExportNamespaceFromClause) _exportClauseNode() {}
func (s *ExportNamedFromClause) _exportClauseNode()     {}
func (s *ExportNamedClause) _exportClauseNode()         {}
func (s *ExportVarClause) _exportClauseNode()           {}
func (s *ExportFunctionClause) _exportClauseNode()      {}
func (s *ExportClassClause) _exportClauseNode()         {}
func (s *ExportDefaultClause) _exportClauseNode()       {}

func (s *ExportDeclaration) StartAt() file.Idx { return s.Start }

func (s *ExportDeclaration) EndAt() file.Idx { return -1 }
