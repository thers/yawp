package transpiler

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

var undefinedThis = &ast.Js{Code: "void 0"}

type ThisScope struct {
	Parent *ThisScope

	NeedsReplacement bool

	ThisId          *ast.Identifier
	ThisInitializer ast.Expression
}

func (t *Transpiler) pushThisScope() {
	t.thisScope = &ThisScope{
		Parent: t.thisScope,
	}
}

func (t *Transpiler) popThisScope() {
	if t.thisScope.Parent == nil {
		return
	}

	t.thisScope = t.thisScope.Parent
}

func (t *Transpiler) getThisReplacement() *ast.Identifier {
	if t.thisScope.ThisId == nil {
		t.thisScope.ThisId = t.refScope.GhostId()
	}

	if t.thisScope.ThisInitializer == nil {
		t.thisScope.ThisInitializer = undefinedThis
	}

	return t.thisScope.ThisId
}

func (t *Transpiler) getThisDeclaration() ast.Statement {
	if t.thisScope.ThisInitializer == nil {
		return nil
	}

	return &ast.VariableStatement{
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
	}
}
