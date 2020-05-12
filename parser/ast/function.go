package ast

import "yawp/parser/file"

type (
	FunctionLiteral struct {
		Loc            *file.Loc
		Async          bool
		Generator      bool
		Id             *Identifier
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     *FunctionParameters
		Body           *FunctionBody

		DeclarationList []Declaration
	}

	FunctionParameters struct {
		Loc  *file.Loc
		List []FunctionParameter
	}

	FunctionBody struct {
		Loc  *file.Loc
		List Statements
	}

	FunctionParameter interface {
		GetDefaultValue() Expression
		SetDefaultValue(Expression)
		SetTypeAnnotation(FlowType)
		SetFlowTypeOptional(bool)
		_parameterNode()
	}

	IdentifierParameter struct {
		Id           *Identifier
		DefaultValue Expression

		FlowType         FlowType
		FlowTypeOptional bool
	}

	RestParameter struct {
		Binder PatternBinder

		FlowType         FlowType
		FlowTypeOptional bool
	}

	PatternParameter struct {
		Binder       PatternBinder
		DefaultValue Expression

		FlowType         FlowType
		FlowTypeOptional bool
	}
)

func (*FunctionLiteral) _expressionNode() {}
func (*FunctionLiteral) _statementNode()  {}

func (f *FunctionLiteral) GetLoc() *file.Loc { return f.Loc }

func (ip *IdentifierParameter) GetDefaultValue() Expression { return ip.DefaultValue }
func (rp *RestParameter) GetDefaultValue() Expression       { return nil }
func (odp *PatternParameter) GetDefaultValue() Expression   { return odp.DefaultValue }

func (ip *IdentifierParameter) SetDefaultValue(exp Expression) { ip.DefaultValue = exp }
func (rp *RestParameter) SetDefaultValue(_ Expression)         {}
func (odp *PatternParameter) SetDefaultValue(exp Expression)   { odp.DefaultValue = exp }

func (ip *IdentifierParameter) SetTypeAnnotation(flowType FlowType) { ip.FlowType = flowType }
func (rp *RestParameter) SetTypeAnnotation(flowType FlowType)       { rp.FlowType = flowType }
func (rp *PatternParameter) SetTypeAnnotation(flowType FlowType)    { rp.FlowType = flowType }

func (ip *IdentifierParameter) SetFlowTypeOptional(opt bool) { ip.FlowTypeOptional = opt }
func (rp *RestParameter) SetFlowTypeOptional(opt bool)       { rp.FlowTypeOptional = opt }
func (rp *PatternParameter) SetFlowTypeOptional(opt bool)    { rp.FlowTypeOptional = opt }

func (*IdentifierParameter) _parameterNode() {}
func (*RestParameter) _parameterNode()       {}
func (*PatternParameter) _parameterNode()    {}
