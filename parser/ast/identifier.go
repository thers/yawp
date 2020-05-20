package ast

func (i *Identifier) Copy() *Identifier {
	return &Identifier{
		ExprNode:  i.ExprNode.Copy(),
		LegacyRef: i.LegacyRef,
		Name:      i.Name,
	}
}
