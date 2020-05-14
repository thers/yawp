package ast

type (
	PatternBinder interface {
		_patternBinder()
	}

	ExpressionBinder struct {
		Expression IExpr
	}

	IdentifierBinder struct {
		Id *Identifier
	}

	ObjectRestBinder struct {
		Binder         PatternBinder
		OmitProperties []ObjectPropertyName
	}

	ArrayRestBinder struct {
		Binder    PatternBinder
		FromIndex int
	}

	ObjectPropertyBinder struct {
		Binder       PatternBinder
		PropertyName ObjectPropertyName
		DefaultValue IExpr
	}

	ArrayItemBinder struct {
		Binder       PatternBinder
		Index        int
		DefaultValue IExpr
	}
)

func (*ArrayBinding) _patternBinder()         {}
func (*ObjectBinding) _patternBinder()        {}
func (*IdentifierBinder) _patternBinder()     {}
func (*ObjectRestBinder) _patternBinder()     {}
func (*ArrayRestBinder) _patternBinder()      {}
func (*ObjectPropertyBinder) _patternBinder() {}
func (*ArrayItemBinder) _patternBinder()      {}
func (*ExpressionBinder) _patternBinder()     {}

func (s *VariableBinding) Copy() *VariableBinding {
	return &VariableBinding{
		ExprNode:    s.ExprNode.Copy(),
		Kind:        s.Kind,
		Binder:      s.Binder,
		Initializer: s.Initializer,
		FlowType:    s.FlowType,
	}
}
