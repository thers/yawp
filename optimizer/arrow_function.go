package optimizer

import (
	"yawp/options"
	"yawp/parser/ast"
)

func (o *Optimizer) ArrowFunctionExpression(af *ast.ArrowFunctionExpression) *ast.ArrowFunctionExpression {
	if o.options.Target >= options.ES2015 {
		return o.Walker.ArrowFunctionExpression(af)
	}

	if af.Async {
		return o.asyncArrowFunction(af)
	}

	functionLiteral := &ast.FunctionLiteral{
		Loc: af.GetLoc(),
		Parameters: o.FunctionParameters(&ast.FunctionParameters{
			List: af.Parameters,
		}),
	}

	o.thisScope.NeedsReplacement = true
	defer func() {
		o.thisScope.NeedsReplacement = false
	}()

	functionLiteral.Body = o.Statement(af.Body)

	o.Walker.ReplacementExpression = functionLiteral

	return nil
}

func (o *Optimizer) asyncArrowFunction(af *ast.ArrowFunctionExpression) *ast.ArrowFunctionExpression {
	return o.Walker.ArrowFunctionExpression(af)
}
