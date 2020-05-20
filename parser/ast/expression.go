package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

// All expression nodes implement the IExpr interface.
type IExpr interface {
	INode
	_expressionNode()
}
type ExprNode Node

func (*ExprNode) _expressionNode() {}

func (e *ExprNode) GetLoc() *file.Loc { return e.Loc }
func (e *ExprNode) GetNode() *Node    { return (*Node)(e) }

func (e *ExprNode) Copy() ExprNode {
	return (ExprNode)((*Node)(e).Copy())
}

type Expressions []IExpr

func (Expressions) _expressionNode()  {}
func (Expressions) GetLoc() *file.Loc { panic("Can not invoke GetLoc on Expressions") }
func (Expressions) GetNode() *Node    { panic("Can not invoke GetNode on Expressions") }

type (
	AssignmentExpression struct {
		ExprNode
		Operator token.Token
		Left     IExpr
		Right    IExpr
	}

	BinaryExpression struct {
		ExprNode
		Operator   token.Token
		Left       IExpr
		Right      IExpr
		Comparison bool
	}

	BooleanLiteral struct {
		ExprNode
		Literal string
	}

	CallExpression struct {
		ExprNode
		Callee        IExpr
		TypeArguments []FlowType
		ArgumentList  []IExpr
	}

	SpreadExpression struct {
		ExprNode
		Value IExpr
	}

	NewExpression struct {
		ExprNode
		Callee        IExpr
		TypeArguments []FlowType
		ArgumentList  []IExpr
	}

	NullLiteral struct {
		ExprNode
		Literal string
	}

	NumberLiteral struct {
		ExprNode
		Literal string
	}

	RegExpLiteral struct {
		ExprNode
		Literal string
		Pattern string
		Flags   string
	}

	SequenceExpression struct {
		ExprNode
		Sequence []IExpr
	}

	StringLiteral struct {
		ExprNode
		Raw     bool
		Literal string
	}

	ThisExpression struct {
		ExprNode
		Private bool
	}

	UnaryExpression struct {
		ExprNode
		Operator token.Token
		Operand  IExpr
		Postfix  bool
	}

	ArrowFunctionExpression struct {
		ExprNode
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     []FunctionParameter
		Body           *FunctionBody
		Async          bool
	}

	AwaitExpression struct {
		ExprNode
		Expression IExpr
	}

	NewTargetExpression struct {
		ExprNode
	}

	ArrayLiteral struct {
		ExprNode
		List []IExpr
	}

	ArraySpread struct {
		ExprNode
		Expression IExpr
	}

	CoalesceExpression struct {
		ExprNode
		Head       IExpr
		Consequent IExpr
	}

	OptionalObjectMemberAccessExpression struct {
		ExprNode
		Left       IExpr
		Identifier *Identifier
	}

	OptionalArrayMemberAccessExpression struct {
		ExprNode
		Left  IExpr
		Index IExpr
	}

	OptionalCallExpression struct {
		ExprNode
		Left      IExpr
		Arguments []IExpr
	}

	ConditionalExpression struct {
		ExprNode
		Test       IExpr
		Consequent IExpr
		Alternate  IExpr
	}

	TemplateExpression struct {
		ExprNode
		Strings       []string
		Substitutions []IExpr
	}

	TaggedTemplateExpression struct {
		ExprNode
		Tag      IExpr
		Template *TemplateExpression
	}

	YieldExpression struct {
		ExprNode
		Delegate bool
		Argument IExpr
	}

	ClassExpression struct {
		ExprNode
		Name               *Identifier
		TypeParameters     []*FlowTypeParameter
		SuperClass         IExpr
		SuperTypeArguments []FlowType
		Body               IStmt
	}

	SuperExpression struct {
		ExprNode
	}

	ComputedName struct {
		ExprNode
		Expression IExpr
	}

	ArrayBinding struct {
		ExprNode
		List []PatternBinder
	}

	ObjectBinding struct {
		ExprNode
		List []PatternBinder
	}

	VariableBinding struct {
		ExprNode
		Kind        token.Token
		Binder      PatternBinder
		Initializer IExpr
		FlowType    FlowType
	}

	FlowTypeAssertionExpression struct {
		ExprNode
		Left     IExpr
		FlowType FlowType
	}

	JSXElement struct {
		ExprNode
		Name       *JSXElementName
		Attributes []JSXAttribute
		Children   []JSXChild
	}

	JSXFragment struct {
		ExprNode
		Children []JSXChild
	}

	JSXNamespacedName struct {
		ExprNode
		Namespace string
		Name      string
	}

	MemberExpression struct {
		ExprNode
		Left  IExpr
		Right IExpr
		Kind  MemberKind
	}

	ObjectSpread struct {
		ExprNode
		Expression IExpr
	}

	ObjectLiteral struct {
		ExprNode
		Properties []ObjectProperty
	}
)
