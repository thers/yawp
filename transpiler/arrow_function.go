package transpiler

import (
	"yawp/options"
	"yawp/parser/ast"
)

func (t *Transpiler) ArrowFunctionExpression(af *ast.ArrowFunctionExpression) *ast.ArrowFunctionExpression {
	if t.options.Target >= options.ES2015 {
		return t.Walker.ArrowFunctionExpression(af)
	}

	if af.Async {
		return t.asyncArrowFunction(af)
	}

	functionLiteral := &ast.FunctionLiteral{
		Loc: af.GetLoc(),
		Parameters: t.FunctionParameters(&ast.FunctionParameters{
			List: af.Parameters,
		}),
	}

	t.thisScope.NeedsReplacement = true
	defer func() {
		t.thisScope.NeedsReplacement = false
	}()

	functionLiteral.Body = t.FunctionBody(af.Body)

	t.Walker.ReplacementExpression = functionLiteral

	return nil
}

func (t *Transpiler) asyncArrowFunction(af *ast.ArrowFunctionExpression) *ast.ArrowFunctionExpression {
	return t.Walker.ArrowFunctionExpression(af)
}
