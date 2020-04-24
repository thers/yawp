package ast

type Visitor interface {
	// statements
	Statement(stmt Statement) Statement
	ExportDeclaration(stmt *ExportDeclaration) *ExportDeclaration
	ImportDeclaration(stmt *ImportDeclaration) *ImportDeclaration
	FlowTypeStatement(stmt *FlowTypeStatement) *FlowTypeStatement
	FlowInterfaceStatement(stmt *FlowInterfaceStatement) *FlowInterfaceStatement
	BlockStatement(stmt *BlockStatement) *BlockStatement
	ClassStatement(stmt *ClassStatement) *ClassStatement
	ClassFieldStatement(stmt *ClassFieldStatement) *ClassFieldStatement
	ClassAccessorStatement(stmt *ClassAccessorStatement) *ClassAccessorStatement
	ClassMethodStatement(stmt *ClassMethodStatement) *ClassMethodStatement
	LegacyDecoratorStatement(stmt *LegacyDecoratorStatement) *LegacyDecoratorStatement
	ForInStatement(stmt *ForInStatement) *ForInStatement
	ForOfStatement(stmt *ForOfStatement) *ForOfStatement
	ForStatement(stmt *ForStatement) *ForStatement
	BranchStatement(stmt *BranchStatement) *BranchStatement
	CatchStatement(stmt *CatchStatement) *CatchStatement
	DebuggerStatement(stmt *DebuggerStatement) *DebuggerStatement
	DoWhileStatement(stmt *DoWhileStatement) *DoWhileStatement
	EmptyStatement(stmt *EmptyStatement) *EmptyStatement
	ExpressionStatement(stmt *ExpressionStatement) *ExpressionStatement
	IfStatement(stmt *IfStatement) *IfStatement
	LabelledStatement(stmt *LabelledStatement) *LabelledStatement
	ReturnStatement(stmt *ReturnStatement) *ReturnStatement
	SwitchStatement(stmt *SwitchStatement) *SwitchStatement
	CaseStatement(stmt *CaseStatement) *CaseStatement
	ThrowStatement(stmt *ThrowStatement) *ThrowStatement
	TryStatement(stmt *TryStatement) *TryStatement
	VariableStatement(stmt *VariableStatement) *VariableStatement
	WhileStatement(stmt *WhileStatement) *WhileStatement
	WithStatement(stmt *WithStatement) *WithStatement
	YieldStatement(stmt *YieldStatement) *YieldStatement

	// expressions
	Expression(exp Expression) Expression
	Identifier(exp *Identifier) *Identifier

	ImportClause(exp *ImportClause) *ImportClause
	ImportCall(exp *ImportCall) *ImportCall

	FlowTypeAssertion(exp *FlowTypeAssertion) *FlowTypeAssertion

	FunctionLiteral(stmt *FunctionLiteral) *FunctionLiteral
	StringLiteral(exp *StringLiteral) *StringLiteral
	BooleanLiteral(exp *BooleanLiteral) *BooleanLiteral
	ObjectLiteral(exp *ObjectLiteral) *ObjectLiteral
	ArrayLiteral(exp *ArrayLiteral) *ArrayLiteral
	NullLiteral(exp *NullLiteral) *NullLiteral
	NumberLiteral(exp *NumberLiteral) *NumberLiteral
	RegExpLiteral(exp *RegExpLiteral) *RegExpLiteral

	ObjectSpread(exp *ObjectSpread) *ObjectSpread
	ArraySpread(exp *ArraySpread) *ArraySpread

	VariableBinding(exp *VariableBinding) *VariableBinding

	ClassExpression(exp *ClassExpression) *ClassExpression
	ClassSuperExpression(exp *ClassSuperExpression) *ClassSuperExpression
	CoalesceExpression(exp *CoalesceExpression) *CoalesceExpression
	ConditionalExpression(exp *ConditionalExpression) *ConditionalExpression
	JsxElement(exp *JSXElement) *JSXElement
	JsxFragment(exp *JSXFragment) *JSXFragment
	NewTargetExpression(exp *NewTargetExpression) *NewTargetExpression
	AssignExpression(exp *AssignExpression) *AssignExpression
	BinaryExpression(exp *BinaryExpression) *BinaryExpression
	BracketExpression(exp *BracketExpression) *BracketExpression
	CallExpression(exp *CallExpression) *CallExpression
	DotExpression(exp *DotExpression) *DotExpression
	SpreadExpression(exp *SpreadExpression) *SpreadExpression
	NewExpression(exp *NewExpression) *NewExpression
	SequenceExpression(exp *SequenceExpression) *SequenceExpression
	ThisExpression(exp *ThisExpression) *ThisExpression
	UnaryExpression(exp *UnaryExpression) *UnaryExpression
	ArrowFunctionExpression(exp *ArrowFunctionExpression) *ArrowFunctionExpression
	AwaitExpression(exp *AwaitExpression) *AwaitExpression
	OptionalObjectMemberAccessExpression(exp *OptionalObjectMemberAccessExpression) *OptionalObjectMemberAccessExpression
	OptionalArrayMemberAccessExpression(exp *OptionalArrayMemberAccessExpression) *OptionalArrayMemberAccessExpression
	OptionalCallExpression(exp *OptionalCallExpression) *OptionalCallExpression
	TemplateExpression(exp *TemplateExpression) *TemplateExpression
	TaggedTemplateExpression(exp *TaggedTemplateExpression) *TaggedTemplateExpression
	YieldExpression(exp *YieldExpression) *YieldExpression

	// others
	MemberExpression(exp MemberExpression) MemberExpression
	ClassFieldName(name ClassFieldName) ClassFieldName
	PatternBinder(binder PatternBinder) PatternBinder
	IdentifierBinder(b *IdentifierBinder) *IdentifierBinder
	ObjectRestBinder(b *ObjectRestBinder) *ObjectRestBinder
	ArrayRestBinder(b *ArrayRestBinder) *ArrayRestBinder
	ObjectPropertyBinder(b *ObjectPropertyBinder) *ObjectPropertyBinder
	ArrayItemBinder(b *ArrayItemBinder) *ArrayItemBinder
	ArrayBinding(b *ArrayBinding) *ArrayBinding
	ObjectBinding(b *ObjectBinding) *ObjectBinding
	FunctionParameters(params *FunctionParameters) *FunctionParameters
	FunctionParameter(param FunctionParameter) FunctionParameter
	IdentifierParameter(p *IdentifierParameter) *IdentifierParameter
	RestParameter(p *RestParameter) *RestParameter
	ObjectPatternParameter(p *ObjectPatternParameter) *ObjectPatternParameter
	ObjectPatternIdentifierParameter(p *ObjectPatternIdentifierParameter) *ObjectPatternIdentifierParameter
	ArrayPatternParameter(p *ArrayPatternParameter) *ArrayPatternParameter
	ObjectProperty(p ObjectProperty) ObjectProperty
	ObjectPropertySetter(p *ObjectPropertySetter) *ObjectPropertySetter
	ObjectPropertyGetter(p *ObjectPropertyGetter) *ObjectPropertyGetter
	ObjectPropertyValue(p *ObjectPropertyValue) *ObjectPropertyValue
	ObjectPropertyName(n ObjectPropertyName) ObjectPropertyName
	ComputedName(n *ComputedName) *ComputedName
	LegacyDecoratorSubject(s LegacyDecoratorSubject) LegacyDecoratorSubject
	ExportClause(c ExportClause) ExportClause
	ExportNamespaceFromClause(c *ExportNamespaceFromClause) *ExportNamespaceFromClause
	ExportNamedFromClause(c *ExportNamedFromClause) *ExportNamedFromClause
	ExportNamedClause(c *ExportNamedClause) *ExportNamedClause
	ExportVarClause(c *ExportVarClause) *ExportVarClause
	ExportFunctionClause(c *ExportFunctionClause) *ExportFunctionClause
	ExportClassClause(c *ExportClassClause) *ExportClassClause
	ExportDefaultClause(c *ExportDefaultClause) *ExportDefaultClause
}

