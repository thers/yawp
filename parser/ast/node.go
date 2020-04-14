/*
Package ast declares types representing a JavaScript AST.

Warning

The parser and AST interfaces are still works-in-progress (particularly where
node types are concerned) and may change in the future.

*/
package ast

import (
	"github.com/go-sourcemap/sourcemap"
	"yawp/parser/file"
	"yawp/parser/token"
)

// All nodes implement the Node interface.
type Node interface {
	StartAt() file.Loc // The index of the first character belonging to the node
	EndAt() file.Loc   // The index of the first character immediately after the node
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
		Start   file.Loc
		Literal string
		Value   bool
	}

	BracketExpression struct {
		Left         Expression
		Member       Expression
		LeftBracket  file.Loc
		RightBracket file.Loc
	}

	CallExpression struct {
		Callee           Expression
		TypeArguments    []FlowType
		LeftParenthesis  file.Loc
		ArgumentList     []Expression
		RightParenthesis file.Loc
	}

	DotExpression struct {
		Left       Expression
		Identifier *Identifier
	}

	Identifier struct {
		Start file.Loc
		Name  string
	}

	ObjectDestructuring struct {
		Start       file.Loc
		Identifiers []*Identifier
	}

	RestExpression struct {
		Start file.Loc
		Name  string
	}

	SpreadExpression struct {
		Start file.Loc
		Value Expression
	}

	NewExpression struct {
		Start            file.Loc
		Callee           Expression
		TypeArguments    []FlowType
		LeftParenthesis  file.Loc
		ArgumentList     []Expression
		RightParenthesis file.Loc
	}

	NullLiteral struct {
		Start   file.Loc
		Literal string
	}

	NumberLiteral struct {
		Start   file.Loc
		Literal string
		Value   interface{}
	}

	FunctionParameters struct {
		Opening file.Loc
		List    []FunctionParameter
		Closing file.Loc
	}

	RegExpLiteral struct {
		Start   file.Loc
		Literal string
		Pattern string
		Flags   string
	}

	SequenceExpression struct {
		Sequence []Expression
	}

	StringLiteral struct {
		Start   file.Loc
		Literal string
		Value   string
	}

	ThisExpression struct {
		Start   file.Loc
		Private bool
	}

	UnaryExpression struct {
		Operator token.Token
		Start    file.Loc // If a prefix operation
		Operand  Expression
		Postfix  bool
	}

	VariableExpression struct {
		Kind        token.Token
		Name        string
		Start       file.Loc
		Initializer Expression
		FlowType    FlowType
	}

	ArrowFunctionExpression struct {
		Start          file.Loc
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     []FunctionParameter
		Body           Statement
		Async          bool
	}

	AwaitExpression struct {
		Start      file.Loc
		Expression Expression
	}
)

// _expressionNode

func (*AssignExpression) _expressionNode()        {}
func (*BinaryExpression) _expressionNode()        {}
func (*BooleanLiteral) _expressionNode()          {}
func (*BracketExpression) _expressionNode()       {}
func (*CallExpression) _expressionNode()          {}
func (*DotExpression) _expressionNode()           {}
func (*Identifier) _expressionNode()              {}
func (*NewExpression) _expressionNode()           {}
func (*NullLiteral) _expressionNode()             {}
func (*NumberLiteral) _expressionNode()           {}
func (*RegExpLiteral) _expressionNode()           {}
func (*SequenceExpression) _expressionNode()      {}
func (*StringLiteral) _expressionNode()           {}
func (*ThisExpression) _expressionNode()          {}
func (*UnaryExpression) _expressionNode()         {}
func (*VariableExpression) _expressionNode()      {}
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
		LeftBrace  file.Loc
		List       []Statement
		RightBrace file.Loc
	}

	BranchStatement struct {
		Idx   file.Loc
		Token token.Token
		Label *Identifier
	}

	CaseStatement struct {
		Case       file.Loc
		Test       Expression
		Consequent []Statement
	}

	CatchStatement struct {
		Catch     file.Loc
		Parameter *Identifier
		Body      Statement
	}

	DebuggerStatement struct {
		Debugger file.Loc
	}

	DoWhileStatement struct {
		Do   file.Loc
		Test Expression
		Body Statement
	}

	EmptyStatement struct {
		Semicolon file.Loc
	}

	ExpressionStatement struct {
		Expression Expression
	}

	IfStatement struct {
		If         file.Loc
		Test       Expression
		Consequent Statement
		Alternate  Statement
	}

	LabelledStatement struct {
		Label     *Identifier
		Colon     file.Loc
		Statement Statement
	}

	ReturnStatement struct {
		Return   file.Loc
		Argument Expression
	}

	SwitchStatement struct {
		Switch       file.Loc
		Discriminant Expression
		Default      int
		Body         []*CaseStatement
	}

	ThrowStatement struct {
		Throw    file.Loc
		Argument Expression
	}

	TryStatement struct {
		Try     file.Loc
		Body    Statement
		Catch   *CatchStatement
		Finally Statement
	}

	VariableStatement struct {
		Kind token.Token
		Var  file.Loc
		List []Expression
	}

	WhileStatement struct {
		While file.Loc
		Test  Expression
		Body  Statement
	}

	WithStatement struct {
		With   file.Loc
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
		Kind token.Token
		Var  file.Loc
		List []*VariableExpression
	}

	ClassDeclaration struct {
		Name    *Identifier
		Extends MemberExpression
	}
)

