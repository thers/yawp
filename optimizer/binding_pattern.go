package optimizer

import (
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/token"
)

const ghostRefIdName = ""

func (o *Optimizer) pushExtraVariableBinding(vb *ast.VariableBinding) {
	o.extraVariableBindings = append(o.extraVariableBindings, vb)
}

func (o *Optimizer) VariableStatement(vs *ast.VariableStatement) *ast.VariableStatement {
	// resolve variable kind to ref kind early
	o.bindingRefKind = o.resolveTokenToRefKind(vs.Kind)
	defer func(){
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
		ghostRef := o.refScope.GhostRef()
		ghostId := &ast.Identifier{
			Ref:  ghostRef,
			Name: ghostRefIdName, // so no allocation for this happens, meh
		}

		o.pushExtraVariableBinding(&ast.VariableBinding{
			Kind: token.VAR,
			Binder: &ast.IdentifierBinder{
				Id: ghostId,
			},
			Initializer: vb.Initializer,
		})
		vb.Initializer = ghostId
	}

	// we have to proceed binder here as well as we need information about initializer
	// now, if it's an identifier binder it's an easy way
	// so we can leave it as it is and call it a day
	switch binder := vb.Binder.(type) {
	case *ast.IdentifierBinder:
		vb.Binder = o.IdentifierBinder(binder)
	case *ast.ObjectBinding:
		return o.es5ObjectBinding(binder, vb)
	default:
		vb.Binder = o.PatternBinder(binder)
	}

	return vb
}

func (o *Optimizer) es5ObjectBinding(ob *ast.ObjectBinding, vb *ast.VariableBinding) *ast.VariableBinding {
	// we want to keep the last one for the binding place as
	for i := 0; i <= len(ob.List)-2; i++ {

	}

	pb := ob.List[0].(*ast.ObjectPropertyBinder)

	newInitializer := &ast.DotExpression{
		Left:       vb.Initializer,
		Identifier: pb.PropertyName.Clone(),
	}

	return &ast.VariableBinding{
		Loc:         vb.Loc,
		Kind:        token.VAR,
		Binder:      o.IdentifierBinder(&ast.IdentifierBinder{
			Id: pb.PropertyName,
		}),
		Initializer: newInitializer,
	}
}
