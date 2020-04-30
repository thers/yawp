package optimizer

import (
	"yawp/parser/ast"
)

func (o *Optimizer) ThisExpression(te *ast.ThisExpression) *ast.ThisExpression {
	if !o.thisScope.NeedsReplacement {
		return te
	}

	o.Walker.ReplacementExpression = o.getThisReplacement()

	return nil
}
