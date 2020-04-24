package optimizer

import (
	"yawp/ids"
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (o *Optimizer) pushRefScope() *RefScope {
	parentRefScope := o.refScope

	o.refScope = &RefScope{
		Parent: parentRefScope,
		Refs:   make(map[string]*ast.Ref, 0),
		ids:    o.ids,
		minify: o.options.Minify,
	}

	return o.refScope
}

func (o *Optimizer) popRefScope() *RefScope {
	refScope := o.refScope

	if o.refScope.Parent != nil {
		o.refScope = o.refScope.Parent
	}

	return refScope
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

func (r *RefScope) BindRef(tkn token.Token, name string) *ast.Ref {
	var kind ast.RefKind
	var ref *ast.Ref
	var ok bool

	switch tkn {
	case token.VAR:
		kind = ast.RVar
	case token.CONST:
		kind = ast.RConst
	case token.LET:
		kind = ast.RLet
	case token.IMPORT:
		kind = ast.RImport
	default:
		panic("Invalid ref kind")
	}

	// vars can hoist declarations and they're not block-scoped
	if kind == ast.RVar {
		ref = r.GetRef(name)

		if ref != nil {
			if ref.Kind == ast.RUnknown {
				ref.Kind = ast.RVar

				if r.minify {
					ref.Name = r.NextMangledId()
				}
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