type DefaultVisitor struct {
	Specific Visitor
}

func (d *DefaultVisitor) ExportDeclaration(stmt *ExportDeclaration) *ExportDeclaration {
	stmt.Clause = d.Specific.ExportClause(stmt.Clause)

	return stmt
}

func (d *DefaultVisitor) ExportClause(ac ExportClause) ExportClause {
	if ac == nil {
		return nil
	}

	switch c := ac.(type) {
	case *ExportNamespaceFromClause:
		return d.Specific.ExportNamespaceFromClause(c)
	case *ExportNamedFromClause:
		return d.Specific.ExportNamedFromClause(c)
	case *ExportNamedClause:
		return d.Specific.ExportNamedClause(c)
	case *ExportVarClause:
		return d.Specific.ExportVarClause(c)
	case *ExportFunctionClause:
		return d.Specific.ExportFunctionClause(c)
	case *ExportClassClause:
		return d.Specific.ExportClassClause(c)
	case *ExportDefaultClause:
		return d.Specific.ExportDefaultClause(c)
	case *FlowTypeStatement:
		return d.Specific.FlowTypeStatement(c)
	case *FlowInterfaceStatement:
		return d.Specific.FlowInterfaceStatement(c)

	default:
		panic("Unknown export clause type")
	}
}

func (d *DefaultVisitor) ExportNamespaceFromClause(c *ExportNamespaceFromClause) *ExportNamespaceFromClause {
	panic("implement me")
}

func (d *DefaultVisitor) ExportNamedFromClause(c *ExportNamedFromClause) *ExportNamedFromClause {
	panic("implement me")
}

func (d *DefaultVisitor) ExportNamedClause(c *ExportNamedClause) *ExportNamedClause {
	panic("implement me")
}

func (d *DefaultVisitor) ExportVarClause(c *ExportVarClause) *ExportVarClause {
	panic("implement me")
}

