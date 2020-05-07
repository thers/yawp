package optimizer

import (
	"strconv"
	"yawp/builtins"
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (o *Optimizer) pushExtraVariableBinding(vb *ast.VariableBinding) {
	o.extraVariableBindings = append(o.extraVariableBindings, vb)
}

func (o *Optimizer) VariableStatement(vs *ast.VariableStatement) ast.Statement {
	// resolve variable kind to ref kind early
	o.bindingRefKind = o.resolveTokenToRefKind(vs.Kind)
	defer func() {
		o.bindingRefKind = ast.RUnknown
	}()

	if o.options.Target < options.ES2015 {
		// in ES5 only var it is
		vs.Kind = token.VAR

		o.extraVariableBindings = make([]*ast.VariableBinding, 0)

		for index, vb := range vs.List {
			vs.List[index] = o.VariableBinding(vb)
			vs.Kind = token.VAR
		}

		vs.List = append(o.extraVariableBindings, vs.List...)
	}

	return vs
}

func (o *Optimizer) VariableBinding(vb *ast.VariableBinding) *ast.VariableBinding {
	// processing initializer first so we resolve id refs correctly
	// for example, this code:
	// `a=1; { const a=a; log(a) }`
	// should be treated like this:
	// `a=1; { const b=a; lob(b) }`

	// for ES2015+ we can keep both var kinds and destructuring as it is, yay
	// just have to deal with refs and it is
	if o.options.Target >= options.ES2015 {
		o.Expression(vb.Initializer)
		o.PatternBinder(vb.Binder)

		return vb
	}

	vb.Initializer = o.Expression(vb.Initializer)

	// we have to transform destructuring patterns into ES5 stuff
	// so if we have destructuring patter AND initializer is something
	// that is not identifier  introduce new variable with that initializer

	_, initializerIsIdentifier := vb.Initializer.(*ast.Identifier)
	_, binderIsIdentifier := vb.Binder.(*ast.IdentifierBinder)

	if !initializerIsIdentifier && !binderIsIdentifier {
		o.forkVariableBinding(vb)
	}

	return o.es5PatternBinder(vb.Binder, vb)
}

func (o *Optimizer) forkVariableBinding(vb *ast.VariableBinding) {
	ghostId := o.refScope.GhostId()

	o.pushExtraVariableBinding(&ast.VariableBinding{
		Loc:  vb.GetLoc(),
		Kind: token.VAR,
		Binder: &ast.IdentifierBinder{
			Id: ghostId,
		},
		Initializer: vb.Initializer,
	})

	vb.Initializer = ghostId
}

func (o *Optimizer) es5PatternBinder(pb ast.PatternBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	switch b := pb.(type) {
	case *ast.IdentifierBinder:
		vb.Binder = o.IdentifierBinder(b)
	case *ast.ObjectBinding:
		return o.es5ObjectBinding(b, vb)
	case *ast.ObjectPropertyBinder:
		return o.es5ObjectPropertyBinder(b, vb)
	case *ast.ArrayBinding:
		return o.es5ArrayBinding(b, vb)
	case *ast.ArrayItemBinder:
		return o.es5ArrayItemBinder(b, vb)
	case *ast.ArrayRestBinder:
		return o.es5ArrayRestBinder(b, vb)
	case *ast.ObjectRestBinder:
		return o.es5ObjectRestBinder(b, vb)
	}

	return vb
}

func (o *Optimizer) es5ObjectBinding(ob *ast.ObjectBinding, vb *ast.VariableBinding) *ast.VariableBinding {
	for index, propBinder := range ob.List {
		nvb := o.es5PatternBinder(propBinder, vb.Clone())

		// last one should be returned
		if index == len(ob.List)-1 {
			nvb.Loc = vb.Loc
			return nvb
		} else {
			o.pushExtraVariableBinding(nvb)
		}
	}

	panic("No last property?")
}

func (o *Optimizer) es5ArrayBinding(ab *ast.ArrayBinding, vb *ast.VariableBinding) *ast.VariableBinding {
	for index, propBinder := range ab.List {
		nvb := o.es5PatternBinder(propBinder, vb.Clone())

		// last one should be returned
		if index == len(ab.List)-1 {
			nvb.Loc = vb.Loc
			return nvb
		} else {
			o.pushExtraVariableBinding(nvb)
		}
	}

	panic("No last property?")
}

func (o *Optimizer) es5ObjectPropertyBinder(opb *ast.ObjectPropertyBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	vb.Initializer = &ast.MemberExpression{
		Left:  vb.Initializer,
		Right: o.Expression(opb.Id),
		Kind:  ast.MKObject,
	}

	if opb.DefaultValue != nil {
		vb.Initializer = &ast.BinaryExpression{
			Operator: token.LOGICAL_OR,
			Left:     vb.Initializer,
			Right:    o.Expression(opb.DefaultValue),
		}

		// we shouldn't init initializer twice, so we have to fork now
		o.forkVariableBinding(vb)
	}

	return o.es5PatternBinder(opb.Binder, vb)
}

func (o *Optimizer) es5ArrayItemBinder(aib *ast.ArrayItemBinder, vb *ast.VariableBinding) *ast.VariableBinding {
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
			Right:    o.Expression(aib.DefaultValue),
		}

		// we shouldn't init initializer twice, so we have to fork now
		o.forkVariableBinding(vb)
	}

	return o.es5PatternBinder(aib.Binder, vb)
}

func (o *Optimizer) es5ArrayRestBinder(arb *ast.ArrayRestBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	vb.Initializer = &ast.CallExpression{
		Callee: builtins.SlicedArrayRest,
		ArgumentList: []ast.Expression{
			vb.Initializer,
			&ast.NumberLiteral{
				Literal: strconv.Itoa(arb.FromIndex),
			},
		},
	}
	vb.Binder = o.PatternBinder(arb.Binder)

	return vb
}

func (o *Optimizer) es5ObjectRestBinder(orb *ast.ObjectRestBinder, vb *ast.VariableBinding) *ast.VariableBinding {
	o.module.Additions.ObjectOmit = true

	vb.Initializer = &ast.CallExpression{
		Callee: builtins.ObjectRest,
		ArgumentList: []ast.Expression{
			vb.Initializer,
			&ast.ArrayLiteral{
				List: propertiesToStrings(orb.OmitProperties),
			},
		},
	}
	vb.Binder = o.PatternBinder(orb.Binder)

	return vb
}

func propertiesToStrings(list []ast.ObjectPropertyName) []ast.Expression {
	strs := make([]ast.Expression, 0, len(list))

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
