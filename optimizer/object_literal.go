package optimizer

import "yawp/parser/ast"

func (o *Optimizer) ObjectPropertyValue(opv *ast.ObjectPropertyValue) *ast.ObjectPropertyValue {
	if opv == nil {
		return nil
	}

	opv.Value = o.Expression(opv.Value)
	opv.PropertyName = o.ObjectPropertyName(opv.PropertyName)

	return opv
}
