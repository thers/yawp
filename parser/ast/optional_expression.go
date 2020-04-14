package ast

import "yawp/parser/file"

type (
	OptionalObjectMemberAccessExpression struct {
		Left       Expression
		Identifier *Identifier
	}

	OptionalArrayMemberAccessExpression struct {
		Left  Expression
		Index Expression
		End   file.Loc
	}

	OptionalCallExpression struct {
		Left      Expression
		Arguments []Expression
		End       file.Loc
	}
)

func (*OptionalObjectMemberAccessExpression) _expressionNode() {}
func (*OptionalArrayMemberAccessExpression) _expressionNode()  {}
func (*OptionalCallExpression) _expressionNode()               {}

func (s *OptionalObjectMemberAccessExpression) StartAt() file.Loc { return s.Left.StartAt() }
func (s *OptionalArrayMemberAccessExpression) StartAt() file.Loc  { return s.Left.StartAt() }
func (s *OptionalCallExpression) StartAt() file.Loc               { return s.Left.StartAt() }

func (s *OptionalObjectMemberAccessExpression) EndAt() file.Loc { return s.Identifier.EndAt() }
func (s *OptionalArrayMemberAccessExpression) EndAt() file.Loc  { return s.End }
func (s *OptionalCallExpression) EndAt() file.Loc               { return s.End }
