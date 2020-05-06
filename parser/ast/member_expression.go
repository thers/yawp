package ast

import "yawp/parser/file"

type MemberKind int

const (
	MKArray = iota
	MKObject
)

type (
	MemberExpression struct {
		Left  Expression
		Right Expression
		Kind  MemberKind
	}
)

func (m *MemberExpression) GetLoc() *file.Loc { return m.Left.GetLoc() }
func (*MemberExpression) _expressionNode()    {}
