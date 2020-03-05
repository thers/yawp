package ast

import "yawp/parser/file"

type (
	PatternBinder interface {
		_patternBinder()
	}

	IdentifierBinder struct {
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
		Start file.Idx
		End   file.Idx
		List  []*ArrayItemBinder
	}

	ObjectBinding struct {
		Start file.Idx
		End   file.Idx
		List  []*ObjectPropertyBinder
	}

	VariableBinding struct {
		Start       file.Idx
		Binder      PatternBinder
		Initializer Expression
	}
)

func (*ArrayBinding) _patternBinder()         {}
func (*ObjectBinding) _patternBinder()        {}
func (*IdentifierBinder) _patternBinder()     {}
func (*ObjectPropertyBinder) _patternBinder() {}
func (*ArrayItemBinder) _patternBinder()      {}

func (*ArrayBinding) _expressionNode()    {}
func (*ObjectBinding) _expressionNode()   {}
func (*VariableBinding) _expressionNode() {}

func (s *ArrayBinding) StartAt() file.Idx    { return s.Start }
func (s *ObjectBinding) StartAt() file.Idx   { return s.Start }
func (s *VariableBinding) StartAt() file.Idx { return s.Start }

func (s *ArrayBinding) EndAt() file.Idx    { return s.End }
func (s *ObjectBinding) EndAt() file.Idx   { return s.End }
func (s *VariableBinding) EndAt() file.Idx { return s.Initializer.EndAt() }