func (d *DefaultVisitor) ExportFunctionClause(c *ExportFunctionClause) *ExportFunctionClause {
	panic("implement me")
}

func (d *DefaultVisitor) ExportClassClause(c *ExportClassClause) *ExportClassClause {
	panic("implement me")
}

func (d *DefaultVisitor) ExportDefaultClause(c *ExportDefaultClause) *ExportDefaultClause {
	panic("implement me")
}

func (d *DefaultVisitor) ImportDeclaration(stmt *ImportDeclaration) *ImportDeclaration {
	for index, i := range stmt.Imports {
		stmt.Imports[index] = d.Specific.ImportClause(i)
	}

	return stmt
}

func (d *DefaultVisitor) ImportClause(exp *ImportClause) *ImportClause {
	exp.LocalIdentifier = d.Specific.Identifier(exp.LocalIdentifier)
	exp.ModuleIdentifier = d.Specific.Identifier(exp.ModuleIdentifier)

	return exp
}

func (d *DefaultVisitor) ImportCall(exp *ImportCall) *ImportCall {
	exp.Expression = d.Specific.Expression(exp.Expression)

	return exp
}

func (d *DefaultVisitor) FlowTypeStatement(stmt *FlowTypeStatement) *FlowTypeStatement {
	return stmt
}

func (d *DefaultVisitor) FlowInterfaceStatement(stmt *FlowInterfaceStatement) *FlowInterfaceStatement {
	return stmt
}

func (d *DefaultVisitor) FlowTypeAssertion(exp *FlowTypeAssertion) *FlowTypeAssertion {
	return exp
}

func (d *DefaultVisitor) Statement(stmt Statement) Statement {
	if stmt == nil {
		return nil
	}

	switch s := stmt.(type) {
	case *ExportDeclaration:
		return d.Specific.ExportDeclaration(s)
	case *ImportDeclaration:
		return d.Specific.ImportDeclaration(s)
	case *FlowTypeStatement:
		return d.Specific.FlowTypeStatement(s)
	case *FlowInterfaceStatement:
		return d.Specific.FlowInterfaceStatement(s)
	case *BlockStatement:
		return d.Specific.BlockStatement(s)
	case *ClassStatement:
		return d.Specific.ClassStatement(s)
	case *ClassFieldStatement:
		return d.Specific.ClassFieldStatement(s)
	case *ClassAccessorStatement:
		return d.Specific.ClassAccessorStatement(s)
	case *ClassMethodStatement:
		return d.Specific.ClassMethodStatement(s)
	case *LegacyDecoratorStatement:
		return d.Specific.LegacyDecoratorStatement(s)
	case *ForInStatement:
		return d.Specific.ForInStatement(s)
	case *ForOfStatement:
		return d.Specific.ForOfStatement(s)
	case *ForStatement:
		return d.Specific.ForStatement(s)
	case *BranchStatement:
		return d.Specific.BranchStatement(s)
	case *CatchStatement:
		return d.Specific.CatchStatement(s)
	case *DebuggerStatement:
		return d.Specific.DebuggerStatement(s)
	case *DoWhileStatement:
		return d.Specific.DoWhileStatement(s)
	case *EmptyStatement:
		return d.Specific.EmptyStatement(s)
	case *ExpressionStatement:
		return d.Specific.ExpressionStatement(s)
	case *IfStatement:
		return d.Specific.IfStatement(s)
	case *LabelledStatement:
		return d.Specific.LabelledStatement(s)
	case *ReturnStatement:
		return d.Specific.ReturnStatement(s)
	case *SwitchStatement:
		return d.Specific.SwitchStatement(s)
	case *ThrowStatement:
		return d.Specific.ThrowStatement(s)
	case *TryStatement:
		return d.Specific.TryStatement(s)
	case *VariableStatement:
		return d.Specific.VariableStatement(s)
	case *WhileStatement:
		return d.Specific.WhileStatement(s)
	case *WithStatement:
		return d.Specific.WithStatement(s)
	case *YieldStatement:
		return d.Specific.YieldStatement(s)
	case *FunctionLiteral:
		return d.Specific.FunctionLiteral(s)

	default:
		panic("Unknown statement")
	}
}

