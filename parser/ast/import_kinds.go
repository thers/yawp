package ast

type ImportKind int

const (
	IKValue ImportKind = iota
	IKType
	IKTypeOf
)
