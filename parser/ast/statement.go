package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

// All statement nodes implement the IStmt interface.
type IStmt interface {
	INode
	_statementNode()
}
type StmtNode Node

func (*StmtNode) _statementNode() {}

func (e *StmtNode) GetNode() *Node    { return (*Node)(e) }
func (s *StmtNode) GetLoc() *file.Loc { return s.Loc }

type Statements []IStmt

func (Statements) _statementNode()   {}
func (Statements) GetLoc() *file.Loc { panic("Can not invoke GetLoc on Statements") }
func (Statements) GetNode() *Node    { panic("Can not invoke GetNode on Statements") }

type (
	BlockStatement struct {
		StmtNode
		List []IStmt
	}

	BranchStatement struct {
		StmtNode
		Token token.Token
		Label *Identifier
	}

	CaseStatement struct {
		StmtNode
		Test       IExpr
		Consequent []IStmt
	}

	DebuggerStatement struct {
		StmtNode
	}

	DoWhileStatement struct {
		StmtNode
		Test IExpr
		Body IStmt
	}

	EmptyStatement struct {
		StmtNode
	}

	ExpressionStatement struct {
		StmtNode
		Expression IExpr
	}

	IfStatement struct {
		StmtNode
		Test       IExpr
		Consequent IStmt
		Alternate  IStmt
	}

	LabelledStatement struct {
		StmtNode
		Label     *Identifier
		Colon     file.Idx
		Statement IStmt
	}

	ReturnStatement struct {
		StmtNode
		Argument IExpr
	}

	SwitchStatement struct {
		StmtNode
		Discriminant IExpr
		Default      int
		Body         []IStmt
	}

	ThrowStatement struct {
		StmtNode
		Argument IExpr
	}

	VariableStatement struct {
		StmtNode
		Kind token.Token
		List []*VariableBinding
	}

	WhileStatement struct {
		StmtNode
		Test IExpr
		Body IStmt
	}

	WithStatement struct {
		StmtNode
		Object IExpr
		Body   IStmt
	}

	ForInStatement struct {
		StmtNode
		Left  IStmt
		Right IExpr
		Body  IStmt
	}

	ForOfStatement struct {
		StmtNode
		Left  IStmt
		Right IExpr
		Body  IStmt
	}

	ForStatement struct {
		StmtNode
		Initializer IStmt
		Update      IExpr
		Test        IExpr
		Body        IStmt
	}

	LegacyDecoratorStatement struct {
		StmtNode
		Decorators []IExpr
		Subject    IStmt
	}

	TryStatement struct {
		StmtNode
		Body    IStmt
		Catch   IStmt
		Finally IStmt
	}

	CatchStatement struct {
		StmtNode
		Parameter PatternBinder
		Body      IStmt
	}

	ClassStatement struct {
		StmtNode
		Expression *ClassExpression
	}

	ClassFieldStatement struct {
		StmtNode
		Name        ObjectPropertyName
		Static      bool
		Private     bool
		Type        FlowType
		Initializer IExpr
	}

	ClassAccessorStatement struct {
		StmtNode
		Field ObjectPropertyName
		Kind  string
		Body  *FunctionLiteral
	}

	ClassMethodStatement struct {
		StmtNode
		Name           ObjectPropertyName
		Static         bool
		Private        bool
		Async          bool
		Generator      bool
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     *FunctionParameters
		Body           *FunctionBody
	}

	ExportStatement struct {
		StmtNode
		Clause ExportClause
	}

	ImportStatement struct {
		StmtNode
		Kind               ImportKind
		Imports            []*ImportClause
		From               string
		HasNamespaceClause bool
		HasDefaultClause   bool
		HasNamedClause     bool
	}

	FlowTypeStatement struct {
		StmtNode
		Name           *FlowIdentifier
		Opaque         bool
		Type           FlowType
		TypeParameters []*FlowTypeParameter
	}

	FlowInterfaceStatement struct {
		StmtNode
		Name           *FlowIdentifier
		TypeParameters []*FlowTypeParameter
		Body           []FlowInterfaceBodyStatement
	}
)
