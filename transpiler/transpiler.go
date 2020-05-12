package transpiler

import (
	"yawp/ids"
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func Transpile(module *ast.Module, options *options.Options) {
	transpiler := &Transpiler{
		Walker:  ast.Walker{},
		module:  module,
		options: options,
		ids:     module.Ids,
	}
	transpiler.Walker.Visitor = transpiler
	transpiler.pushRefScope()
	transpiler.pushThisScope()

	module.Visit(transpiler)
}

type Transpiler struct {
	ast.Walker

	ids     *ids.Ids
	module  *ast.Module
	options *options.Options

	refScope  *RefScope
	thisScope *ThisScope

	bindingRefKind ast.RefKind

	extraVariables []*ast.VariableBinding

	functionScope *FunctionScope
}

func (t *Transpiler) pushFunctionScope() func() {
	functionScope := t.functionScope

	t.functionScope = &FunctionScope{
		ExtraVariables: make([]*ast.VariableBinding, 0),
		ParameterIndex: 0,
	}

	return func() {
		t.functionScope = functionScope
	}
}

func (t *Transpiler) Body(stmts []ast.Statement) []ast.Statement {
	stmts = t.Walker.Body(stmts)

	extras := make([]ast.Statement, 0)

	if t.thisScope.ThisInitializer != nil {
		extras = append(extras, &ast.VariableStatement{
			Kind: token.VAR,
			List: []*ast.VariableBinding{
				{
					Kind: token.VAR,
					Binder: &ast.IdentifierBinder{
						Id: t.thisScope.ThisId,
					},
					Initializer: t.thisScope.ThisInitializer,
				},
			},
		})
	}

	return append(extras, stmts...)
}

func (t *Transpiler) BlockStatement(bs *ast.BlockStatement) ast.Statement {
	t.pushRefScope()
	defer t.popRefScope()

	return t.Walker.BlockStatement(bs)
}

func (t *Transpiler) MemberExpression(me *ast.MemberExpression) ast.Expression {
	if leftIdentifier, ok := me.Left.(*ast.Identifier); ok {
		me.Left = t.Identifier(leftIdentifier)
	} else {
		me.Left = t.Expression(me.Left)
	}

	return me
}
