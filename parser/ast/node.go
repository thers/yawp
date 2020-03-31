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
	StartAt() file.Idx // The index of the first character belonging to the node
	EndAt() file.Idx   // The index of the first character immediately after the node
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

	BadExpression struct {
		From file.Idx
		To   file.Idx
	}

	BinaryExpression struct {
		Operator   token.Token
		Left       Expression
		Right      Expression
		Comparison bool
	}

	BooleanLiteral struct {
		Start   file.Idx
		Literal string
		Value   bool
	}

	BracketExpression struct {
		Left         Expression
		Member       Expression
		LeftBracket  file.Idx
		RightBracket file.Idx
	}

	CallExpression struct {
		Callee           Expression
		LeftParenthesis  file.Idx
		ArgumentList     []Expression
		RightParenthesis file.Idx
	}

	DotExpression struct {
		Left       Expression
		Identifier *Identifier
	}

	Identifier struct {
		Start file.Idx
		Name  string
	}

	ObjectDestructuring struct {
		Start       file.Idx
		Identifiers []*Identifier
	}

	RestExpression struct {
		Start file.Idx
		Name  string
	}

	Spread struct {
		Start file.Idx
		Value Expression
	}

	NewExpression struct {
		Start            file.Idx
		Callee           Expression
		LeftParenthesis  file.Idx
		ArgumentList     []Expression
		RightParenthesis file.Idx
	}

	NullLiteral struct {
		Start   file.Idx
		Literal string
	}

	NumberLiteral struct {
		Start   file.Idx
		Literal string
		Value   interface{}
	}

	FunctionParameters struct {
		Opening file.Idx
		List    []FunctionParameter
		Closing file.Idx
	}

	RegExpLiteral struct {
		Start   file.Idx
		Literal string
		Pattern string
		Flags   string
	}

	SequenceExpression struct {
		Sequence []Expression
	}

	StringLiteral struct {
		Start   file.Idx
		Literal string
		Value   string
	}

	ThisExpression struct {
		Start   file.Idx
		Private bool
	}

	UnaryExpression struct {
		Operator token.Token
		Start    file.Idx // If a prefix operation
		Operand  Expression
		Postfix  bool
	}

	VariableExpression struct {
		Kind        token.Token
		Name        string
		Start       file.Idx
		Initializer Expression
		FlowType    FlowType
	}

	ArrowFunctionExpression struct {
		Start          file.Idx
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     []FunctionParameter
		Body           Statement
		Async          bool
	}

	ClassExpression struct {
		Start   file.Idx
		Name    *Identifier
		Extends *Identifier
		Body    Statement
	}

	ClassSuperExpression struct {
		Start     file.Idx
		End       file.Idx
		Arguments []Expression
	}

	AwaitExpression struct {
		Start      file.Idx
		Expression Expression
	}
)

// _expressionNode

func (*AssignExpression) _expressionNode()        {}
func (*BadExpression) _expressionNode()           {}
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
func (*ClassExpression) _expressionNode()         {}
func (*ClassSuperExpression) _expressionNode()    {}
func (*AwaitExpression) _expressionNode()         {}

// ========= //
// Statement //
// ========= //

type (
	// All statement nodes implement the Statement interface.
	Statement interface {
		Node
		_statementNode()
	}

	BadStatement struct {
		From file.Idx
		To   file.Idx
	}

	BlockStatement struct {
		LeftBrace  file.Idx
		List       []Statement
		RightBrace file.Idx
	}

	BranchStatement struct {
		Idx   file.Idx
		Token token.Token
		Label *Identifier
	}

	CaseStatement struct {
		Case       file.Idx
		Test       Expression
		Consequent []Statement
	}

	CatchStatement struct {
		Catch     file.Idx
		Parameter *Identifier
		Body      Statement
	}

	DebuggerStatement struct {
		Debugger file.Idx
	}

	DoWhileStatement struct {
		Do   file.Idx
		Test Expression
		Body Statement
	}

	EmptyStatement struct {
		Semicolon file.Idx
	}

	ExpressionStatement struct {
		Expression Expression
	}

	IfStatement struct {
		If         file.Idx
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
		Return   file.Idx
		Argument Expression
	}

	SwitchStatement struct {
		Switch       file.Idx
		Discriminant Expression
		Default      int
		Body         []*CaseStatement
	}

	ThrowStatement struct {
		Throw    file.Idx
		Argument Expression
	}

	TryStatement struct {
		Try     file.Idx
		Body    Statement
		Catch   *CatchStatement
		Finally Statement
	}

	VariableStatement struct {
		Kind token.Token
		Var  file.Idx
		List []Expression
	}

	WhileStatement struct {
		While file.Idx
		Test  Expression
		Body  Statement
	}

	WithStatement struct {
		With   file.Idx
		Object Expression
		Body   Statement
	}
)

