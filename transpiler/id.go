package transpiler

import "yawp/parser/ast"

func (t *Transpiler) IdentifierBinder(vb *ast.IdentifierBinder) *ast.IdentifierBinder {
	if t.bindingRefKind != ast.SRUnknown {
		// Binding yet unknown id
		vb.Id.LegacyRef = t.refScope.BindRef(t.bindingRefKind, vb.Id.Name)
	}

	return vb
}

func (t *Transpiler) Identifier(id *ast.Identifier) *ast.Identifier {
	id.LegacyRef = t.refScope.UseRef(id.Name)

	return id
}
