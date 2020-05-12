package transpiler

import "yawp/parser/ast"

func (t *Transpiler) IdentifierBinder(vb *ast.IdentifierBinder) *ast.IdentifierBinder {
	if t.bindingRefKind != ast.RUnknown {
		// Binding yet unknown id
		vb.Id.Ref = t.refScope.BindRef(t.bindingRefKind, vb.Id.Name)
	}

	return vb
}

func (t *Transpiler) Identifier(id *ast.Identifier) *ast.Identifier {
	id.Ref = t.refScope.UseRef(id.Name)

	return id
}
