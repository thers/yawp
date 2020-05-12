package transpiler

import "yawp/parser/ast"

func (t *Transpiler) ObjectPropertyValue(opv *ast.ObjectPropertyValue) *ast.ObjectPropertyValue {
	if opv == nil {
		return nil
	}

	opv.Value = t.Expression(opv.Value)
	opv.PropertyName = t.ObjectPropertyName(opv.PropertyName)

	return opv
}
