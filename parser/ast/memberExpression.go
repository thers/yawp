package ast

type (
	MemberExpression interface {
		Expression
		_memberExpressionNode()
	}
)

func (d *DotExpression) _memberExpressionNode()     {}
func (b *BracketExpression) _memberExpressionNode() {}
func (i *Identifier) _memberExpressionNode()        {}
