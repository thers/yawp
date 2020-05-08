package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type (
	PatternBinder interface {
		_patternBinder()
	}

	ExpressionBinder struct {
		Expression Expression
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
		DefaultValue Expression
	}

	ArrayItemBinder struct {
		Binder       PatternBinder
		Index        int
		DefaultValue Expression
	}

	ArrayBinding struct {
		Loc  *file.Loc
		List []PatternBinder
	}

	ObjectBinding struct {
		Loc  *file.Loc
		List []PatternBinder
	}

	VariableBinding struct {
		Loc         *file.Loc
		Kind        token.Token
		Binder      PatternBinder
		Initializer Expression
		FlowType    FlowType
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

func (*ArrayBinding) _expressionNode()    {}
func (*ObjectBinding) _expressionNode()   {}
func (*VariableBinding) _expressionNode() {}

func (s *ArrayBinding) GetLoc() *file.Loc    { return s.Loc }
func (s *ObjectBinding) GetLoc() *file.Loc   { return s.Loc }
func (s *VariableBinding) GetLoc() *file.Loc { return s.Loc }

func (s *VariableBinding) Clone() *VariableBinding {
	return &VariableBinding{
		Loc:         s.Loc,
		Kind:        s.Kind,
		Binder:      s.Binder,
		Initializer: s.Initializer,
		FlowType:    s.FlowType,
	}
}
