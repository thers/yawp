package ast

import "yawp/parser/file"

type (
	FunctionLiteral struct {
		Start          file.Idx
		Async          bool
		Generator      bool
		Name           *Identifier
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     *FunctionParameters
		Body           Statement
		Source         string

		DeclarationList []Declaration
	}

	FunctionParameter interface {
		GetDefaultValue() Expression
		SetDefaultValue(Expression)
		SetTypeAnnotation(FlowType)
		SetOptional(bool)
		_parameterNode()
	}

	IdentifierParameter struct {
		Name         *Identifier
		DefaultValue Expression
		FlowType     FlowType
		Optional     bool
	}

	RestParameter struct {
		Binder   PatternBinder
		FlowType FlowType
		Optional bool
	}

	ObjectPatternIdentifierParameter struct {
		Parameter    FunctionParameter
		PropertyName string
	}

	ObjectPatternParameter struct {
		List         []*ObjectPatternIdentifierParameter
		DefaultValue Expression
		FlowType     FlowType
		Optional     bool
	}

	ArrayPatternParameter struct {
		List         []FunctionParameter
		DefaultValue Expression
		FlowType     FlowType
		Optional     bool
	}
)

func (*FunctionLiteral) _expressionNode() {}
func (*FunctionLiteral) _statementNode()  {}

func (f *FunctionLiteral) StartAt() file.Idx { return f.Start }

func (f *FunctionLiteral) EndAt() file.Idx { return f.Body.EndAt() }

func (ip *IdentifierParameter) GetDefaultValue() Expression     { return ip.DefaultValue }
func (rp *RestParameter) GetDefaultValue() Expression           { return nil }
func (odp *ObjectPatternParameter) GetDefaultValue() Expression { return odp.DefaultValue }
func (adp *ArrayPatternParameter) GetDefaultValue() Expression  { return adp.DefaultValue }

func (ip *IdentifierParameter) SetDefaultValue(exp Expression)     { ip.DefaultValue = exp }
func (rp *RestParameter) SetDefaultValue(_ Expression)             {}
func (odp *ObjectPatternParameter) SetDefaultValue(exp Expression) { odp.DefaultValue = exp }
func (adp *ArrayPatternParameter) SetDefaultValue(exp Expression)  { adp.DefaultValue = exp }

func (ip *IdentifierParameter) SetTypeAnnotation(flowType FlowType)     { ip.FlowType = flowType }
func (rp *RestParameter) SetTypeAnnotation(flowType FlowType)           { rp.FlowType = flowType }
func (odp *ObjectPatternParameter) SetTypeAnnotation(flowType FlowType) { odp.FlowType = flowType }
func (adp *ArrayPatternParameter) SetTypeAnnotation(flowType FlowType)  { adp.FlowType = flowType }

func (ip *IdentifierParameter) SetOptional(opt bool)     { ip.Optional = opt }
func (rp *RestParameter) SetOptional(opt bool)           { rp.Optional = opt }
func (odp *ObjectPatternParameter) SetOptional(opt bool) { odp.Optional = opt }
func (adp *ArrayPatternParameter) SetOptional(opt bool)  { adp.Optional = opt }

func (*IdentifierParameter) _parameterNode()    {}
func (*RestParameter) _parameterNode()          {}
func (*ObjectPatternParameter) _parameterNode() {}
func (*ArrayPatternParameter) _parameterNode()  {}
