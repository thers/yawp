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
)

type (
	// Primitives
	FlowTrueType struct {
		Start file.Idx
	}
	FlowFalseType struct {
		Start file.Idx
	}
	FlowBooleanType struct {
		Start file.Idx
	}
	FlowVoidType struct {
		Start file.Idx
	}
	FlowNullType struct {
		Start file.Idx
	}

	FlowStringType struct {
		Start  file.Idx
		String string
	}

	FlowNumberType struct {
		Start  file.Idx
		Number interface{}
	}
)

func (*FlowTrueType) _flowType()    {}
func (*FlowFalseType) _flowType()   {}
func (*FlowBooleanType) _flowType() {}
func (*FlowStringType) _flowType()  {}
func (*FlowNumberType) _flowType()  {}
