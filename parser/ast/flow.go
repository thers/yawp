package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type (
	FlowTypeStatement struct {
		Start          file.Idx
		Name           *FlowIdentifier
		Type           FlowType
		TypeParameters []*FlowTypeParameter
	}

	FlowInterfaceStatement struct {
		Start          file.Idx
		End            file.Idx
		Name           *FlowIdentifier
		TypeParameters []*FlowTypeParameter
		Body           []FlowInterfaceBodyStatement
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

	FlowTypeParameter struct {
		Start         file.Idx
		Name          *FlowIdentifier
		Covariant     bool
		Contravariant bool
		Boundary      FlowType
		DefaultValue  FlowType
	}

	FlowInterfaceMethod struct {
		Start          file.Idx
		Name           *FlowIdentifier
		TypeParameters []*FlowTypeParameter
		Parameters     []interface{}
		ReturnType     FlowType
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
		End        file.Idx
		Properties []FlowObjectProperty
	}

	FlowExactObject struct {
		Start      file.Idx
		End        file.Idx
		Properties []FlowObjectProperty
	}

	FlowNamedObjectProperty struct {
		Start         file.Idx
		Optional      bool
		Covariant     bool
		Contravariant bool
		Name          string
		Value         FlowType
	}

	FlowIndexerObjectProperty struct {
		Start   file.Idx
		KeyName string
		KeyType FlowType
		Value   FlowType
	}

	FlowInexactSpecifierProperty struct {
		Start file.Idx
	}

	FlowSpreadObjectProperty struct {
		Start    file.Idx
		FlowType FlowType
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
		Left     Expression
		FlowType FlowType
	}

	FlowOptionalType struct {
		FlowType FlowType
	}

	FlowFunctionType struct {
		Start      file.Idx
		Parameters []FlowType
		ReturnType FlowType
	}
)

type (
	// Primitives
	FlowPrimitiveType struct {
		Start file.Idx
		End   file.Idx
		Kind  token.Token
	}

	FlowTrueType struct {
		Start file.Idx
		End   file.Idx
	}
	FlowFalseType struct {
		Start file.Idx
		End   file.Idx
	}

	FlowExistentialType struct {
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

	FlowTupleType struct {
		Start    file.Idx
		End      file.Idx
		Elements []FlowType
	}
)

func (*FlowPrimitiveType) _flowType()     {}
func (*FlowTrueType) _flowType()          {}
func (*FlowFalseType) _flowType()         {}
func (*FlowStringLiteralType) _flowType() {}
func (*FlowNumberLiteralType) _flowType() {}
func (*FlowIdentifier) _flowType()        {}
func (*FlowTypeOfType) _flowType()        {}
func (*FlowOptionalType) _flowType()      {}
func (*FlowInexactObject) _flowType()     {}
func (*FlowExactObject) _flowType()       {}
func (*FlowTupleType) _flowType()         {}
func (*FlowExistentialType) _flowType()   {}
func (*FlowFunctionType) _flowType()      {}

func (*FlowNamedObjectProperty) _flowObjectProperty()      {}
func (*FlowIndexerObjectProperty) _flowObjectProperty()    {}
func (*FlowInexactSpecifierProperty) _flowObjectProperty() {}
func (*FlowSpreadObjectProperty) _flowObjectProperty()     {}

func (*FlowNamedObjectProperty) _flowInterfaceBodyStatement()   {}
func (*FlowIndexerObjectProperty) _flowInterfaceBodyStatement() {}
func (*FlowInterfaceMethod) _flowInterfaceBodyStatement()       {}

func (*FlowTypeAssertion) _expressionNode() {}

func (*FlowTypeStatement) _statementNode()      {}
func (*FlowInterfaceStatement) _statementNode() {}

func (f *FlowTypeAssertion) StartAt() file.Idx      { return f.Left.StartAt() }
func (f *FlowTypeStatement) StartAt() file.Idx      { return f.Start }
func (f *FlowInterfaceStatement) StartAt() file.Idx { return f.Start }

func (f *FlowTypeStatement) EndAt() file.Idx      { return f.Type.EndAt() }
func (f *FlowInterfaceStatement) EndAt() file.Idx { return f.End }

func (f *FlowPrimitiveType) EndAt() file.Idx     { return f.End }
func (f *FlowTrueType) EndAt() file.Idx          { return f.End }
func (f *FlowFalseType) EndAt() file.Idx         { return f.End }
func (f *FlowStringLiteralType) EndAt() file.Idx { return f.End }
func (f *FlowNumberLiteralType) EndAt() file.Idx { return f.End }
func (f *FlowIdentifier) EndAt() file.Idx        { return f.Start + file.Idx(len(f.Name)) }
func (f *FlowTypeOfType) EndAt() file.Idx {
	return f.Identifier.Start + file.Idx(len(f.Identifier.Name))
}
func (f *FlowTypeAssertion) EndAt() file.Idx   { return f.FlowType.EndAt() }
func (f *FlowOptionalType) EndAt() file.Idx    { return f.FlowType.EndAt() }
func (f *FlowInexactObject) EndAt() file.Idx   { return f.End }
func (f *FlowExactObject) EndAt() file.Idx     { return f.End }
func (f *FlowTupleType) EndAt() file.Idx       { return f.End }
func (f *FlowExistentialType) EndAt() file.Idx { return f.Start + 1 }
func (f *FlowFunctionType) EndAt() file.Idx    { return f.ReturnType.EndAt() }
