package ast

import (
	"yawp/parser/file"
	"yawp/parser/token"
)

type (
	FlowType interface {
		INode
		_flowType()
	}

	FlowInterfaceBodyStatement interface {
		_flowInterfaceBodyStatement()
	}

	FlowObjectProperty interface {
		_flowObjectProperty()
	}

	FlowTypeParameter struct {
		Loc           *file.Loc
		Name          *FlowIdentifier
		Covariant     bool
		Contravariant bool
		Boundary      FlowType
		DefaultValue  FlowType
	}

	FlowInterfaceMethod struct {
		Loc            *file.Loc
		Name           *FlowIdentifier
		TypeParameters []*FlowTypeParameter
		Parameters     []interface{}
		ReturnType     FlowType
	}

	FlowIdentifier struct {
		Loc           *file.Loc
		Name          string
		Qualification *FlowIdentifier
	}

	FlowObjectIndexer struct {
		Key   FlowType
		Value FlowType
	}

	FlowInexactObject struct {
		Loc        *file.Loc
		Properties []FlowObjectProperty
	}

	FlowExactObject struct {
		Loc        *file.Loc
		Properties []FlowObjectProperty
	}

	FlowNamedObjectProperty struct {
		Loc           *file.Loc
		Optional      bool
		Covariant     bool
		Contravariant bool
		Name          string
		Value         FlowType
	}

	FlowIndexerObjectProperty struct {
		Loc     *file.Loc
		KeyName string
		KeyType FlowType
		Value   FlowType
	}

	FlowInexactSpecifierProperty struct {
		Loc *file.Loc
	}

	FlowSpreadObjectProperty struct {
		Loc      *file.Loc
		FlowType FlowType
	}

	FlowUnionType struct {
		Loc   *file.Loc
		Types []FlowType
	}

	FlowIntersectionType struct {
		Loc   *file.Loc
		Types []FlowType
	}

	FlowTypeOfType struct {
		Loc        *file.Loc
		Identifier *FlowIdentifier
	}

	FlowOptionalType struct {
		FlowType FlowType
	}

	FlowFunctionType struct {
		Loc            *file.Loc
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
		Loc  *file.Loc
		Kind token.Token
	}

	FlowTrueType struct {
		Loc *file.Loc
	}
	FlowFalseType struct {
		Loc *file.Loc
	}

	FlowExistentialType struct {
		Loc *file.Loc
	}

	FlowStringLiteralType struct {
		Loc    *file.Loc
		String string
	}

	FlowNumberLiteralType struct {
		Loc    *file.Loc
		Number interface{}
	}

	FlowTupleType struct {
		Loc      *file.Loc
		Elements []FlowType
	}

	FlowArrayType struct {
		Loc         *file.Loc
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

func (*FlowTypeStatement) _exportClauseNode()      {}
func (*FlowInterfaceStatement) _exportClauseNode() {}

func (f *FlowPrimitiveType) GetLoc() *file.Loc     { return f.Loc }
func (f *FlowTrueType) GetLoc() *file.Loc          { return f.Loc }
func (f *FlowFalseType) GetLoc() *file.Loc         { return f.Loc }
func (f *FlowStringLiteralType) GetLoc() *file.Loc { return f.Loc }
func (f *FlowNumberLiteralType) GetLoc() *file.Loc { return f.Loc }
func (f *FlowIdentifier) GetLoc() *file.Loc        { return f.Loc }
func (f *FlowTypeOfType) GetLoc() *file.Loc        { return f.Loc.Add(f.Identifier.GetLoc()) }
func (f *FlowOptionalType) GetLoc() *file.Loc      { return f.FlowType.GetLoc() }
func (f *FlowInexactObject) GetLoc() *file.Loc     { return f.Loc }
func (f *FlowExactObject) GetLoc() *file.Loc       { return f.Loc }
func (f *FlowTupleType) GetLoc() *file.Loc         { return f.Loc }
func (f *FlowArrayType) GetLoc() *file.Loc         { return f.Loc }
func (f *FlowExistentialType) GetLoc() *file.Loc   { return f.Loc }
func (f *FlowFunctionType) GetLoc() *file.Loc      { return f.Loc }

func (f *FlowUnionType) GetLoc() *file.Loc {
	return f.Loc.Add(f.Types[len(f.Types)-1].GetLoc())
}
func (f *FlowIntersectionType) GetLoc() *file.Loc {
	return f.Loc.Add(f.Types[len(f.Types)-1].GetLoc())
}

func (f *FlowGenericType) GetLoc() *file.Loc {
	return f.TypeArguments[len(f.TypeArguments)-1].GetLoc()
}

func (fi *FlowIdentifier) Identifier() *Identifier {
	return &Identifier{
		ExprNode: ExprNode{
			Loc: fi.Loc,
		},
		Name: fi.Name,
	}
}

func (f *FlowPrimitiveType) GetNode() *Node     { return nil }
func (f *FlowTrueType) GetNode() *Node          { return nil }
func (f *FlowFalseType) GetNode() *Node         { return nil }
func (f *FlowStringLiteralType) GetNode() *Node { return nil }
func (f *FlowNumberLiteralType) GetNode() *Node { return nil }
func (f *FlowIdentifier) GetNode() *Node        { return nil }
func (f *FlowTypeOfType) GetNode() *Node        { return nil }
func (f *FlowOptionalType) GetNode() *Node      { return nil }
func (f *FlowInexactObject) GetNode() *Node     { return nil }
func (f *FlowExactObject) GetNode() *Node       { return nil }
func (f *FlowTupleType) GetNode() *Node         { return nil }
func (f *FlowExistentialType) GetNode() *Node   { return nil }
func (f *FlowFunctionType) GetNode() *Node      { return nil }
func (f *FlowUnionType) GetNode() *Node         { return nil }
func (f *FlowIntersectionType) GetNode() *Node  { return nil }
func (f *FlowGenericType) GetNode() *Node       { return nil }
func (f *FlowArrayType) GetNode() *Node         { return nil }
