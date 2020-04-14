package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type (
	FlowTypeStatement struct {
		Start          file.Loc
		Name           *FlowIdentifier
		Opaque         bool
		Type           FlowType
		TypeParameters []*FlowTypeParameter
	}

	FlowInterfaceStatement struct {
		Start          file.Loc
		End            file.Loc
		Name           *FlowIdentifier
		TypeParameters []*FlowTypeParameter
		Body           []FlowInterfaceBodyStatement
	}

	FlowType interface {
		_flowType()
		EndAt() file.Loc
	}

	FlowInterfaceBodyStatement interface {
		_flowInterfaceBodyStatement()
	}

	FlowObjectProperty interface {
		_flowObjectProperty()
	}

	FlowTypeParameter struct {
		Start         file.Loc
		Name          *FlowIdentifier
		Covariant     bool
		Contravariant bool
		Boundary      FlowType
		DefaultValue  FlowType
	}

	FlowInterfaceMethod struct {
		Start          file.Loc
		Name           *FlowIdentifier
		TypeParameters []*FlowTypeParameter
		Parameters     []interface{}
		ReturnType     FlowType
	}

	FlowIdentifier struct {
		Start         file.Loc
		Name          string
		Qualification *FlowIdentifier
	}

	FlowObjectIndexer struct {
		Key   FlowType
		Value FlowType
	}

	FlowInexactObject struct {
		Start      file.Loc
		End        file.Loc
		Properties []FlowObjectProperty
	}

	FlowExactObject struct {
		Start      file.Loc
		End        file.Loc
		Properties []FlowObjectProperty
	}

	FlowNamedObjectProperty struct {
		Start         file.Loc
		Optional      bool
		Covariant     bool
		Contravariant bool
		Name          string
		Value         FlowType
	}

	FlowIndexerObjectProperty struct {
		Start   file.Loc
		KeyName string
		KeyType FlowType
		Value   FlowType
	}

	FlowInexactSpecifierProperty struct {
		Start file.Loc
	}

	FlowSpreadObjectProperty struct {
		Start    file.Loc
		FlowType FlowType
	}

	FlowUnionType struct {
		Start file.Loc
		Types []FlowType
	}

	FlowIntersectionType struct {
		Start file.Loc
		Types []FlowType
	}

	FlowTypeOfType struct {
		Start      file.Loc
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
		Start          file.Loc
		Parameters     []*FlowFunctionParameter
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
	}

	FlowFunctionParameter struct {
		Identifier *Identifier
		Type       FlowType
	}

	FlowGenericType struct {
		Name          *FlowIdentifier
		TypeArguments []FlowType
	}
)

type (
	// Primitives
	FlowPrimitiveType struct {
		Start file.Loc
		End   file.Loc
		Kind  token.Token
	}

	FlowTrueType struct {
		Start file.Loc
		End   file.Loc
	}
	FlowFalseType struct {
		Start file.Loc
		End   file.Loc
	}

	FlowExistentialType struct {
		Start file.Loc
	}

	FlowStringLiteralType struct {
		Start  file.Loc
		End    file.Loc
		String string
	}

	FlowNumberLiteralType struct {
		Start  file.Loc
		End    file.Loc
		Number interface{}
	}

	FlowTupleType struct {
		Start    file.Loc
		End      file.Loc
		Elements []FlowType
	}

	FlowArrayType struct {
		End         file.Loc
		ElementType FlowType
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
func (*FlowUnionType) _flowType()         {}
func (*FlowIntersectionType) _flowType()  {}
func (*FlowGenericType) _flowType()       {}
func (*FlowArrayType) _flowType()         {}

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

func (*FlowTypeStatement) _exportClauseNode()      {}
func (*FlowInterfaceStatement) _exportClauseNode() {}

func (f *FlowTypeAssertion) StartAt() file.Loc      { return f.Left.StartAt() }
func (f *FlowTypeStatement) StartAt() file.Loc      { return f.Start }
func (f *FlowInterfaceStatement) StartAt() file.Loc { return f.Start }

func (f *FlowTypeStatement) EndAt() file.Loc      { return f.Type.EndAt() }
func (f *FlowInterfaceStatement) EndAt() file.Loc { return f.End }

func (f *FlowPrimitiveType) EndAt() file.Loc     { return f.End }
func (f *FlowTrueType) EndAt() file.Loc          { return f.End }
func (f *FlowFalseType) EndAt() file.Loc         { return f.End }
func (f *FlowStringLiteralType) EndAt() file.Loc { return f.End }
func (f *FlowNumberLiteralType) EndAt() file.Loc { return f.End }
func (f *FlowIdentifier) EndAt() file.Loc        { return f.Start + file.Loc(len(f.Name)) }
func (f *FlowTypeOfType) EndAt() file.Loc {
	return f.Identifier.Start + file.Loc(len(f.Identifier.Name))
}
func (f *FlowTypeAssertion) EndAt() file.Loc    { return f.FlowType.EndAt() }
func (f *FlowOptionalType) EndAt() file.Loc     { return f.FlowType.EndAt() }
func (f *FlowInexactObject) EndAt() file.Loc    { return f.End }
func (f *FlowExactObject) EndAt() file.Loc      { return f.End }
func (f *FlowTupleType) EndAt() file.Loc        { return f.End }
func (f *FlowArrayType) EndAt() file.Loc        { return f.End }
func (f *FlowExistentialType) EndAt() file.Loc  { return f.Start + 1 }
func (f *FlowFunctionType) EndAt() file.Loc     { return f.ReturnType.EndAt() }
func (f *FlowUnionType) EndAt() file.Loc        { return f.Types[len(f.Types)-1].EndAt() }
func (f *FlowIntersectionType) EndAt() file.Loc { return f.Types[len(f.Types)-1].EndAt() }
func (f *FlowGenericType) EndAt() file.Loc {
	return f.TypeArguments[len(f.TypeArguments)-1].EndAt() + 1
}

func (fi *FlowIdentifier) Identifier() *Identifier {
	return &Identifier{
		Start: fi.Start,
		Name:  fi.Name,
	}
}
