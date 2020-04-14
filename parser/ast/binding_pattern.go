package ast

import "yawp/parser/file"

type (
	PatternBinder interface {
		_patternBinder()
	}

	IdentifierBinder struct {
		Name *Identifier
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
		Start file.Loc
		End   file.Loc
		List  []PatternBinder
	}

	ObjectBinding struct {
		Start file.Loc
		End   file.Loc
		List  []PatternBinder
	}

	VariableBinding struct {
		Start       file.Loc
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

func (s *ArrayBinding) StartAt() file.Loc    { return s.Start }
func (s *ObjectBinding) StartAt() file.Loc   { return s.Start }
func (s *VariableBinding) StartAt() file.Loc { return s.Start }

func (s *ArrayBinding) EndAt() file.Loc    { return s.End }
func (s *ObjectBinding) EndAt() file.Loc   { return s.End }
func (s *VariableBinding) EndAt() file.Loc { return s.Initializer.EndAt() }
