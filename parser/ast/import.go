package ast

type ImportKind int

const (
	IKValue ImportKind = iota
	IKType
	IKTypeOf
)

type (
	ImportClause struct {
		ExprNode
		Namespace        bool
		ModuleIdentifier *Identifier
		LocalIdentifier  *Identifier
	}

	ImportCall struct {
		ExprNode
		Expression IExpr
	}
)
