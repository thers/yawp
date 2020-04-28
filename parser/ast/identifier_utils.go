package ast

func (i *Identifier) Clone() *Identifier {
	return &Identifier{
		Loc:  i.Loc,
		Ref:  i.Ref,
		Name: i.Name,
	}
}