// _statementNode

func (*BadStatement) _statementNode()        {}
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
		Var  file.Idx
		List []*VariableExpression
	}

	ClassDeclaration struct {
		Name    *Identifier
		Extends *Identifier
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

func (self *AssignExpression) StartAt() file.Idx        { return self.Left.StartAt() }
func (self *BadExpression) StartAt() file.Idx           { return self.From }
func (self *BinaryExpression) StartAt() file.Idx        { return self.Left.StartAt() }
func (self *BooleanLiteral) StartAt() file.Idx          { return self.Start }
func (self *BracketExpression) StartAt() file.Idx       { return self.Left.StartAt() }
func (self *CallExpression) StartAt() file.Idx          { return self.Callee.StartAt() }
func (self *DotExpression) StartAt() file.Idx           { return self.Left.StartAt() }
func (self *Identifier) StartAt() file.Idx              { return self.Start }
func (self *NewExpression) StartAt() file.Idx           { return self.Start }
func (self *NullLiteral) StartAt() file.Idx             { return self.Start }
func (self *NumberLiteral) StartAt() file.Idx           { return self.Start }
func (self *RegExpLiteral) StartAt() file.Idx           { return self.Start }
func (self *SequenceExpression) StartAt() file.Idx      { return self.Sequence[0].StartAt() }
func (self *StringLiteral) StartAt() file.Idx           { return self.Start }
func (self *ThisExpression) StartAt() file.Idx          { return self.Start }
func (self *UnaryExpression) StartAt() file.Idx         { return self.Start }
func (self *VariableExpression) StartAt() file.Idx      { return self.Start }
func (self *ArrowFunctionExpression) StartAt() file.Idx { return self.Start }
func (self *ClassExpression) StartAt() file.Idx         { return self.Start }
func (self *ClassSuperExpression) StartAt() file.Idx    { return self.Start }
func (self *AwaitExpression) StartAt() file.Idx         { return self.Start }

func (self *BadStatement) StartAt() file.Idx        { return self.From }
func (self *BlockStatement) StartAt() file.Idx      { return self.LeftBrace }
func (self *BranchStatement) StartAt() file.Idx     { return self.Idx }
func (self *CaseStatement) StartAt() file.Idx       { return self.Case }
func (self *CatchStatement) StartAt() file.Idx      { return self.Catch }
func (self *DebuggerStatement) StartAt() file.Idx   { return self.Debugger }
func (self *DoWhileStatement) StartAt() file.Idx    { return self.Do }
func (self *EmptyStatement) StartAt() file.Idx      { return self.Semicolon }
func (self *ExpressionStatement) StartAt() file.Idx { return self.Expression.StartAt() }
func (self *IfStatement) StartAt() file.Idx         { return self.If }
func (self *LabelledStatement) StartAt() file.Idx   { return self.Label.StartAt() }
func (self *Program) StartAt() file.Idx             { return self.Body[0].StartAt() }
func (self *ReturnStatement) StartAt() file.Idx     { return self.Return }
func (self *SwitchStatement) StartAt() file.Idx     { return self.Switch }
func (self *ThrowStatement) StartAt() file.Idx      { return self.Throw }
func (self *TryStatement) StartAt() file.Idx        { return self.Try }
func (self *VariableStatement) StartAt() file.Idx   { return self.Var }
func (self *WhileStatement) StartAt() file.Idx      { return self.While }
func (self *WithStatement) StartAt() file.Idx       { return self.With }

// ==== //
// EndAt //
// ==== //

func (self *AssignExpression) EndAt() file.Idx   { return self.Right.EndAt() }
func (self *BadExpression) EndAt() file.Idx      { return self.To }
func (self *BinaryExpression) EndAt() file.Idx   { return self.Right.EndAt() }
func (self *BooleanLiteral) EndAt() file.Idx     { return file.Idx(int(self.Start) + len(self.Literal)) }
func (self *BracketExpression) EndAt() file.Idx  { return self.RightBracket + 1 }
func (self *CallExpression) EndAt() file.Idx     { return self.RightParenthesis + 1 }
func (self *DotExpression) EndAt() file.Idx      { return self.Identifier.EndAt() }
func (self *Identifier) EndAt() file.Idx         { return file.Idx(int(self.Start) + len(self.Name)) }
func (self *NewExpression) EndAt() file.Idx      { return self.RightParenthesis + 1 }
func (self *NullLiteral) EndAt() file.Idx        { return file.Idx(int(self.Start) + 4) } // "null"
func (self *NumberLiteral) EndAt() file.Idx      { return file.Idx(int(self.Start) + len(self.Literal)) }
func (self *RegExpLiteral) EndAt() file.Idx      { return file.Idx(int(self.Start) + len(self.Literal)) }
func (self *SequenceExpression) EndAt() file.Idx { return self.Sequence[0].EndAt() }
func (self *StringLiteral) EndAt() file.Idx      { return file.Idx(int(self.Start) + len(self.Literal)) }
func (self *ThisExpression) EndAt() file.Idx     { return self.Start }
func (self *UnaryExpression) EndAt() file.Idx {
	if self.Postfix {
		return self.Operand.EndAt() + 2 // ++ --
	}
	return self.Operand.EndAt()
}
func (self *VariableExpression) EndAt() file.Idx {
	if self.Initializer == nil {
		return file.Idx(int(self.Start) + len(self.Name) + 1)
	}
	return self.Initializer.EndAt()
}
func (self *ArrowFunctionExpression) EndAt() file.Idx { return self.Body.EndAt() }
func (self *ClassExpression) EndAt() file.Idx         { return self.Body.EndAt() }
func (self *ClassSuperExpression) EndAt() file.Idx    { return self.End }
func (self *AwaitExpression) EndAt() file.Idx         { return self.Expression.EndAt() }

func (self *BadStatement) EndAt() file.Idx        { return self.To }
func (self *BlockStatement) EndAt() file.Idx      { return self.RightBrace + 1 }
func (self *BranchStatement) EndAt() file.Idx     { return self.Idx }
func (self *CaseStatement) EndAt() file.Idx       { return self.Consequent[len(self.Consequent)-1].EndAt() }
func (self *CatchStatement) EndAt() file.Idx      { return self.Body.EndAt() }
func (self *DebuggerStatement) EndAt() file.Idx   { return self.Debugger + 8 }
func (self *DoWhileStatement) EndAt() file.Idx    { return self.Test.EndAt() }
func (self *EmptyStatement) EndAt() file.Idx      { return self.Semicolon + 1 }
func (self *ExpressionStatement) EndAt() file.Idx { return self.Expression.EndAt() }
func (self *IfStatement) EndAt() file.Idx {
	if self.Alternate != nil {
		return self.Alternate.EndAt()
	}
	return self.Consequent.EndAt()
}
func (self *LabelledStatement) EndAt() file.Idx { return self.Colon + 1 }
func (self *Program) EndAt() file.Idx           { return self.Body[len(self.Body)-1].EndAt() }
func (self *ReturnStatement) EndAt() file.Idx   { return self.Return }
func (self *SwitchStatement) EndAt() file.Idx   { return self.Body[len(self.Body)-1].EndAt() }
func (self *ThrowStatement) EndAt() file.Idx    { return self.Throw }
func (self *TryStatement) EndAt() file.Idx      { return self.Try }
func (self *VariableStatement) EndAt() file.Idx { return self.List[len(self.List)-1].EndAt() }
func (self *WhileStatement) EndAt() file.Idx    { return self.Body.EndAt() }
func (self *WithStatement) EndAt() file.Idx     { return self.Body.EndAt() }
