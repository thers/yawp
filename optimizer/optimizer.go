package optimizer

import (
	"yawp/ids"
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/token"
)

/**
Here we're performing:
	- ids resolving
	- ids mangling
*/

type Optimizer struct {
	ast.DefaultVisitor

	ids     *ids.Ids
	program *ast.Module
	options *options.Options

	refScope *RefScope

	variableBindingKind token.Token
}

func NewOptimizer(program *ast.Module, options *options.Options) *Optimizer {
	optimizer := &Optimizer{
		DefaultVisitor: ast.DefaultVisitor{},
		program:        program,
		options:        options,
		ids:            program.Ids,
	}
	optimizer.DefaultVisitor.Specific = optimizer
	optimizer.pushRefScope()

	return optimizer
}

func (o *Optimizer) VariableBinding(vb *ast.VariableBinding) *ast.VariableBinding {
	o.variableBindingKind = vb.Kind
	vb = o.DefaultVisitor.VariableBinding(vb)
	o.variableBindingKind = -1

	return vb
}

func (o *Optimizer) IdentifierBinder(vb *ast.IdentifierBinder) *ast.IdentifierBinder {
	if o.variableBindingKind > -1 {
		vb.Name.Ref = o.refScope.BindRef(o.variableBindingKind, vb.Name.Name)
	}

	return vb
}

func (o *Optimizer) BlockStatement(bs *ast.BlockStatement) *ast.BlockStatement {
	o.pushRefScope()
	defer o.popRefScope()

	return o.DefaultVisitor.BlockStatement(bs)
}

func (o *Optimizer) Identifier(id *ast.Identifier) *ast.Identifier {
	id.Ref = o.refScope.UseRef(id.Name)

	return id
}

func (o *Optimizer) DotExpression(de *ast.DotExpression) *ast.DotExpression {
	if leftIdentifier, ok := de.Left.(*ast.Identifier); ok {
		de.Left = o.Identifier(leftIdentifier)
	} else {
		de.Left = o.Expression(de.Left)
	}

	return de
}
