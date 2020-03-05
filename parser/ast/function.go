package ast

import "yawp/parser/file"

type (
	FunctionLiteral struct {
		Start      file.Idx
		Async      bool
		Name       *Identifier
		Parameters *FunctionParameters
		Body       Statement
		Source     string

		DeclarationList []Declaration
	}

	FunctionParameter interface {
		GetDefaultValue() Expression
		SetDefaultValue(Expression)
		_parameterNode()
	}

	IdentifierParameter struct {
		Name         *Identifier
		DefaultValue Expression
	}

	RestParameter struct {
		Binder PatternBinder
	}

	ObjectPatternIdentifierParameter struct {
		Parameter    FunctionParameter
		PropertyName string
	}

	ObjectPatternParameter struct {
		List         []*ObjectPatternIdentifierParameter
		DefaultValue Expression
	}

	ArrayPatternParameter struct {
		List         []FunctionParameter
		DefaultValue Expression
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
func (rp *RestParameter) SetDefaultValue(_ Expression)           {}
func (odp *ObjectPatternParameter) SetDefaultValue(exp Expression) { odp.DefaultValue = exp }
func (adp *ArrayPatternParameter) SetDefaultValue(exp Expression)  { adp.DefaultValue = exp }

func (*IdentifierParameter) _parameterNode()    {}
func (*RestParameter) _parameterNode()          {}
func (*ObjectPatternParameter) _parameterNode() {}
func (*ArrayPatternParameter) _parameterNode()  {}