func (d *DefaultVisitor) Expression(exp Expression) Expression {
	if exp == nil {
		return nil
	}

	switch s := exp.(type) {
	case *Identifier:
		return d.Specific.Identifier(s)
	case *ImportClause:
		return d.Specific.ImportClause(s)
	case *ImportCall:
		return d.Specific.ImportCall(s)
	case *FlowTypeAssertion:
		return d.Specific.FlowTypeAssertion(s)
	case *FunctionLiteral:
		return d.Specific.FunctionLiteral(s)
	case *StringLiteral:
		return d.Specific.StringLiteral(s)
	case *BooleanLiteral:
		return d.Specific.BooleanLiteral(s)
	case *ObjectLiteral:
		return d.Specific.ObjectLiteral(s)
	case *ArrayLiteral:
		return d.Specific.ArrayLiteral(s)
	case *NullLiteral:
		return d.Specific.NullLiteral(s)
	case *NumberLiteral:
		return d.Specific.NumberLiteral(s)
	case *RegExpLiteral:
		return d.Specific.RegExpLiteral(s)
	case *ObjectSpread:
		return d.Specific.ObjectSpread(s)
	case *ArraySpread:
		return d.Specific.ArraySpread(s)
	case *ArrayBinding:
		return d.Specific.ArrayBinding(s)
	case *ObjectBinding:
		return d.Specific.ObjectBinding(s)
	case *VariableBinding:
		return d.Specific.VariableBinding(s)
	case *ClassExpression:
		return d.Specific.ClassExpression(s)
	case *ClassSuperExpression:
		return d.Specific.ClassSuperExpression(s)
	case *CoalesceExpression:
		return d.Specific.CoalesceExpression(s)
	case *ConditionalExpression:
		return d.Specific.ConditionalExpression(s)
	case *JSXElement:
		return d.Specific.JsxElement(s)
	case *JSXFragment:
		return d.Specific.JsxFragment(s)
	case *NewTargetExpression:
		return d.Specific.NewTargetExpression(s)
	case *AssignExpression:
		return d.Specific.AssignExpression(s)
	case *BinaryExpression:
		return d.Specific.BinaryExpression(s)
	case *BracketExpression:
		return d.Specific.BracketExpression(s)
	case *CallExpression:
		return d.Specific.CallExpression(s)
	case *DotExpression:
		return d.Specific.DotExpression(s)
	case *SpreadExpression:
		return d.Specific.SpreadExpression(s)
	case *NewExpression:
		return d.Specific.NewExpression(s)
	case *SequenceExpression:
		return d.Specific.SequenceExpression(s)
	case *ThisExpression:
		return d.Specific.ThisExpression(s)
	case *UnaryExpression:
		return d.Specific.UnaryExpression(s)
	case *ArrowFunctionExpression:
		return d.Specific.ArrowFunctionExpression(s)
	case *AwaitExpression:
		return d.Specific.AwaitExpression(s)
	case *OptionalObjectMemberAccessExpression:
		return d.Specific.OptionalObjectMemberAccessExpression(s)
	case *OptionalArrayMemberAccessExpression:
		return d.Specific.OptionalArrayMemberAccessExpression(s)
	case *OptionalCallExpression:
		return d.Specific.OptionalCallExpression(s)
	case *TemplateExpression:
		return d.Specific.TemplateExpression(s)
	case *TaggedTemplateExpression:
		return d.Specific.TaggedTemplateExpression(s)
	case *YieldExpression:
		return d.Specific.YieldExpression(s)

	default:
		panic("Unknown expression")
	}
}

func (d *DefaultVisitor) BlockStatement(block *BlockStatement) *BlockStatement {
	for index, stmt := range block.List {
		block.List[index] = d.Specific.Statement(stmt)
	}

	return block
}

func (d *DefaultVisitor) ClassStatement(stmt *ClassStatement) *ClassStatement {
	stmt.Expression = d.Specific.ClassExpression(stmt.Expression)

	return stmt
}

func (d *DefaultVisitor) ClassFieldStatement(stmt *ClassFieldStatement) *ClassFieldStatement {
	stmt.Initializer = d.Specific.Expression(stmt.Initializer)
	stmt.Name = d.Specific.ClassFieldName(stmt.Name)

	return stmt
}

func (d *DefaultVisitor) ClassFieldName(name ClassFieldName) ClassFieldName {
	switch n := name.(type) {
	case *ComputedName:
		return d.Specific.ComputedName(n)
	case *Identifier:
		return d.Specific.Identifier(n)

	default:
		return n
	}
}

func (d *DefaultVisitor) ComputedName(n *ComputedName) *ComputedName {
	n.Expression = d.Specific.Expression(n.Expression)
	return n
}

func (d *DefaultVisitor) ClassAccessorStatement(stmt *ClassAccessorStatement) *ClassAccessorStatement {
	stmt.Body = d.Specific.FunctionLiteral(stmt.Body)
	stmt.Field = d.Specific.ClassFieldName(stmt.Field)

	return stmt
}

func (d *DefaultVisitor) ClassMethodStatement(stmt *ClassMethodStatement) *ClassMethodStatement {
	stmt.Body = d.Specific.Statement(stmt.Body)
	stmt.Name = d.Specific.ClassFieldName(stmt.Name)
	stmt.Parameters = d.Specific.FunctionParameters(stmt.Parameters)

	return stmt
}

func (d *DefaultVisitor) LegacyDecoratorStatement(stmt *LegacyDecoratorStatement) *LegacyDecoratorStatement {
	for index, dec := range stmt.Decorators {
		stmt.Decorators[index] = d.Specific.Expression(dec)
	}

	stmt.Subject = d.Specific.LegacyDecoratorSubject(stmt.Subject)

	return stmt
}

