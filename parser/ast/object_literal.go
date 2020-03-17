package ast

import "yawp/parser/file"

type (
	Property struct {
		Key   string
		Value Expression
	}

	ObjectLiteral struct {
		LeftBrace  file.Idx
		RightBrace file.Idx
		Properties []ObjectProperty
	}

	ObjectProperty interface {
		_objectProperty()
	}

	ObjectPropertySetter struct {
		PropertyName ObjectPropertyName
		Setter       *FunctionLiteral
	}

	ObjectPropertyGetter struct {
		PropertyName ObjectPropertyName
		Getter       *FunctionLiteral
	}

	ObjectPropertyValue struct {
		PropertyName ObjectPropertyName
		Value        Expression
	}

	ObjectSpread struct {
		Expression
	}

	ObjectPropertyName interface {
		_objectPropertyName()
	}
)

func (*ObjectLiteral) _expressionNode() {}

func (*ObjectPropertySetter) _objectProperty() {}
func (*ObjectPropertyGetter) _objectProperty() {}
func (*ObjectPropertyValue) _objectProperty()  {}
func (*ObjectSpread) _objectProperty()         {}

func (*Identifier) _objectPropertyName()   {}
func (*ComputedName) _objectPropertyName() {}

func (o *ObjectLiteral) StartAt() file.Idx { return o.LeftBrace }

func (o *ObjectLiteral) EndAt() file.Idx { return o.RightBrace }
