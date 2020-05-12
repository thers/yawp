package transpiler

import (
	"yawp/parser/ast"
)

func (t *Transpiler) ThisExpression(te *ast.ThisExpression) *ast.ThisExpression {
	if !t.thisScope.NeedsReplacement {
		return te
	}

	t.Walker.ReplacementExpression = t.getThisReplacement()

	return nil
}
