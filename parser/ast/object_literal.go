package ast

type (
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
		Value        IExpr
	}

	ObjectPropertyName interface {
		IExpr
		_objectPropertyName()
	}
)

func (*ObjectPropertySetter) _objectProperty() {}
func (*ObjectPropertyGetter) _objectProperty() {}
func (*ObjectPropertyValue) _objectProperty()  {}
func (*ObjectSpread) _objectProperty()         {}

func (*Identifier) _objectPropertyName()    {}
func (*ComputedName) _objectPropertyName()  {}
func (*StringLiteral) _objectPropertyName() {}
func (*NumberLiteral) _objectPropertyName() {}
