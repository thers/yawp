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
	- class lowering
*/

func Optimize(module *ast.Module, options *options.Options) {
	optimizer := &Optimizer{
		Walker:  ast.Walker{},
		module:  module,
		options: options,
		ids:     module.Ids,
	}
	optimizer.Walker.Visitor = optimizer
	optimizer.pushRefScope()
	optimizer.pushThisScope()

	module.Visit(optimizer)
}

type Optimizer struct {
	ast.Walker

	ids     *ids.Ids
	module  *ast.Module
	options *options.Options

	refScope  *RefScope
	thisScope *ThisScope

	bindingRefKind ast.RefKind

	extraVariableBindings []*ast.VariableBinding
}

func (o *Optimizer) Body(stmts []ast.Statement) []ast.Statement {
	stmts = o.Walker.Body(stmts)

	extras := make([]ast.Statement, 0)

	if o.thisScope.ThisInitializer != nil {
		extras = append(extras, &ast.VariableStatement{
			Kind: token.VAR,
			List: []*ast.VariableBinding{
				{
					Kind: token.VAR,
					Binder: &ast.IdentifierBinder{
						Id: o.thisScope.ThisId,
					},
					Initializer: o.thisScope.ThisInitializer,
				},
			},
		})
	}

	return append(extras, stmts...)
}

func (o *Optimizer) IdentifierBinder(vb *ast.IdentifierBinder) *ast.IdentifierBinder {
	if o.bindingRefKind != ast.RUnknown {
		vb.Id.Ref = o.refScope.BindRef(o.bindingRefKind, vb.Id.Name)
	}

	return vb
}

func (o *Optimizer) BlockStatement(bs *ast.BlockStatement) *ast.BlockStatement {
	o.pushRefScope()
	defer o.popRefScope()

	return o.Walker.BlockStatement(bs)
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

func (o *Optimizer) BracketExpression(be *ast.BracketExpression) *ast.BracketExpression {
	if leftIdentifier, ok := be.Left.(*ast.Identifier); ok {
		be.Left = o.Identifier(leftIdentifier)
	} else {
		be.Left = o.Expression(be.Left)
	}

	be.Member = o.Expression(be.Member)

	return be
}

func (o *Optimizer) FunctionLiteral(fl *ast.FunctionLiteral) *ast.FunctionLiteral {
	if fl.Id != nil {
		fl.Id.Ref = o.refScope.BindRef(ast.RFn, fl.Id.Name)
	}

	o.pushRefScope()
	defer o.popRefScope()

	fl.Parameters = o.FunctionParameters(fl.Parameters)
	fl.Body = o.Statement(fl.Body)

	return fl
}

func (o *Optimizer) IdentifierParameter(ip *ast.IdentifierParameter) *ast.IdentifierParameter {
	ip.DefaultValue = o.Expression(ip.DefaultValue)
	ip.Id.Ref = o.refScope.BindRef(ast.RFnParam, ip.Id.Name)

	return ip
}
