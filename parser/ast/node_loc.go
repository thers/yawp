package ast

import "yawp/parser/file"

func (i *Identifier) GetLoc() *file.Loc              { return i.Loc }
func (v *VariableStatement) GetLoc() *file.Loc       { return v.Loc }
func (s *SequenceExpression) GetLoc() *file.Loc      { return s.Loc }
func (d *UnaryExpression) GetLoc() *file.Loc         { return d.Loc }
func (d *SpreadExpression) GetLoc() *file.Loc        { return d.Loc }
func (d *ThisExpression) GetLoc() *file.Loc          { return d.Loc }
func (d *BlockStatement) GetLoc() *file.Loc          { return d.Loc }
func (d *AwaitExpression) GetLoc() *file.Loc         { return d.Loc }
func (d *ArrowFunctionExpression) GetLoc() *file.Loc { return d.Loc }
func (d *ReturnStatement) GetLoc() *file.Loc         { return d.Loc }
func (d *NullLiteral) GetLoc() *file.Loc             { return d.Loc }
func (d *BooleanLiteral) GetLoc() *file.Loc          { return d.Loc }
func (d *StringLiteral) GetLoc() *file.Loc           { return d.Loc }
func (d *NumberLiteral) GetLoc() *file.Loc           { return d.Loc }
func (d *RegExpLiteral) GetLoc() *file.Loc           { return d.Loc }
func (d *NewExpression) GetLoc() *file.Loc           { return d.Loc }

func (d *EmptyStatement) GetLoc() *file.Loc    { return d.Loc }
func (d *DebuggerStatement) GetLoc() *file.Loc { return d.Loc }
func (d *BranchStatement) GetLoc() *file.Loc   { return d.Loc }
func (d *SwitchStatement) GetLoc() *file.Loc   { return d.Loc }
func (d *CaseStatement) GetLoc() *file.Loc     { return d.Loc }
func (d *WithStatement) GetLoc() *file.Loc     { return d.Loc }
func (d *DoWhileStatement) GetLoc() *file.Loc  { return d.Loc.Add(d.Body.GetLoc()) }
func (d *IfStatement) GetLoc() *file.Loc {
	loc := d.Loc.Add(d.Consequent.GetLoc())

	if d.Alternate != nil {
		loc = loc.Add(d.Alternate.GetLoc())
	}

	return loc
}
func (d *ThrowStatement) GetLoc() *file.Loc      { return d.Loc }
func (d *WhileStatement) GetLoc() *file.Loc      { return d.Loc.Add(d.Body.GetLoc()) }
func (d *ExpressionStatement) GetLoc() *file.Loc { return d.Expression.GetLoc() }
func (d *LabelledStatement) GetLoc() *file.Loc   { return d.Statement.GetLoc() }
func (d *CatchStatement) GetLoc() *file.Loc      { return d.Loc.Add(d.Body.GetLoc()) }
func (d *TryStatement) GetLoc() *file.Loc {
	loc := d.Loc.Add(d.Body.GetLoc())

	if d.Catch != nil {
		loc = loc.Add(d.Catch.GetLoc())
	}

	if d.Finally != nil {
		loc = loc.Add(d.Finally.GetLoc())
	}

	return loc
}

func (d *BinaryExpression) GetLoc() *file.Loc {
	return d.Left.GetLoc().Add(d.Right.GetLoc())
}
func (d *AssignExpression) GetLoc() *file.Loc {
	return d.Left.GetLoc().Add(d.Right.GetLoc())
}
func (b *BracketExpression) GetLoc() *file.Loc { return b.Left.GetLoc().Add(b.Member.GetLoc()) }
func (d *DotExpression) GetLoc() *file.Loc     { return d.Left.GetLoc().Add(d.Identifier.GetLoc()) }
func (d *CallExpression) GetLoc() *file.Loc {
	if len(d.ArgumentList) < 1 {
		return d.Callee.GetLoc()
	}

	lastArgLoc := d.ArgumentList[len(d.ArgumentList)-1].GetLoc()

	return d.Callee.GetLoc().Add(lastArgLoc)
}