func (d *DefaultVisitor) LegacyDecoratorSubject(as LegacyDecoratorSubject) LegacyDecoratorSubject {
	switch s := as.(type) {
	case *ClassStatement:
		return d.Specific.ClassStatement(s)
	case *ClassFieldStatement:
		return d.Specific.ClassFieldStatement(s)
	case *ClassMethodStatement:
		return d.Specific.ClassMethodStatement(s)

	default:
		panic("Unknown legacy decorator subject type")
	}
}

func (d *DefaultVisitor) ForInStatement(stmt *ForInStatement) *ForInStatement {
	stmt.Into = d.Specific.Expression(stmt.Into)
	stmt.Source = d.Specific.Expression(stmt.Source)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) ForOfStatement(stmt *ForOfStatement) *ForOfStatement {
	stmt.Binder = d.Specific.Expression(stmt.Binder)
	stmt.Iterator = d.Specific.Expression(stmt.Iterator)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) ForStatement(stmt *ForStatement) *ForStatement {
	stmt.Initializer = d.Specific.VariableStatement(stmt.Initializer)
	stmt.Test = d.Specific.Expression(stmt.Test)
	stmt.Update = d.Specific.Expression(stmt.Update)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) BranchStatement(stmt *BranchStatement) *BranchStatement {
	stmt.Label = d.Specific.Identifier(stmt.Label)

	return stmt
}

func (d *DefaultVisitor) CatchStatement(stmt *CatchStatement) *CatchStatement {
	if stmt == nil {
		return nil
	}

	stmt.Parameter = d.Specific.Identifier(stmt.Parameter)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) DebuggerStatement(stmt *DebuggerStatement) *DebuggerStatement {
	return stmt
}

func (d *DefaultVisitor) DoWhileStatement(stmt *DoWhileStatement) *DoWhileStatement {
	stmt.Body = d.Specific.Statement(stmt.Body)
	stmt.Test = d.Specific.Expression(stmt.Test)

	return stmt
}

func (d *DefaultVisitor) EmptyStatement(stmt *EmptyStatement) *EmptyStatement {
	return stmt
}

func (d *DefaultVisitor) ExpressionStatement(stmt *ExpressionStatement) *ExpressionStatement {
	stmt.Expression = d.Specific.Expression(stmt.Expression)

	return stmt
}

func (d *DefaultVisitor) IfStatement(stmt *IfStatement) *IfStatement {
	stmt.Test = d.Specific.Expression(stmt.Test)
	stmt.Consequent = d.Specific.Statement(stmt.Consequent)
	stmt.Alternate = d.Specific.Statement(stmt.Alternate)

	return stmt
}

func (d *DefaultVisitor) LabelledStatement(stmt *LabelledStatement) *LabelledStatement {
	stmt.Label = d.Specific.Identifier(stmt.Label)
	stmt.Statement = d.Specific.Statement(stmt.Statement)

	return stmt
}

func (d *DefaultVisitor) ReturnStatement(stmt *ReturnStatement) *ReturnStatement {
	stmt.Argument = d.Specific.Expression(stmt.Argument)

	return stmt
}

func (d *DefaultVisitor) SwitchStatement(stmt *SwitchStatement) *SwitchStatement {
	stmt.Discriminant = d.Specific.Expression(stmt.Discriminant)

	for index, c := range stmt.Body {
		stmt.Body[index] = d.Specific.CaseStatement(c)
	}

	return stmt
}

func (d *DefaultVisitor) CaseStatement(s *CaseStatement) *CaseStatement {
	s.Test = d.Specific.Expression(s.Test)

	for index, st := range s.Consequent {
		s.Consequent[index] = d.Specific.Statement(st)
	}

	return s
}

func (d *DefaultVisitor) ThrowStatement(stmt *ThrowStatement) *ThrowStatement {
	stmt.Argument = d.Specific.Expression(stmt.Argument)

	return stmt
}

func (d *DefaultVisitor) TryStatement(stmt *TryStatement) *TryStatement {
	stmt.Body = d.Specific.Statement(stmt.Body)
	stmt.Catch = d.Specific.CatchStatement(stmt.Catch)
	stmt.Finally = d.Specific.Statement(stmt.Finally)

	return stmt
}

func (d *DefaultVisitor) VariableStatement(stmt *VariableStatement) *VariableStatement {
	for index, v := range stmt.List {
		stmt.List[index] = d.Specific.VariableBinding(v)
	}

	return stmt
}

func (d *DefaultVisitor) WhileStatement(stmt *WhileStatement) *WhileStatement {
	stmt.Test = d.Specific.Expression(stmt.Test)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) WithStatement(stmt *WithStatement) *WithStatement {
	stmt.Object = d.Specific.Expression(stmt.Object)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) YieldStatement(stmt *YieldStatement) *YieldStatement {
	stmt.Expression = d.Specific.YieldExpression(stmt.Expression)

	return stmt
}

