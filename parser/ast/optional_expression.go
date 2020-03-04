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
		End   file.Idx
	}

	OptionalCallExpression struct {
		Left      Expression
		Arguments []Expression
		End       file.Idx
	}
)

func (*OptionalObjectMemberAccessExpression) _expressionNode() {}
func (*OptionalArrayMemberAccessExpression) _expressionNode()  {}
func (*OptionalCallExpression) _expressionNode()               {}

func (s *OptionalObjectMemberAccessExpression) StartAt() file.Idx { return s.Left.StartAt() }
func (s *OptionalArrayMemberAccessExpression) StartAt() file.Idx  { return s.Left.StartAt() }
func (s *OptionalCallExpression) StartAt() file.Idx               { return s.Left.StartAt() }

func (s *OptionalObjectMemberAccessExpression) EndAt() file.Idx { return s.Identifier.EndAt() }
func (s *OptionalArrayMemberAccessExpression) EndAt() file.Idx  { return s.End }
func (s *OptionalCallExpression) EndAt() file.Idx               { return s.End }
