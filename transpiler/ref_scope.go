package transpiler

import (
	"yawp/ids"
	"yawp/parser/ast"
	"yawp/parser/token"
)

const ghostIdName = ""

func (t *Transpiler) pushRefScope() *RefScope {
	parentRefScope := t.refScope

	t.refScope = &RefScope{
		Parent: parentRefScope,
		Refs:   make(map[string]*ast.SymbolRef, 0),
		ids:    t.ids,
		minify: t.options.Minify,
	}

	return t.refScope
}

func (t *Transpiler) popRefScope() *RefScope {
	refScope := t.refScope

	if t.refScope.Parent != nil {
		t.refScope = t.refScope.Parent
	}

	return refScope
}

func (t *Transpiler) resolveTokenToRefKind(tkn token.Token) (kind ast.SymbolRefType) {
	switch tkn {
	case token.VAR:
		kind = ast.SRVar
	case token.CONST:
		kind = ast.SRConst
	case token.LET:
		kind = ast.SRLet
	default:
		kind = ast.SRUnknown
	}

	return
}

type RefScope struct {
	Parent *RefScope
	Refs   map[string]*ast.SymbolRef

	ids    *ids.Ids
	minify bool
}

func (r *RefScope) NextMangledId() string {
	return r.ids.Next()
}

func (r *RefScope) createRef(name string) *ast.SymbolRef {
	var shadowedRef *ast.SymbolRef

	ref := &ast.SymbolRef{
		Name:   name,
		Usages: 0,
	}

	if r.Parent != nil {
		shadowedRef = r.Parent.GetRef(name)

		if shadowedRef != nil {
			ref.ShadowsRef = shadowedRef
			shadowedRef.ShadowedByRef = ref
		}
	}

	r.Refs[name] = ref

	return ref
}

func (r *RefScope) GhostRef() *ast.SymbolRef {
	return &ast.SymbolRef{
		Name: r.NextMangledId(),
		Type: ast.SRVar,
	}
}

func (r *RefScope) GhostId() *ast.Identifier {
	return &ast.Identifier{
		LegacyRef: r.GhostRef(),
		Name:      ghostIdName,
	}
}

func (r *RefScope) GetRef(name string) *ast.SymbolRef {
	var ref *ast.SymbolRef
	var ok bool

	if ref, ok = r.Refs[name]; !ok {
		if r.Parent != nil {
			return r.Parent.GetRef(name)
		}
	}

	return ref
}

func (r *RefScope) BindRef(kind ast.SymbolRefType, name string) *ast.SymbolRef {
	var ref *ast.SymbolRef
	var ok bool

	// vars can hoist declarations and they're not block-scoped
	// we also don't even bother mangling them
	if kind == ast.SRVar {
		ref = r.GetRef(name)

		if ref != nil {
			if ref.Type == ast.SRUnknown {
				ref.Type = ast.SRVar
			}

			return ref
		}
	}

	if ref, ok = r.Refs[name]; !ok {
		ref = r.createRef(name)
	}

	if r.minify {
		ref.Name = r.NextMangledId()
	}

	ref.Type = kind

	return ref
}

func (r *RefScope) UseRef(name string) *ast.SymbolRef {
	var ref *ast.SymbolRef

	ref = r.GetRef(name)

	if ref == nil {
		ref = r.createRef(name)
	}

	ref.Usages++

	return ref
}
