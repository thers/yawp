package ast

import "yawp/parser/file"

type (
	Property struct {
		Key   string
		Value Expression
	}

	ObjectLiteral struct {
		Loc        *file.Loc
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
		Expression
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

func (o *ObjectLiteral) GetLoc() *file.Loc { return o.Loc }
