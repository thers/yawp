package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

// All nodes implement the Node interface.
type Node interface {
	GetLoc() *file.Loc
}

// ========== //
// Expression //
// ========== //

type (
	// All expression nodes implement the Expression interface.
	Expression interface {
		Node
		_expressionNode()
	}

	AssignExpression struct {
		Operator token.Token
		Left     Expression
		Right    Expression
	}

	BinaryExpression struct {
		Operator   token.Token
		Left       Expression
		Right      Expression
		Comparison bool
	}

	BooleanLiteral struct {
		Loc     *file.Loc
		Literal string
	}

	CallExpression struct {
		Callee        Expression
		TypeArguments []FlowType
		ArgumentList  []Expression
	}

	Identifier struct {
		Loc  *file.Loc
		Ref  *Ref
		Name string
	}

	ObjectDestructuring struct {
		Start       file.Idx
		Identifiers []*Identifier
	}

	RestExpression struct {
		Start file.Idx
		Name  string
	}

	SpreadExpression struct {
		Loc   *file.Loc
		Value Expression
	}

	NewExpression struct {
		Loc           *file.Loc
		Callee        Expression
		TypeArguments []FlowType
		ArgumentList  []Expression
	}

	NullLiteral struct {
		Loc     *file.Loc
		Literal string
	}

	NumberLiteral struct {
		Loc     *file.Loc
		Literal string
	}

	FunctionParameters struct {
		Loc  *file.Loc
		List []FunctionParameter
	}

	RegExpLiteral struct {
		Loc     *file.Loc
		Literal string
		Pattern string
		Flags   string
	}

	SequenceExpression struct {
		Loc      *file.Loc
		Sequence []Expression
	}

	StringLiteral struct {
		Loc     *file.Loc
		Raw     bool
		Literal string
	}

	ThisExpression struct {
		Loc     *file.Loc
		Private bool
	}

	UnaryExpression struct {
		Loc      *file.Loc
		Operator token.Token
		Operand  Expression
		Postfix  bool
	}

	ArrowFunctionExpression struct {
		Loc            *file.Loc
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     []FunctionParameter
		Body           Statement
		Async          bool
	}

	AwaitExpression struct {
		Loc        *file.Loc
		Expression Expression
	}
)

// _expressionNode

func (*AssignExpression) _expressionNode()        {}
func (*BinaryExpression) _expressionNode()        {}
func (*BooleanLiteral) _expressionNode()          {}
func (*CallExpression) _expressionNode()          {}
func (*Identifier) _expressionNode()              {}
func (*NewExpression) _expressionNode()           {}
func (*NullLiteral) _expressionNode()             {}
func (*NumberLiteral) _expressionNode()           {}
func (*RegExpLiteral) _expressionNode()           {}
func (*SequenceExpression) _expressionNode()      {}
func (*StringLiteral) _expressionNode()           {}
func (*ThisExpression) _expressionNode()          {}
func (*UnaryExpression) _expressionNode()         {}
func (*ArrowFunctionExpression) _expressionNode() {}
func (*AwaitExpression) _expressionNode()         {}
func (*SpreadExpression) _expressionNode()        {}

// ========= //
// Statement //
// ========= //

type (
	// All statement nodes implement the Statement interface.
	Statement interface {
		Node
		_statementNode()
	}

	BlockStatement struct {
		Loc  *file.Loc
		List []Statement
	}

	BranchStatement struct {
		Loc   *file.Loc
		Token token.Token
		Label *Identifier
	}

	CaseStatement struct {
		Loc        *file.Loc
		Test       Expression
		Consequent []Statement
	}

	CatchStatement struct {
		Loc       *file.Loc
		Parameter *Identifier
		Body      Statement
	}

	DebuggerStatement struct {
		Loc *file.Loc
	}

	DoWhileStatement struct {
		Loc  *file.Loc
		Test Expression
		Body Statement
	}

	EmptyStatement struct {
		Loc *file.Loc
	}

	ExpressionStatement struct {
		Expression Expression
	}

	IfStatement struct {
		Loc        *file.Loc
		Test       Expression
		Consequent Statement
		Alternate  Statement
	}

	LabelledStatement struct {
		Label     *Identifier
		Colon     file.Idx
		Statement Statement
	}

	ReturnStatement struct {
		Loc      *file.Loc
		Argument Expression
	}

	SwitchStatement struct {
		Loc          *file.Loc
		Discriminant Expression
		Default      int
		Body         []Statement
	}

	ThrowStatement struct {
		Loc      *file.Loc
		Argument Expression
	}

	TryStatement struct {
		Loc     *file.Loc
		Body    Statement
		Catch   Statement
		Finally Statement
	}

	VariableStatement struct {
		Loc  *file.Loc
		Kind token.Token
		List []*VariableBinding
	}

	WhileStatement struct {
		Loc  *file.Loc
		Test Expression
		Body Statement
	}

	WithStatement struct {
		Loc    *file.Loc
		Object Expression
		Body   Statement
	}
)

// _statementNode

func (*BlockStatement) _statementNode()      {}
func (*BranchStatement) _statementNode()     {}
func (*CaseStatement) _statementNode()       {}
func (*CatchStatement) _statementNode()      {}
func (*DebuggerStatement) _statementNode()   {}
func (*DoWhileStatement) _statementNode()    {}
func (*EmptyStatement) _statementNode()      {}
func (*ExpressionStatement) _statementNode() {}
func (*IfStatement) _statementNode()         {}
func (*LabelledStatement) _statementNode()   {}
func (*ReturnStatement) _statementNode()     {}
func (*SwitchStatement) _statementNode()     {}
func (*ThrowStatement) _statementNode()      {}
func (*TryStatement) _statementNode()        {}
func (*VariableStatement) _statementNode()   {}
func (*WhileStatement) _statementNode()      {}
func (*WithStatement) _statementNode()       {}

// =========== //
// Declaration //
// =========== //

type (
	// All declaration nodes implement the Declaration interface.
	Declaration interface {
		_declarationNode()
	}

	FunctionDeclaration struct {
		Function *FunctionLiteral
	}

	VariableDeclaration struct {
		Loc  *file.Loc
		Kind token.Token
		List []*VariableBinding
	}

	ClassDeclaration struct {
		Name    *Identifier
		Extends Expression
	}
)

// _declarationNode

func (*FunctionDeclaration) _declarationNode() {}
func (*VariableDeclaration) _declarationNode() {}
func (*ClassDeclaration) _declarationNode()    {}