func (d *DefaultVisitor) Identifier(exp *Identifier) *Identifier {
	if exp == nil {
		return nil
	}

	return exp
}

func (d *DefaultVisitor) FunctionLiteral(stmt *FunctionLiteral) *FunctionLiteral {
	stmt.Name = d.Specific.Identifier(stmt.Name)
	stmt.Parameters = d.Specific.FunctionParameters(stmt.Parameters)
	stmt.Body = d.Specific.Statement(stmt.Body)

	return stmt
}

func (d *DefaultVisitor) StringLiteral(exp *StringLiteral) *StringLiteral {
	return exp
}

func (d *DefaultVisitor) BooleanLiteral(exp *BooleanLiteral) *BooleanLiteral {
	return exp
}

func (d *DefaultVisitor) ObjectLiteral(exp *ObjectLiteral) *ObjectLiteral {
	for index, prop := range exp.Properties {
		exp.Properties[index] = d.Specific.ObjectProperty(prop)
	}

	return exp
}

func (d *DefaultVisitor) ObjectProperty(ap ObjectProperty) ObjectProperty {
	if ap == nil {
		return nil
	}

	switch p := ap.(type) {
	case *ObjectPropertySetter:
		return d.Specific.ObjectPropertySetter(p)
	case *ObjectPropertyGetter:
		return d.Specific.ObjectPropertyGetter(p)
	case *ObjectPropertyValue:
		return d.Specific.ObjectPropertyValue(p)
	case *ObjectSpread:
		return d.Specific.ObjectSpread(p)

	default:
		panic("Unknown object property type")
	}
}

func (d *DefaultVisitor) ObjectPropertySetter(p *ObjectPropertySetter) *ObjectPropertySetter {
	p.PropertyName = d.Specific.ObjectPropertyName(p.PropertyName)
	p.Setter = d.Specific.FunctionLiteral(p.Setter)

	return p
}

func (d *DefaultVisitor) ObjectPropertyGetter(p *ObjectPropertyGetter) *ObjectPropertyGetter {
	p.PropertyName = d.Specific.ObjectPropertyName(p.PropertyName)
	p.Getter = d.Specific.FunctionLiteral(p.Getter)

	return p
}

func (d *DefaultVisitor) ObjectPropertyValue(p *ObjectPropertyValue) *ObjectPropertyValue {
	p.PropertyName = d.Specific.ObjectPropertyName(p.PropertyName)
	p.Value = d.Specific.Expression(p.Value)

	return p
}

func (d *DefaultVisitor) ObjectPropertyName(an ObjectPropertyName) ObjectPropertyName {
	switch n := an.(type) {
	case *ComputedName:
		return d.Specific.ComputedName(n)
	case *Identifier:
		return d.Specific.Identifier(n)

	default:
		panic("Unknown object property name type")
	}
}

func (d *DefaultVisitor) ArrayLiteral(exp *ArrayLiteral) *ArrayLiteral {
	for index, item := range exp.Value {
		exp.Value[index] = d.Specific.Expression(item)
	}

	return exp
}

func (d *DefaultVisitor) NullLiteral(exp *NullLiteral) *NullLiteral {
	return exp
}

func (d *DefaultVisitor) NumberLiteral(exp *NumberLiteral) *NumberLiteral {
	return exp
}

func (d *DefaultVisitor) RegExpLiteral(exp *RegExpLiteral) *RegExpLiteral {
	return exp
}

func (d *DefaultVisitor) ObjectSpread(exp *ObjectSpread) *ObjectSpread {
	exp.Expression = d.Specific.Expression(exp.Expression)

	return exp
}

func (d *DefaultVisitor) ArraySpread(exp *ArraySpread) *ArraySpread {
	exp.Expression = d.Specific.Expression(exp.Expression)

	return exp
}

func (d *DefaultVisitor) ArrayBinding(exp *ArrayBinding) *ArrayBinding {
	for index, b := range exp.List {
		exp.List[index] = d.Specific.PatternBinder(b)
	}

	return exp
}

func (d *DefaultVisitor) ObjectBinding(exp *ObjectBinding) *ObjectBinding {
	for index, b := range exp.List {
		exp.List[index] = d.Specific.PatternBinder(b)
	}

	return exp
}

func (d *DefaultVisitor) VariableBinding(exp *VariableBinding) *VariableBinding {
	exp.Binder = d.Specific.PatternBinder(exp.Binder)
	exp.Initializer = d.Specific.Expression(exp.Initializer)

	return exp
}

func (d *DefaultVisitor) ClassExpression(exp *ClassExpression) *ClassExpression {
	exp.Name = d.Specific.Identifier(exp.Name)
	exp.SuperClass = d.Specific.MemberExpression(exp.SuperClass)
	exp.Body = d.Specific.Statement(exp.Body)

	return exp
}

func (d *DefaultVisitor) MemberExpression(exp MemberExpression) MemberExpression {
	if exp == nil {
		return nil
	}

	switch e := exp.(type) {
	case *BracketExpression:
		return d.Specific.BracketExpression(e)
	case *DotExpression:
		return d.Specific.DotExpression(e)
	case *Identifier:
		return d.Specific.Identifier(e)

	default:
		panic("Unknown member expression type")
	}
}

