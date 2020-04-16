package ast

import "yawp/parser/file"

type (
	OptionalObjectMemberAccessExpression struct {
		Left       Expression
		Identifier *Identifier
	}

	OptionalArrayMemberAccessExpression struct {
		Loc   *file.Loc
		Left  Expression
		Index Expression
	}

	OptionalCallExpression struct {
		Loc       *file.Loc
		Left      Expression
		Arguments []Expression
	}
)

func (*OptionalObjectMemberAccessExpression) _expressionNode() {}
func (*OptionalArrayMemberAccessExpression) _expressionNode()  {}
func (*OptionalCallExpression) _expressionNode()               {}

func (s *OptionalObjectMemberAccessExpression) GetLoc() *file.Loc {
	return s.Left.GetLoc().Add(s.Identifier.GetLoc())
}
func (s *OptionalArrayMemberAccessExpression) GetLoc() *file.Loc { return s.Loc }
func (s *OptionalCallExpression) GetLoc() *file.Loc              { return s.Loc }