// _declarationNode

func (*FunctionDeclaration) _declarationNode() {}
func (*VariableDeclaration) _declarationNode() {}
func (*ClassDeclaration) _declarationNode()    {}

// ==== //
// Node //
// ==== //

type Program struct {
	Body []Statement

	DeclarationList []Declaration

	File *file.File

	SourceMap *sourcemap.Consumer
}

// ==== //
// StartAt //
// ==== //

func (self *AssignExpression) StartAt() file.Loc        { return self.Left.StartAt() }
func (self *BinaryExpression) StartAt() file.Loc        { return self.Left.StartAt() }
func (self *BooleanLiteral) StartAt() file.Loc          { return self.Start }
func (self *BracketExpression) StartAt() file.Loc       { return self.Left.StartAt() }
func (self *CallExpression) StartAt() file.Loc          { return self.Callee.StartAt() }
func (self *DotExpression) StartAt() file.Loc           { return self.Left.StartAt() }
func (self *Identifier) StartAt() file.Loc              { return self.Start }
func (self *NewExpression) StartAt() file.Loc           { return self.Start }
func (self *NullLiteral) StartAt() file.Loc             { return self.Start }
func (self *NumberLiteral) StartAt() file.Loc           { return self.Start }
func (self *RegExpLiteral) StartAt() file.Loc           { return self.Start }
func (self *SequenceExpression) StartAt() file.Loc      { return self.Sequence[0].StartAt() }
func (self *StringLiteral) StartAt() file.Loc           { return self.Start }
func (self *ThisExpression) StartAt() file.Loc          { return self.Start }
func (self *UnaryExpression) StartAt() file.Loc         { return self.Start }
func (self *VariableExpression) StartAt() file.Loc      { return self.Start }
func (self *ArrowFunctionExpression) StartAt() file.Loc { return self.Start }
func (self *AwaitExpression) StartAt() file.Loc         { return self.Start }
func (self *SpreadExpression) StartAt() file.Loc        { return self.Start }

func (self *BlockStatement) StartAt() file.Loc      { return self.LeftBrace }
func (self *BranchStatement) StartAt() file.Loc     { return self.Idx }
func (self *CaseStatement) StartAt() file.Loc       { return self.Case }
func (self *CatchStatement) StartAt() file.Loc      { return self.Catch }
func (self *DebuggerStatement) StartAt() file.Loc   { return self.Debugger }
func (self *DoWhileStatement) StartAt() file.Loc    { return self.Do }
func (self *EmptyStatement) StartAt() file.Loc      { return self.Semicolon }
func (self *ExpressionStatement) StartAt() file.Loc { return self.Expression.StartAt() }
func (self *IfStatement) StartAt() file.Loc         { return self.If }
func (self *LabelledStatement) StartAt() file.Loc   { return self.Label.StartAt() }
func (self *Program) StartAt() file.Loc             { return self.Body[0].StartAt() }
func (self *ReturnStatement) StartAt() file.Loc     { return self.Return }
func (self *SwitchStatement) StartAt() file.Loc     { return self.Switch }
func (self *ThrowStatement) StartAt() file.Loc      { return self.Throw }
func (self *TryStatement) StartAt() file.Loc        { return self.Try }
func (self *VariableStatement) StartAt() file.Loc   { return self.Var }
func (self *WhileStatement) StartAt() file.Loc      { return self.While }
func (self *WithStatement) StartAt() file.Loc       { return self.With }