func (d *DefaultVisitor) ClassSuperExpression(exp *ClassSuperExpression) *ClassSuperExpression {
	for index, arg := range exp.Arguments {
		exp.Arguments[index] = d.Specific.Expression(arg)
	}

	return exp
}

func (d *DefaultVisitor) CoalesceExpression(exp *CoalesceExpression) *CoalesceExpression {
	exp.Head = d.Specific.Expression(exp.Head)
	exp.Consequent = d.Specific.Expression(exp.Consequent)

	return exp
}

func (d *DefaultVisitor) ConditionalExpression(exp *ConditionalExpression) *ConditionalExpression {
	exp.Test = d.Specific.Expression(exp.Test)
	exp.Consequent = d.Specific.Expression(exp.Consequent)
	exp.Alternate = d.Specific.Expression(exp.Alternate)

	return exp
}

func (d *DefaultVisitor) JsxElement(exp *JSXElement) *JSXElement {
	return exp
}

func (d *DefaultVisitor) JsxFragment(exp *JSXFragment) *JSXFragment {
	return exp
}

func (d *DefaultVisitor) NewTargetExpression(exp *NewTargetExpression) *NewTargetExpression {
	return exp
}

func (d *DefaultVisitor) AssignExpression(exp *AssignExpression) *AssignExpression {
	exp.Left = d.Specific.Expression(exp.Left)
	exp.Right = d.Specific.Expression(exp.Right)

	return exp
}

func (d *DefaultVisitor) BinaryExpression(exp *BinaryExpression) *BinaryExpression {
	exp.Left = d.Specific.Expression(exp.Left)
	exp.Right = d.Specific.Expression(exp.Right)

	return exp
}

func (d *DefaultVisitor) BracketExpression(exp *BracketExpression) *BracketExpression {
	exp.Left = d.Specific.Expression(exp.Left)
	exp.Member = d.Specific.Expression(exp.Member)

	return exp
}

func (d *DefaultVisitor) CallExpression(exp *CallExpression) *CallExpression {
	exp.Callee = d.Specific.Expression(exp.Callee)

	for index, arg := range exp.ArgumentList {
		exp.ArgumentList[index] = d.Specific.Expression(arg)
	}

	return exp
}

func (d *DefaultVisitor) DotExpression(exp *DotExpression) *DotExpression {
	exp.Left = d.Specific.Expression(exp.Left)
	exp.Identifier = d.Specific.Identifier(exp.Identifier)

	return exp
}

func (d *DefaultVisitor) SpreadExpression(exp *SpreadExpression) *SpreadExpression {
	exp.Value = d.Specific.Expression(exp.Value)

	return exp
}

func (d *DefaultVisitor) NewExpression(exp *NewExpression) *NewExpression {
	exp.Callee = d.Specific.Expression(exp.Callee)

	for index, arg := range exp.ArgumentList {
		exp.ArgumentList[index] = d.Specific.Expression(arg)
	}

	return exp
}

func (d *DefaultVisitor) SequenceExpression(exp *SequenceExpression) *SequenceExpression {
	for index, e := range exp.Sequence {
		exp.Sequence[index] = d.Specific.Expression(e)
	}

	return exp
}

func (d *DefaultVisitor) ThisExpression(exp *ThisExpression) *ThisExpression {
	return exp
}

func (d *DefaultVisitor) UnaryExpression(exp *UnaryExpression) *UnaryExpression {
	exp.Operand = d.Specific.Expression(exp.Operand)

	return exp
}

func (d *DefaultVisitor) ArrowFunctionExpression(exp *ArrowFunctionExpression) *ArrowFunctionExpression {
	exp.Body = d.Specific.Statement(exp.Body)

	for index, param := range exp.Parameters {
		exp.Parameters[index] = d.Specific.FunctionParameter(param)
	}

	return exp
}

func (d *DefaultVisitor) AwaitExpression(exp *AwaitExpression) *AwaitExpression {
	exp.Expression = d.Specific.Expression(exp.Expression)

	return exp
}

func (d *DefaultVisitor) OptionalObjectMemberAccessExpression(exp *OptionalObjectMemberAccessExpression) *OptionalObjectMemberAccessExpression {
	exp.Left = d.Specific.Expression(exp.Left)
	exp.Identifier = d.Specific.Identifier(exp.Identifier)

	return exp
}

func (d *DefaultVisitor) OptionalArrayMemberAccessExpression(exp *OptionalArrayMemberAccessExpression) *OptionalArrayMemberAccessExpression {
	exp.Left = d.Specific.Expression(exp.Left)
	exp.Index = d.Specific.Expression(exp.Index)

	return exp
}

