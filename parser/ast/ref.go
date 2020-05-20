package ast

import "yawp/parser/token"

type SymbolRefType int

const (
	SRUnknown SymbolRefType = iota
	SRVar
	SRLet
	SRConst
	SRClass
	SRImport
	SRExport
	SRFn
	SRFnParam
	SRLabel
	SRBuiltin
)

type SymbolRef struct {
	Name   string
	Type   SymbolRefType
	Usages int

	ShadowsRef    *SymbolRef
	ShadowedByRef *SymbolRef

	Mangled bool
}

func SymbolRefTypeFromToken(value token.Token) SymbolRefType {
	switch value {
	case token.VAR:
		return SRVar
	case token.CONST:
		return SRConst
	case token.LET:
		return SRLet
	default:
		return SRUnknown
	}
}
