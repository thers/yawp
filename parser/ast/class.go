package ast

import "yawp/parser/file"

type (
	ClassStatement struct {
		Expression *ClassExpression
	}

	ClassFieldName interface {
		_classFieldName()
	}

	ClassFieldStatement struct {
		Start       file.Idx
		Name        ClassFieldName
		Static      bool
		Private     bool
		Initializer Expression
	}

	ClassAccessorStatement struct {
		Start file.Idx
		Field ClassFieldName
		Kind  string
		Body  *FunctionLiteral
	}

	ClassMethodStatement struct {
		Method     file.Idx
		Name       ClassFieldName
		Static     bool
		Private    bool
		Async      bool
		Generator  bool
		Parameters *FunctionParameters
		Body       Statement
		Source     string
	}
)

func (*ClassStatement) _statementNode()         {}
func (*ClassFieldStatement) _statementNode()    {}
func (*ClassAccessorStatement) _statementNode() {}
func (*ClassMethodStatement) _statementNode()   {}

func (*Identifier) _classFieldName()   {}
func (*ComputedName) _classFieldName() {}

func (self *ClassStatement) StartAt() file.Idx         { return self.Expression.Start }
func (self *ClassFieldStatement) StartAt() file.Idx    { return self.Start }
func (self *ClassAccessorStatement) StartAt() file.Idx { return self.Start }
func (self *ClassMethodStatement) StartAt() file.Idx   { return self.Method }

func (self *ClassStatement) EndAt() file.Idx         { return self.Expression.Body.EndAt() }
func (self *ClassFieldStatement) EndAt() file.Idx    { return self.Initializer.EndAt() }
func (self *ClassMethodStatement) EndAt() file.Idx   { return self.Body.EndAt() }
func (self *ClassAccessorStatement) EndAt() file.Idx { return self.Body.EndAt() }
