package generator

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (g *Generator) pushRefScope() *RefScope {
	parentRefScope := g.refScope

	g.refScope = &RefScope{
		Parent:      parentRefScope,
		Refs:        make(map[string]*ast.Ref, 0),
		idGenerator: newIdGenerator(),
	}

	return g.refScope
}

func (g *Generator) popRefScope() *RefScope {
	refScope := g.refScope

	if g.refScope.Parent != nil {
		g.refScope = g.refScope.Parent
	}

	return refScope
}

type RefScope struct {
	Parent *RefScope
	Refs   map[string]*ast.Ref

	idGenerator *IdGenerator
}

func (r *RefScope) NextMangledId() string {
	return r.idGenerator.Next()
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
	if ref, ok := r.Refs[name]; ok {
		return ref
	}

	return nil
}

func (r *RefScope) SetRef(tkn token.Token, name string) *ast.Ref {
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

	if ref, ok = r.Refs[name]; ok {
		ref.Kind = kind

		// previously we didn't know that this ref is declared
		ref.Usages++

		return ref
	}

	ref = r.createRef(name)
	ref.Kind = kind

	return ref
}

func (r *RefScope) UseRef(name string) *ast.Ref {
	if ref, ok := r.Refs[name]; ok {
		ref.Usages++

		return ref
	}

	return r.createRef(name)
}
