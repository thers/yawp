package ast

type RefKind int

const (
	RUnknown RefKind = iota
	RVar
	RLet
	RConst
	RImport
	RLabel
)

type Ref struct {
	Name   string
	Kind   RefKind
	Usages int

	ShadowsRef    *Ref
	ShadowedByRef *Ref

	Mangled bool
}
