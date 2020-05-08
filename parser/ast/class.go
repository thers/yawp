package ast

import "yawp/parser/file"

type (
	ClassStatement struct {
		Expression *ClassExpression
	}

	ClassExpression struct {
		Loc                *file.Loc
		Name               *Identifier
		TypeParameters     []*FlowTypeParameter
		SuperClass         Expression
		SuperTypeArguments []FlowType
		Body               Statement
	}

	SuperExpression struct {
		Loc *file.Loc
	}

	ClassFieldName interface {
		_classFieldName()
	}

	ClassFieldStatement struct {
		Loc         *file.Loc
		Name        ClassFieldName
		Static      bool
		Private     bool
		Type        FlowType
		Initializer Expression
	}

	ClassAccessorStatement struct {
		Loc   *file.Loc
		Field ClassFieldName
		Kind  string
		Body  *FunctionLiteral
	}

	ClassMethodStatement struct {
		Loc            *file.Loc
		Name           ClassFieldName
		Static         bool
		Private        bool
		Async          bool
		Generator      bool
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     *FunctionParameters
		Body           Statement
	}
)

func (*ClassStatement) _statementNode()         {}
func (*ClassFieldStatement) _statementNode()    {}
func (*ClassAccessorStatement) _statementNode() {}
func (*ClassMethodStatement) _statementNode()   {}

func (*ClassExpression) _expressionNode() {}
func (*SuperExpression) _expressionNode() {}

func (*Identifier) _classFieldName()    {}
func (*StringLiteral) _classFieldName() {}
func (*ComputedName) _classFieldName()  {}

func (c *ClassStatement) GetLoc() *file.Loc         { return c.Expression.GetLoc() }
func (c *ClassExpression) GetLoc() *file.Loc        { return c.Loc }
func (c *SuperExpression) GetLoc() *file.Loc        { return c.Loc }
func (c *ClassFieldStatement) GetLoc() *file.Loc    { return c.Loc }
func (c *ClassAccessorStatement) GetLoc() *file.Loc { return c.Loc }
func (c *ClassMethodStatement) GetLoc() *file.Loc   { return c.Loc }
