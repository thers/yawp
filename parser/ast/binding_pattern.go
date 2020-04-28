package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type (
	PatternBinder interface {
		_patternBinder()
	}

	IdentifierBinder struct {
		Id *Identifier
	}

	ObjectRestBinder struct {
		Name *Identifier
	}

	ArrayRestBinder struct {
		Name *Identifier
	}

	ObjectPropertyBinder struct {
		Property     PatternBinder
		PropertyName *Identifier
		DefaultValue Expression
	}

	ArrayItemBinder struct {
		Item         PatternBinder
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

func (*ArrayBinding) _expressionNode()    {}
func (*ObjectBinding) _expressionNode()   {}
func (*VariableBinding) _expressionNode() {}

func (s *ArrayBinding) GetLoc() *file.Loc    { return s.Loc }
func (s *ObjectBinding) GetLoc() *file.Loc   { return s.Loc }
func (s *VariableBinding) GetLoc() *file.Loc { return s.Loc }
