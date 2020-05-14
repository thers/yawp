package ast

func (i *Identifier) Copy() *Identifier {
	return &Identifier{
		ExprNode: i.ExprNode.Copy(),
		Ref:      i.Ref,
		Name:     i.Name,
	}
}
