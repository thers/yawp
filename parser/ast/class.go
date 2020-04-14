package ast

import "yawp/parser/file"

type (
	ClassStatement struct {
		Expression *ClassExpression
	}

	ClassExpression struct {
		Start              file.Loc
		Name               *Identifier
		TypeParameters     []*FlowTypeParameter
		SuperClass         MemberExpression
		SuperTypeArguments []FlowType
		Body               Statement
	}

	ClassSuperExpression struct {
		Start     file.Loc
		End       file.Loc
		Arguments []Expression
	}

	ClassFieldName interface {
		_classFieldName()
	}

	ClassFieldStatement struct {
		Start       file.Loc
		Name        ClassFieldName
		Static      bool
		Private     bool
		Type        FlowType
		Initializer Expression
	}

	ClassAccessorStatement struct {
		Start file.Loc
		Field ClassFieldName
		Kind  string
		Body  *FunctionLiteral
	}

	ClassMethodStatement struct {
		Method         file.Loc
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

func (*ClassExpression) _expressionNode()      {}
func (*ClassSuperExpression) _expressionNode() {}

func (*Identifier) _classFieldName()   {}
func (*ComputedName) _classFieldName() {}

func (self *ClassStatement) StartAt() file.Loc         { return self.Expression.Start }
func (self *ClassExpression) StartAt() file.Loc        { return self.Start }
func (self *ClassSuperExpression) StartAt() file.Loc   { return self.Start }
func (self *ClassFieldStatement) StartAt() file.Loc    { return self.Start }
func (self *ClassAccessorStatement) StartAt() file.Loc { return self.Start }
func (self *ClassMethodStatement) StartAt() file.Loc   { return self.Method }

func (self *ClassStatement) EndAt() file.Loc         { return self.Expression.Body.EndAt() }
func (self *ClassExpression) EndAt() file.Loc        { return self.Body.EndAt() }
func (self *ClassSuperExpression) EndAt() file.Loc   { return self.End }
func (self *ClassFieldStatement) EndAt() file.Loc    { return self.Initializer.EndAt() }
func (self *ClassMethodStatement) EndAt() file.Loc   { return self.Body.EndAt() }
func (self *ClassAccessorStatement) EndAt() file.Loc { return self.Body.EndAt() }
