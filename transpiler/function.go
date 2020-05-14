package transpiler

import (
	"yawp/builtins"
	"yawp/parser/ast"
	"yawp/parser/token"
)

type FunctionScope struct {
	ExtraVariables []*ast.VariableBinding
	ParameterIndex int
}

func (t *Transpiler) pushExtraVariableToFunctionScope(vb *ast.VariableBinding) {
	if t.functionScope == nil {
		panic("Attempt to push extra variable to nil function scope")
	}

	t.functionScope.ExtraVariables = append(t.functionScope.ExtraVariables, vb)
}

func (t *Transpiler) unshiftExtraVariableToFunctionScope(vb *ast.VariableBinding) {
	if t.functionScope == nil {
		panic("Attempt to push extra variable to nil function scope")
	}

	t.functionScope.ExtraVariables = append([]*ast.VariableBinding{vb}, t.functionScope.ExtraVariables...)
}

func (t *Transpiler) FunctionLiteral(fl *ast.FunctionLiteral) *ast.FunctionLiteral {
	if fl.Id != nil {
		fl.Id.Ref = t.refScope.BindRef(ast.RFn, fl.Id.Name)
	}

	// Ref scope starts from arguments
	t.pushRefScope()
	defer t.popRefScope()

	popFunctionScope := t.pushFunctionScope()
	defer popFunctionScope()

	fl.Parameters = t.FunctionParameters(fl.Parameters)

	body := fl.Body

	if len(t.functionScope.ExtraVariables) > 0 {
		body.List = append(ast.Statements{
			&ast.VariableStatement{
				Kind: token.VAR,
				List: t.functionScope.ExtraVariables,
			},
		}, body.List)
	} else {
		body = fl.Body
	}

	fl.Body = t.FunctionBody(body)

	return fl
}

func (t *Transpiler) FunctionParameters(fp *ast.FunctionParameters) *ast.FunctionParameters {
	if t.functionScope == nil {
		panic("Can not transpile function parameters while not in function scope")
	}

	list := make([]ast.FunctionParameter, 0)

	for index, parameter := range fp.List {
		t.functionScope.ParameterIndex = index
		parameter = t.FunctionParameter(parameter)

		if parameter != nil {
			list = append(list, parameter)
		}
	}

	return &ast.FunctionParameters{
		Node: fp.Node,
		List: list,
	}
}

func (t *Transpiler) IdentifierParameter(ip *ast.IdentifierParameter) ast.FunctionParameter {
	ip.Id.Ref = t.refScope.BindRef(ast.RFnParam, ip.Id.Name)

	if ip.DefaultValue != nil {
		// Unshift as it must come before other bindings
		t.unshiftExtraVariableToFunctionScope(&ast.VariableBinding{
			Kind: token.VAR,
			Binder: &ast.IdentifierBinder{
				Id: ip.Id,
			},
			Initializer: &ast.BinaryExpression{
				Operator:   token.LOGICAL_OR,
				Left:       ip.Id,
				Right:      t.Expression(ip.DefaultValue),
				Comparison: false,
			},
		})

		ip.DefaultValue = nil
	}

	return ip
}

func (t *Transpiler) RestParameter(rp *ast.RestParameter) ast.FunctionParameter {
	// Pushing variable declaration into function body
	t.pushExtraVariableToFunctionScope(&ast.VariableBinding{
		Kind:        token.VAR,
		Binder:      rp.Binder,
		Initializer: createArraySlice(builtins.Arguments, t.functionScope.ParameterIndex),
	})

	// Since we don't have this arg in fact now, remove it
	return nil
}

func (t *Transpiler) PatternParameter(pp *ast.PatternParameter) ast.FunctionParameter {
	// Pattern parameter become a simple id parameter
	// and we're pushing pattern binding to the function body
	paramId := t.refScope.GhostId()
	paramId.Ref.Kind = ast.RFnParam

	t.pushExtraVariableToFunctionScope(&ast.VariableBinding{
		Kind:        token.VAR,
		Binder:      pp.Binder,
		Initializer: paramId,
	})

	return t.IdentifierParameter(&ast.IdentifierParameter{
		Id:           paramId,
		DefaultValue: pp.DefaultValue,
	})
}