func (d *DefaultVisitor) OptionalCallExpression(exp *OptionalCallExpression) *OptionalCallExpression {
	exp.Left = d.Specific.Expression(exp.Left)

	for index, arg := range exp.Arguments {
		exp.Arguments[index] = d.Specific.Expression(arg)
	}

	return exp
}

func (d *DefaultVisitor) TemplateExpression(exp *TemplateExpression) *TemplateExpression {
	for index, s := range exp.Substitutions {
		exp.Substitutions[index] = d.Specific.Expression(s)
	}

	return exp
}

func (d *DefaultVisitor) TaggedTemplateExpression(exp *TaggedTemplateExpression) *TaggedTemplateExpression {
	exp.Tag = d.Specific.Expression(exp.Tag)
	exp.Template = d.Specific.TemplateExpression(exp.Template)

	return exp
}

func (d *DefaultVisitor) YieldExpression(exp *YieldExpression) *YieldExpression {
	exp.Expression = d.Specific.Expression(exp.Expression)

	return exp
}

func (d *DefaultVisitor) FunctionParameters(params *FunctionParameters) *FunctionParameters {
	for index, param := range params.List {
		params.List[index] = d.Specific.FunctionParameter(param)
	}

	return params
}

func (d *DefaultVisitor) FunctionParameter(param FunctionParameter) FunctionParameter {
	switch p := param.(type) {
	case *IdentifierParameter:
		return d.Specific.IdentifierParameter(p)
	case *RestParameter:
		return d.Specific.RestParameter(p)
	case *ObjectPatternParameter:
		return d.Specific.ObjectPatternParameter(p)
	case *ArrayPatternParameter:
		return d.Specific.ArrayPatternParameter(p)

	default:
		panic("Unknown function parameter type")
	}
}

func (d *DefaultVisitor) IdentifierParameter(p *IdentifierParameter) *IdentifierParameter {
	p.Name = d.Specific.Identifier(p.Name)
	p.DefaultValue = d.Specific.Expression(p.DefaultValue)

	return p
}

func (d *DefaultVisitor) RestParameter(p *RestParameter) *RestParameter {
	p.Binder = d.Specific.PatternBinder(p.Binder)

	return p
}

func (d *DefaultVisitor) ObjectPatternParameter(p *ObjectPatternParameter) *ObjectPatternParameter {
	p.DefaultValue = d.Specific.Expression(p.DefaultValue)

	for index, prop := range p.List {
		p.List[index] = d.Specific.ObjectPatternIdentifierParameter(prop)
	}

	return p
}

func (d *DefaultVisitor) ObjectPatternIdentifierParameter(p *ObjectPatternIdentifierParameter) *ObjectPatternIdentifierParameter {
	p.Parameter = d.Specific.FunctionParameter(p.Parameter)

	return p
}

func (d *DefaultVisitor) ArrayPatternParameter(p *ArrayPatternParameter) *ArrayPatternParameter {
	p.DefaultValue = d.Specific.Expression(p.DefaultValue)

	for index, item := range p.List {
		p.List[index] = d.Specific.FunctionParameter(item)
	}

	return p
}

func (d *DefaultVisitor) PatternBinder(p PatternBinder) PatternBinder {
	if p == nil {
		return nil
	}

	switch b := p.(type) {
	case *IdentifierBinder:
		return d.Specific.IdentifierBinder(b)
	case *ObjectRestBinder:
		return d.Specific.ObjectRestBinder(b)
	case *ArrayRestBinder:
		return d.Specific.ArrayRestBinder(b)
	case *ObjectPropertyBinder:
		return d.Specific.ObjectPropertyBinder(b)
	case *ArrayItemBinder:
		return d.Specific.ArrayItemBinder(b)
	case *ObjectBinding:
		return d.Specific.ObjectBinding(b)
	case *ArrayBinding:
		return d.Specific.ArrayBinding(b)

	default:
		panic("Unknown pattern binder type")
	}
}

func (d *DefaultVisitor) IdentifierBinder(b *IdentifierBinder) *IdentifierBinder {
	b.Name = d.Specific.Identifier(b.Name)

	return b
}

func (d *DefaultVisitor) ObjectRestBinder(b *ObjectRestBinder) *ObjectRestBinder {
	b.Name = d.Specific.Identifier(b.Name)

	return b
}

func (d *DefaultVisitor) ArrayRestBinder(b *ArrayRestBinder) *ArrayRestBinder {
	b.Name = d.Specific.Identifier(b.Name)

	return b
}

func (d *DefaultVisitor) ObjectPropertyBinder(b *ObjectPropertyBinder) *ObjectPropertyBinder {
	b.PropertyName = d.Specific.Identifier(b.PropertyName)
	b.Property = d.Specific.PatternBinder(b.Property)
	b.DefaultValue = d.Specific.Expression(b.DefaultValue)

	return b
}

func (d *DefaultVisitor) ArrayItemBinder(b *ArrayItemBinder) *ArrayItemBinder {
	b.Item = d.Specific.PatternBinder(b.Item)
	b.DefaultValue = d.Specific.Expression(b.DefaultValue)

	return b
}
