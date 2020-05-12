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
		Refs:   make(map[string]*ast.Ref, 0),
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

func (t *Transpiler) resolveTokenToRefKind(tkn token.Token) (kind ast.RefKind) {
	switch tkn {
	case token.VAR:
		kind = ast.RVar
	case token.CONST:
		kind = ast.RConst
	case token.LET:
		kind = ast.RLet
	default:
		kind = ast.RUnknown
	}

	return
}

type RefScope struct {
	Parent *RefScope
	Refs   map[string]*ast.Ref

	ids    *ids.Ids
	minify bool
}

func (r *RefScope) NextMangledId() string {
	return r.ids.Next()
}

func (r *RefScope) createRef(name string) *ast.Ref {
	var shadowedRef *ast.Ref

	ref := &ast.Ref{
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

func (r *RefScope) GhostRef() *ast.Ref {
	return &ast.Ref{
		Name:   r.NextMangledId(),
		Kind:   ast.RVar,
	}
}

func (r *RefScope) GhostId() *ast.Identifier {
	return &ast.Identifier{
		Ref:  r.GhostRef(),
		Name: ghostIdName,
	}
}

func (r *RefScope) GetRef(name string) *ast.Ref {
	var ref *ast.Ref
	var ok bool

	if ref, ok = r.Refs[name]; !ok {
		if r.Parent != nil {
			return r.Parent.GetRef(name)
		}
	}

	return ref
}

func (r *RefScope) BindRef(kind ast.RefKind, name string) *ast.Ref {
	var ref *ast.Ref
	var ok bool

	// vars can hoist declarations and they're not block-scoped
	// we also don't even bother mangling them
	if kind == ast.RVar {
		ref = r.GetRef(name)

		if ref != nil {
			if ref.Kind == ast.RUnknown {
				ref.Kind = ast.RVar
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

	ref.Kind = kind

	return ref
}

func (r *RefScope) UseRef(name string) *ast.Ref {
	var ref *ast.Ref

	ref = r.GetRef(name)

	if ref == nil {
		ref = r.createRef(name)
	}

	ref.Usages++

	return ref
}
