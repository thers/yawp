package optimizer

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

func (o *Optimizer) pushThisScope() {
	o.thisScope = &ThisScope{
		Parent: o.thisScope,
	}
}

func (o *Optimizer) popThisScope() {
	if o.thisScope.Parent == nil {
		return
	}

	o.thisScope = o.thisScope.Parent
}

func (o *Optimizer) getThisReplacement() *ast.Identifier {
	if o.thisScope.ThisId == nil {
		o.thisScope.ThisId = o.refScope.GhostId()
	}

	if o.thisScope.ThisInitializer == nil {
		o.thisScope.ThisInitializer = undefinedThis
	}

	return o.thisScope.ThisId
}

func (o *Optimizer) getThisDeclaration() ast.Statement {
	if o.thisScope.ThisInitializer == nil {
		return nil
	}

	return &ast.VariableStatement{
		Kind: token.VAR,
		List: []*ast.VariableBinding{
			{
				Kind: token.VAR,
				Binder: &ast.IdentifierBinder{
					Id: o.thisScope.ThisId,
				},
				Initializer: o.thisScope.ThisInitializer,
			},
		},
	}
}
