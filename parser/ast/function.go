package ast

import "yawp/parser/file"

type (
	FunctionLiteral struct {
		Node
		Async          bool
		Generator      bool
		Id             *Identifier
		TypeParameters []*FlowTypeParameter
		ReturnType     FlowType
		Parameters     *FunctionParameters
		Body           *FunctionBody
	}

	FunctionParameters struct {
		Node
		List []FunctionParameter
	}

	FunctionBody struct {
		Node
		List Statements
	}

	FunctionParameter interface {
		GetDefaultValue() IExpr
		SetDefaultValue(IExpr)
		SetTypeAnnotation(FlowType)
		SetFlowTypeOptional(bool)
		_parameterNode()
	}

	IdentifierParameter struct {
		Id           *Identifier
		DefaultValue IExpr

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
		DefaultValue IExpr

		FlowType         FlowType
		FlowTypeOptional bool
	}
)

func (*FunctionLiteral) _expressionNode()    {}
func (*FunctionLiteral) _statementNode()     {}
func (f *FunctionLiteral) GetLoc() *file.Loc { return f.Loc }
func (f *FunctionLiteral) GetNode() *Node    { return &f.Node }

func (ip *IdentifierParameter) GetDefaultValue() IExpr { return ip.DefaultValue }
func (rp *RestParameter) GetDefaultValue() IExpr       { return nil }
func (odp *PatternParameter) GetDefaultValue() IExpr   { return odp.DefaultValue }

func (ip *IdentifierParameter) SetDefaultValue(exp IExpr) { ip.DefaultValue = exp }
func (rp *RestParameter) SetDefaultValue(_ IExpr)         {}
func (odp *PatternParameter) SetDefaultValue(exp IExpr)   { odp.DefaultValue = exp }

func (ip *IdentifierParameter) SetTypeAnnotation(flowType FlowType) { ip.FlowType = flowType }
func (rp *RestParameter) SetTypeAnnotation(flowType FlowType)       { rp.FlowType = flowType }
func (rp *PatternParameter) SetTypeAnnotation(flowType FlowType)    { rp.FlowType = flowType }

func (ip *IdentifierParameter) SetFlowTypeOptional(opt bool) { ip.FlowTypeOptional = opt }
func (rp *RestParameter) SetFlowTypeOptional(opt bool)       { rp.FlowTypeOptional = opt }
func (rp *PatternParameter) SetFlowTypeOptional(opt bool)    { rp.FlowTypeOptional = opt }

func (*IdentifierParameter) _parameterNode() {}
func (*RestParameter) _parameterNode()       {}
func (*PatternParameter) _parameterNode()    {}