// ==== //
// EndAt //
// ==== //

func (self *AssignExpression) EndAt() file.Loc   { return self.Right.EndAt() }
func (self *BinaryExpression) EndAt() file.Loc   { return self.Right.EndAt() }
func (self *BooleanLiteral) EndAt() file.Loc     { return file.Loc(int(self.Start) + len(self.Literal)) }
func (self *BracketExpression) EndAt() file.Loc  { return self.RightBracket + 1 }
func (self *CallExpression) EndAt() file.Loc     { return self.RightParenthesis + 1 }
func (self *DotExpression) EndAt() file.Loc      { return self.Identifier.EndAt() }
func (self *Identifier) EndAt() file.Loc         { return file.Loc(int(self.Start) + len(self.Name)) }
func (self *NewExpression) EndAt() file.Loc      { return self.RightParenthesis + 1 }
func (self *NullLiteral) EndAt() file.Loc        { return file.Loc(int(self.Start) + 4) } // "null"
func (self *NumberLiteral) EndAt() file.Loc      { return file.Loc(int(self.Start) + len(self.Literal)) }
func (self *RegExpLiteral) EndAt() file.Loc      { return file.Loc(int(self.Start) + len(self.Literal)) }
func (self *SequenceExpression) EndAt() file.Loc { return self.Sequence[0].EndAt() }
func (self *StringLiteral) EndAt() file.Loc      { return file.Loc(int(self.Start) + len(self.Literal)) }
func (self *ThisExpression) EndAt() file.Loc     { return self.Start }
func (self *UnaryExpression) EndAt() file.Loc {
	if self.Postfix {
		return self.Operand.EndAt() + 2 // ++ --
	}
	return self.Operand.EndAt()
}
func (self *VariableExpression) EndAt() file.Loc {
	if self.Initializer == nil {
		return file.Loc(int(self.Start) + len(self.Name) + 1)
	}
	return self.Initializer.EndAt()
}
func (self *ArrowFunctionExpression) EndAt() file.Loc { return self.Body.EndAt() }
func (self *AwaitExpression) EndAt() file.Loc         { return self.Expression.EndAt() }
func (self *SpreadExpression) EndAt() file.Loc        { return self.Value.EndAt() }

func (self *BlockStatement) EndAt() file.Loc      { return self.RightBrace + 1 }
func (self *BranchStatement) EndAt() file.Loc     { return self.Idx }
func (self *CaseStatement) EndAt() file.Loc       { return self.Consequent[len(self.Consequent)-1].EndAt() }
func (self *CatchStatement) EndAt() file.Loc      { return self.Body.EndAt() }
func (self *DebuggerStatement) EndAt() file.Loc   { return self.Debugger + 8 }
func (self *DoWhileStatement) EndAt() file.Loc    { return self.Test.EndAt() }
func (self *EmptyStatement) EndAt() file.Loc      { return self.Semicolon + 1 }
func (self *ExpressionStatement) EndAt() file.Loc { return self.Expression.EndAt() }
func (self *IfStatement) EndAt() file.Loc {
	if self.Alternate != nil {
		return self.Alternate.EndAt()
	}
	return self.Consequent.EndAt()
}
func (self *LabelledStatement) EndAt() file.Loc { return self.Colon + 1 }
func (self *Program) EndAt() file.Loc           { return self.Body[len(self.Body)-1].EndAt() }
func (self *ReturnStatement) EndAt() file.Loc   { return self.Return }
func (self *SwitchStatement) EndAt() file.Loc   { return self.Body[len(self.Body)-1].EndAt() }
func (self *ThrowStatement) EndAt() file.Loc    { return self.Throw }
func (self *TryStatement) EndAt() file.Loc      { return self.Try }
func (self *VariableStatement) EndAt() file.Loc { return self.List[len(self.List)-1].EndAt() }
func (self *WhileStatement) EndAt() file.Loc    { return self.Body.EndAt() }
func (self *WithStatement) EndAt() file.Loc     { return self.Body.EndAt() }
