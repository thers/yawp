package ast

type RefKind int

const (
	RUnknown RefKind = iota
	RVar
	RLet
	RConst
	RImport
	RFn
	RFnParam
	RLabel
	RBuiltin
)

type Ref struct {
	Name   string
	Kind   RefKind
	Usages int

	ShadowsRef    *Ref
	ShadowedByRef *Ref

	Mangled bool
}
