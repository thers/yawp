package ast

import "yawp/parser/file"

type (
	FlowTypeStatement struct {
		Name FlowIdentifier
		Type FlowType
	}

	FlowInterfaceStatement struct {
		Name FlowIdentifier
		Body []FlowInterfaceBodyStatement
	}

	FlowType interface {
		_flowType()
		EndAt() file.Idx
	}

	FlowInterfaceBodyStatement interface {
		_flowInterfaceBodyStatement()
	}

	FlowObjectProperty interface {
		_flowObjectProperty()
	}

	FlowIdentifier struct {
		Start file.Idx
		Name  string
	}

	FlowObjectIndexer struct {
		Key   FlowType
		Value FlowType
	}

	FlowInexactObject struct {
		Start      file.Idx
		Properties []FlowObjectProperty
	}

	FlowExactObject struct {
		Start      file.Idx
		Properties []FlowObjectProperty
	}

	FlowArrayType struct {
		Start    file.Idx
		Elements []FlowType
	}

	FlowUnionType struct {
		Start file.Idx
		Types []FlowType
	}

	FlowIntersectionType struct {
		Start file.Idx
		Types []FlowType
	}

	FlowTypeOfType struct {
		Start      file.Idx
		Identifier *FlowIdentifier
	}

	FlowTypeAssertion struct {
		Left Expression
		FlowType
	}

	FlowOptionalType struct {
		FlowType
	}
)

type (
	// Primitives
	FlowTrueType struct {
		Start file.Idx
		End   file.Idx
	}
	FlowFalseType struct {
		Start file.Idx
		End   file.Idx
	}
	FlowBooleanType struct {
		Start file.Idx
		End   file.Idx
	}
	FlowVoidType struct {
		Start file.Idx
		End   file.Idx
	}
	FlowNullType struct {
		Start file.Idx
		End   file.Idx
	}

	FlowStringType struct {
		Start file.Idx
	}

	FlowNumberType struct {
		Start file.Idx
	}

	FlowStringLiteralType struct {
		Start  file.Idx
		End    file.Idx
		String string
	}

	FlowNumberLiteralType struct {
		Start  file.Idx
		End    file.Idx
		Number interface{}
	}
)

func (*FlowTrueType) _flowType()          {}
func (*FlowFalseType) _flowType()         {}
func (*FlowBooleanType) _flowType()       {}
func (*FlowStringType) _flowType()        {}
func (*FlowNumberType) _flowType()        {}
func (*FlowStringLiteralType) _flowType() {}
func (*FlowNumberLiteralType) _flowType() {}
func (*FlowIdentifier) _flowType()        {}
func (*FlowTypeOfType) _flowType()        {}
func (*FlowOptionalType) _flowType()      {}

func (*FlowTypeAssertion) _expressionNode() {}

func (f *FlowTypeAssertion) StartAt() file.Idx { return f.Left.StartAt() }

func (f *FlowTrueType) EndAt() file.Idx          { return f.End }
func (f *FlowFalseType) EndAt() file.Idx         { return f.End }
func (f *FlowBooleanType) EndAt() file.Idx       { return f.End }
func (f *FlowStringType) EndAt() file.Idx        { return f.Start + 6 }
func (f *FlowNumberType) EndAt() file.Idx        { return f.Start + 6 }
func (f *FlowStringLiteralType) EndAt() file.Idx { return f.End }
func (f *FlowNumberLiteralType) EndAt() file.Idx { return f.End }
func (f *FlowIdentifier) EndAt() file.Idx        { return f.Start + file.Idx(len(f.Name)) }
func (f *FlowTypeOfType) EndAt() file.Idx {
	return f.Start + f.Identifier.Start + file.Idx(len(f.Identifier.Name))
}
func (f *FlowTypeAssertion) EndAt() file.Idx { return f.FlowType.EndAt() }
func (f *FlowOptionalType) EndAt() file.Idx  { return f.FlowType.EndAt() }
