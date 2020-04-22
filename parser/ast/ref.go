package ast

type RefKind int

const (
	RVar RefKind = iota
	RLet
	RConst
	RImport
)

type Ref struct {
	Name          string
	Kind          RefKind
	Usages        int
	ShadowsRef    *Ref
	ShadowedByRef *Ref

	IsMangled bool
}
