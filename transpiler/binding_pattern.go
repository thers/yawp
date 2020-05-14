package transpiler

import (
	"strconv"
	"yawp/builtins"
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (t *Transpiler) pushExtraVariableBinding(vb *ast.VariableBinding) {
	t.extraVariables = append(t.extraVariables, vb)
}

func (t *Transpiler) VariableStatement(vs *ast.VariableStatement) ast.IStmt {
	// resolve variable kind to ref kind early
	t.bindingRefKind = t.resolveTokenToRefKind(vs.Kind)
	defer func() {
		t.bindingRefKind = ast.RUnknown
	}()

	if t.options.Target < options.ES2015 {
		// in ES5 only var it is
		vs.Kind = token.VAR

		t.extraVariables = make([]*ast.VariableBinding, 0)

		for index, vb := range vs.List {
			vs.List[index] = t.VariableBinding(vb)
			vs.Kind = token.VAR
		}

		vs.List = append(t.extraVariables, vs.List...)
	}

	return vs
}

func (t *Transpiler) VariableBinding(vb *ast.VariableBinding) *ast.VariableBinding {
	// processing initializer first so we resolve id refs correctly
	// for example, this code:
	// `a=1; { const a=a; log(a) }`
	// should be treated like this:
	// `a=1; { const b=a; lob(b) }`

	// for ES2015+ we can keep both var kinds and destructuring as it is, yay
	// just have to deal with refs and it is
	if t.options.Target >= options.ES2015 {
		t.Expression(vb.Initializer)
		t.PatternBinder(vb.Binder)

		return vb
	}

	vb.Initializer = t.Expression(vb.Initializer)

	// we have to transform destructuring patterns into ES5 stuff
	// so if we have destructuring patter AND initializer is something
	// that is not identifier  introduce new variable with that initializer

	_, initializerIsIdentifier := vb.Initializer.(*ast.Identifier)
	_, binderIsIdentifier := vb.Binder.(*ast.IdentifierBinder)

	if !initializerIsIdentifier && !binderIsIdentifier {
		t.forkVariableBinding(vb)
	}

	return t.es5PatternBinder(vb.Binder, vb)
}

func (t *Transpiler) forkVariableBinding(vb *ast.VariableBinding) {
	ghostId := t.refScope.GhostId()

	t.pushExtraVariableBinding(&ast.VariableBinding{
		ExprNode: vb.ExprNode.Copy(),
		Kind: token.VAR,
		Binder: &ast.IdentifierBinder{
			Id: ghostId,
		},
		Initializer: vb.Initializer,
	})

	vb.Initializer = ghostId
}

func (t *Transpiler) es5PatternBinder(pb ast.PatternBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	switch b := pb.(type) {
	case *ast.IdentifierBinder:
		vb.Binder = t.IdentifierBinder(b)
	case *ast.ObjectBinding:
		return t.es5ObjectBinding(b, vb)
	case *ast.ObjectPropertyBinder:
		return t.es5ObjectPropertyBinder(b, vb)
	case *ast.ArrayBinding:
		return t.es5ArrayBinding(b, vb)
	case *ast.ArrayItemBinder:
		return t.es5ArrayItemBinder(b, vb)
	case *ast.ArrayRestBinder:
		return t.es5ArrayRestBinder(b, vb)
	case *ast.ObjectRestBinder:
		return t.es5ObjectRestBinder(b, vb)
	}

	return vb
}

func (t *Transpiler) es5ObjectBinding(ob *ast.ObjectBinding, vb *ast.VariableBinding) *ast.VariableBinding {
	for index, propBinder := range ob.List {
		nvb := t.es5PatternBinder(propBinder, vb.Copy())

		// last one should be returned
		if index == len(ob.List)-1 {
			nvb.Loc = vb.Loc
			return nvb
		} else {
			t.pushExtraVariableBinding(nvb)
		}
	}

	panic("No last property?")
}

func (t *Transpiler) es5ArrayBinding(ab *ast.ArrayBinding, vb *ast.VariableBinding) *ast.VariableBinding {
	for index, propBinder := range ab.List {
		nvb := t.es5PatternBinder(propBinder, vb.Copy())

		// last one should be returned
		if index == len(ab.List)-1 {
			nvb.Loc = vb.Loc
			return nvb
		} else {
			t.pushExtraVariableBinding(nvb)
		}
	}

	panic("No last property?")
}

func (t *Transpiler) es5ObjectPropertyBinder(opb *ast.ObjectPropertyBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	vb.Initializer = &ast.MemberExpression{
		Left:  vb.Initializer,
		Right: t.Expression(opb.PropertyName),
		Kind:  ast.MKObject,
	}

	if opb.DefaultValue != nil {
		vb.Initializer = &ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left:     vb.Initializer,
			Right:    t.Expression(opb.DefaultValue),
		}

		// we shouldn't init initializer twice, so we have to fork now
		t.forkVariableBinding(vb)
	}

	return t.es5PatternBinder(opb.Binder, vb)
}

func (t *Transpiler) es5ArrayItemBinder(aib *ast.ArrayItemBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	vb.Initializer = &ast.MemberExpression{
		Left: vb.Initializer,
		Right: &ast.NumberLiteral{
			Literal: strconv.Itoa(aib.Index),
		},
		Kind: ast.MKArray,
	}

	if aib.DefaultValue != nil {
		vb.Initializer = &ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left:     vb.Initializer,
			Right:    t.Expression(aib.DefaultValue),
		}

		// we shouldn't init initializer twice, so we have to fork now
		t.forkVariableBinding(vb)
	}

	return t.es5PatternBinder(aib.Binder, vb)
}

func (t *Transpiler) es5ArrayRestBinder(arb *ast.ArrayRestBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	vb.Initializer = createArraySlice(vb.Initializer, arb.FromIndex)
	vb.Binder = t.PatternBinder(arb.Binder)

	return vb
}

func (t *Transpiler) es5ObjectRestBinder(orb *ast.ObjectRestBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	t.module.Additions.ObjectOmit = true

	vb.Initializer = &ast.CallExpression{
		Callee: builtins.ObjectRest,
		ArgumentList: []ast.IExpr{
			vb.Initializer,
			&ast.ArrayLiteral{
				List: propertiesToStrings(orb.OmitProperties),
			},
		},
	}
	vb.Binder = t.PatternBinder(orb.Binder)

	return vb
}

func propertiesToStrings(list []ast.ObjectPropertyName) []ast.IExpr {
	strs := make([]ast.IExpr, 0, len(list))

	for _, id := range list {
		if identifier, ok := id.(*ast.Identifier); ok {
			strs = append(strs, &ast.StringLiteral{
				Literal: identifier.Name,
				Raw:     true,
			})
		}
	}

	return strs
}
