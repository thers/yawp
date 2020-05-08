package ast

type Visitor interface {
	Body(stmts []Statement) []Statement

	// statements
	Statement(stmt Statement) Statement
	ExportDeclaration(stmt *ExportStatement) Statement
	ImportDeclaration(stmt *ImportStatement) Statement
	FlowTypeStatement(stmt *FlowTypeStatement) *FlowTypeStatement
	FlowInterfaceStatement(stmt *FlowInterfaceStatement) *FlowInterfaceStatement
	BlockStatement(stmt *BlockStatement) Statement
	ClassStatement(stmt *ClassStatement) Statement
	ClassFieldStatement(stmt *ClassFieldStatement) Statement
	ClassAccessorStatement(stmt *ClassAccessorStatement) Statement
	ClassMethodStatement(stmt *ClassMethodStatement) Statement
	LegacyDecoratorStatement(stmt *LegacyDecoratorStatement) Statement
	ForInStatement(stmt *ForInStatement) Statement
	ForOfStatement(stmt *ForOfStatement) Statement
	ForStatement(stmt *ForStatement) Statement
	BranchStatement(stmt *BranchStatement) Statement
	CatchStatement(stmt *CatchStatement) Statement
	DebuggerStatement(stmt *DebuggerStatement) Statement
	DoWhileStatement(stmt *DoWhileStatement) Statement
	EmptyStatement(stmt *EmptyStatement) Statement
	ExpressionStatement(stmt *ExpressionStatement) Statement
	IfStatement(stmt *IfStatement) Statement
	LabelledStatement(stmt *LabelledStatement) Statement
	ReturnStatement(stmt *ReturnStatement) Statement
	SwitchStatement(stmt *SwitchStatement) Statement
	CaseStatement(stmt *CaseStatement) Statement
	ThrowStatement(stmt *ThrowStatement) Statement
	TryStatement(stmt *TryStatement) Statement
	VariableStatement(stmt *VariableStatement) Statement
	WhileStatement(stmt *WhileStatement) Statement
	WithStatement(stmt *WithStatement) Statement
	YieldStatement(stmt *YieldStatement) Statement

	// expressions
	Expression(exp Expression) Expression
	Identifier(exp *Identifier) *Identifier
	MemberExpression(exp *MemberExpression) Expression

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
	CallExpression(exp *CallExpression) *CallExpression
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
	ExportClause(c ExportClause) ExportClause
	ExportNamespaceFromClause(c *ExportNamespaceFromClause) *ExportNamespaceFromClause
	ExportNamedFromClause(c *ExportNamedFromClause) *ExportNamedFromClause
	ExportNamedClause(c *ExportNamedClause) *ExportNamedClause
	ExportVarClause(c *ExportVarClause) *ExportVarClause
	ExportFunctionClause(c *ExportFunctionClause) *ExportFunctionClause
	ExportClassClause(c *ExportClassClause) *ExportClassClause
	ExportDefaultClause(c *ExportDefaultClause) *ExportDefaultClause

	Js(js *Js) *Js
}

type Walker struct {
	Visitor Visitor

	ReplacementStatement  Statement
	ReplacementExpression Expression
}

func (w *Walker) Body(stmts []Statement) []Statement {
	for index, stmt := range stmts {
		stmts[index] = w.Visitor.Statement(stmt)
	}

	return stmts
}

func (w *Walker) Statement(stmt Statement) Statement {
	if stmt == nil {
		return nil
	}

	switch s := stmt.(type) {
	case *ExportStatement:
		stmt = w.Visitor.ExportDeclaration(s)
	case *ImportStatement:
		stmt = w.Visitor.ImportDeclaration(s)
	case *FlowTypeStatement:
		stmt = w.Visitor.FlowTypeStatement(s)
	case *FlowInterfaceStatement:
		stmt = w.Visitor.FlowInterfaceStatement(s)
	case *BlockStatement:
		stmt = w.Visitor.BlockStatement(s)
	case *ClassStatement:
		stmt = w.Visitor.ClassStatement(s)
	case *ClassFieldStatement:
		stmt = w.Visitor.ClassFieldStatement(s)
	case *ClassAccessorStatement:
		stmt = w.Visitor.ClassAccessorStatement(s)
	case *ClassMethodStatement:
		stmt = w.Visitor.ClassMethodStatement(s)
	case *LegacyDecoratorStatement:
		stmt = w.Visitor.LegacyDecoratorStatement(s)
	case *ForInStatement:
		stmt = w.Visitor.ForInStatement(s)
	case *ForOfStatement:
		stmt = w.Visitor.ForOfStatement(s)
	case *ForStatement:
		stmt = w.Visitor.ForStatement(s)
	case *BranchStatement:
		stmt = w.Visitor.BranchStatement(s)
	case *CatchStatement:
		stmt = w.Visitor.CatchStatement(s)
	case *DebuggerStatement:
		stmt = w.Visitor.DebuggerStatement(s)
	case *DoWhileStatement:
		stmt = w.Visitor.DoWhileStatement(s)
	case *EmptyStatement:
		stmt = w.Visitor.EmptyStatement(s)
	case *ExpressionStatement:
		stmt = w.Visitor.ExpressionStatement(s)
	case *IfStatement:
		stmt = w.Visitor.IfStatement(s)
	case *LabelledStatement:
		stmt = w.Visitor.LabelledStatement(s)
	case *ReturnStatement:
		stmt = w.Visitor.ReturnStatement(s)
	case *SwitchStatement:
		stmt = w.Visitor.SwitchStatement(s)
	case *ThrowStatement:
		stmt = w.Visitor.ThrowStatement(s)
	case *TryStatement:
		stmt = w.Visitor.TryStatement(s)
	case *VariableStatement:
		stmt = w.Visitor.VariableStatement(s)
	case *WhileStatement:
		stmt = w.Visitor.WhileStatement(s)
	case *WithStatement:
		stmt = w.Visitor.WithStatement(s)
	case *YieldStatement:
		stmt = w.Visitor.YieldStatement(s)
	case *FunctionLiteral:
		stmt = w.Visitor.FunctionLiteral(s)
	case *Js:
		stmt = w.Visitor.Js(s)

	default:
		panic("Unknown statement")
	}

	if w.ReplacementStatement != nil {
		defer func() {
			w.ReplacementStatement = nil
		}()

		return w.ReplacementStatement
	}

	return stmt
}

func (w *Walker) Expression(exp Expression) Expression {
	if exp == nil {
		return nil
	}

	switch s := exp.(type) {
	case *Identifier:
		exp = w.Visitor.Identifier(s)
	case *MemberExpression:
		exp = w.Visitor.MemberExpression(s)
	case *ImportClause:
		exp = w.Visitor.ImportClause(s)
	case *ImportCall:
		exp = w.Visitor.ImportCall(s)
	case *FlowTypeAssertion:
		exp = w.Visitor.FlowTypeAssertion(s)
	case *FunctionLiteral:
		exp = w.Visitor.FunctionLiteral(s)
	case *StringLiteral:
		exp = w.Visitor.StringLiteral(s)
	case *BooleanLiteral:
		exp = w.Visitor.BooleanLiteral(s)
	case *ObjectLiteral:
		exp = w.Visitor.ObjectLiteral(s)
	case *ArrayLiteral:
		exp = w.Visitor.ArrayLiteral(s)
	case *NullLiteral:
		exp = w.Visitor.NullLiteral(s)
	case *NumberLiteral:
		exp = w.Visitor.NumberLiteral(s)
	case *RegExpLiteral:
		exp = w.Visitor.RegExpLiteral(s)
	case *ObjectSpread:
		exp = w.Visitor.ObjectSpread(s)
	case *ArraySpread:
		exp = w.Visitor.ArraySpread(s)
	case *ArrayBinding:
		exp = w.Visitor.ArrayBinding(s)
	case *ObjectBinding:
		exp = w.Visitor.ObjectBinding(s)
	case *VariableBinding:
		exp = w.Visitor.VariableBinding(s)
	case *ClassExpression:
		exp = w.Visitor.ClassExpression(s)
	case *ClassSuperExpression:
		exp = w.Visitor.ClassSuperExpression(s)
	case *CoalesceExpression:
		exp = w.Visitor.CoalesceExpression(s)
	case *ConditionalExpression:
		exp = w.Visitor.ConditionalExpression(s)
	case *JSXElement:
		exp = w.Visitor.JsxElement(s)
	case *JSXFragment:
		exp = w.Visitor.JsxFragment(s)
	case *NewTargetExpression:
		exp = w.Visitor.NewTargetExpression(s)
	case *AssignExpression:
		exp = w.Visitor.AssignExpression(s)
	case *BinaryExpression:
		exp = w.Visitor.BinaryExpression(s)
	case *CallExpression:
		exp = w.Visitor.CallExpression(s)
	case *SpreadExpression:
		exp = w.Visitor.SpreadExpression(s)
	case *NewExpression:
		exp = w.Visitor.NewExpression(s)
	case *SequenceExpression:
		exp = w.Visitor.SequenceExpression(s)
	case *ThisExpression:
		exp = w.Visitor.ThisExpression(s)
	case *UnaryExpression:
		exp = w.Visitor.UnaryExpression(s)
	case *ArrowFunctionExpression:
		exp = w.Visitor.ArrowFunctionExpression(s)
	case *AwaitExpression:
		exp = w.Visitor.AwaitExpression(s)
	case *OptionalObjectMemberAccessExpression:
		exp = w.Visitor.OptionalObjectMemberAccessExpression(s)
	case *OptionalArrayMemberAccessExpression:
		exp = w.Visitor.OptionalArrayMemberAccessExpression(s)
	case *OptionalCallExpression:
		exp = w.Visitor.OptionalCallExpression(s)
	case *TemplateExpression:
		exp = w.Visitor.TemplateExpression(s)
	case *TaggedTemplateExpression:
		exp = w.Visitor.TaggedTemplateExpression(s)
	case *YieldExpression:
		exp = w.Visitor.YieldExpression(s)
	case *Js:
		exp = w.Visitor.Js(s)

	default:
		panic("Unknown expression")
	}

	if w.ReplacementExpression != nil {
		defer func() {
			w.ReplacementExpression = nil
		}()

		return w.ReplacementExpression
	}

	return exp
}

func (w *Walker) Js(js *Js) *Js {
	return js
}

func (w *Walker) ExportDeclaration(stmt *ExportStatement) Statement {
	stmt.Clause = w.Visitor.ExportClause(stmt.Clause)

	return stmt
}

func (w *Walker) ExportClause(ac ExportClause) ExportClause {
	if ac == nil {
		return nil
	}

	switch c := ac.(type) {
	case *ExportNamespaceFromClause:
		return w.Visitor.ExportNamespaceFromClause(c)
	case *ExportNamedFromClause:
		return w.Visitor.ExportNamedFromClause(c)
	case *ExportNamedClause:
		return w.Visitor.ExportNamedClause(c)
	case *ExportVarClause:
		return w.Visitor.ExportVarClause(c)
	case *ExportFunctionClause:
		return w.Visitor.ExportFunctionClause(c)
	case *ExportClassClause:
		return w.Visitor.ExportClassClause(c)
	case *ExportDefaultClause:
		return w.Visitor.ExportDefaultClause(c)
	case *FlowTypeStatement:
		return w.Visitor.FlowTypeStatement(c)
	case *FlowInterfaceStatement:
		return w.Visitor.FlowInterfaceStatement(c)

	default:
		panic("Unknown export clause type")
	}
}

func (w *Walker) ExportNamespaceFromClause(c *ExportNamespaceFromClause) *ExportNamespaceFromClause {
	panic("implement me")
}

func (w *Walker) ExportNamedFromClause(c *ExportNamedFromClause) *ExportNamedFromClause {
	panic("implement me")
}

func (w *Walker) ExportNamedClause(c *ExportNamedClause) *ExportNamedClause {
	panic("implement me")
}

func (w *Walker) ExportVarClause(c *ExportVarClause) *ExportVarClause {
	panic("implement me")
}

func (w *Walker) ExportFunctionClause(c *ExportFunctionClause) *ExportFunctionClause {
	panic("implement me")
}

func (w *Walker) ExportClassClause(c *ExportClassClause) *ExportClassClause {
	panic("implement me")
}

func (w *Walker) ExportDefaultClause(c *ExportDefaultClause) *ExportDefaultClause {
	panic("implement me")
}

func (w *Walker) ImportDeclaration(stmt *ImportStatement) Statement {
	for index, i := range stmt.Imports {
		stmt.Imports[index] = w.Visitor.ImportClause(i)
	}

	return stmt
}

func (w *Walker) ImportClause(exp *ImportClause) *ImportClause {
	exp.LocalIdentifier = w.Visitor.Identifier(exp.LocalIdentifier)
	exp.ModuleIdentifier = w.Visitor.Identifier(exp.ModuleIdentifier)

	return exp
}

func (w *Walker) ImportCall(exp *ImportCall) *ImportCall {
	exp.Expression = w.Visitor.Expression(exp.Expression)

	return exp
}

func (w *Walker) FlowTypeStatement(stmt *FlowTypeStatement) *FlowTypeStatement {
	return stmt
}

func (w *Walker) FlowInterfaceStatement(stmt *FlowInterfaceStatement) *FlowInterfaceStatement {
	return stmt
}

func (w *Walker) FlowTypeAssertion(exp *FlowTypeAssertion) *FlowTypeAssertion {
	return exp
}

func (w *Walker) BlockStatement(block *BlockStatement) Statement {
	for index, stmt := range block.List {
		block.List[index] = w.Visitor.Statement(stmt)
	}

	return block
}

func (w *Walker) ClassStatement(stmt *ClassStatement) Statement {
	stmt.Expression = w.Visitor.ClassExpression(stmt.Expression)

	return stmt
}

func (w *Walker) ClassFieldStatement(stmt *ClassFieldStatement) Statement {
	stmt.Initializer = w.Visitor.Expression(stmt.Initializer)
	stmt.Name = w.Visitor.ClassFieldName(stmt.Name)

	return stmt
}

func (w *Walker) ClassFieldName(name ClassFieldName) ClassFieldName {
	switch n := name.(type) {
	case *ComputedName:
		return w.Visitor.ComputedName(n)
	case *Identifier:
		return w.Visitor.Identifier(n)

	default:
		return n
	}
}

func (w *Walker) ComputedName(n *ComputedName) *ComputedName {
	n.Expression = w.Visitor.Expression(n.Expression)
	return n
}

func (w *Walker) ClassAccessorStatement(stmt *ClassAccessorStatement) Statement {
	stmt.Body = w.Visitor.FunctionLiteral(stmt.Body)
	stmt.Field = w.Visitor.ClassFieldName(stmt.Field)

	return stmt
}

func (w *Walker) ClassMethodStatement(stmt *ClassMethodStatement) Statement {
	stmt.Body = w.Visitor.Statement(stmt.Body)
	stmt.Name = w.Visitor.ClassFieldName(stmt.Name)
	stmt.Parameters = w.Visitor.FunctionParameters(stmt.Parameters)

	return stmt
}

func (w *Walker) LegacyDecoratorStatement(stmt *LegacyDecoratorStatement) Statement {
	for index, dec := range stmt.Decorators {
		stmt.Decorators[index] = w.Visitor.Expression(dec)
	}

	stmt.Subject = w.Visitor.Statement(stmt.Subject)

	return stmt
}

func (w *Walker) ForInStatement(stmt *ForInStatement) Statement {
	stmt.Left = w.Visitor.Statement(stmt.Left)
	stmt.Right = w.Visitor.Expression(stmt.Right)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) ForOfStatement(stmt *ForOfStatement) Statement {
	stmt.Left = w.Visitor.Statement(stmt.Left)
	stmt.Right = w.Visitor.Expression(stmt.Right)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) ForStatement(stmt *ForStatement) Statement {
	stmt.Initializer = w.Visitor.Statement(stmt.Initializer)
	stmt.Test = w.Visitor.Expression(stmt.Test)
	stmt.Update = w.Visitor.Expression(stmt.Update)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) BranchStatement(stmt *BranchStatement) Statement {
	stmt.Label = w.Visitor.Identifier(stmt.Label)

	return stmt
}

func (w *Walker) CatchStatement(stmt *CatchStatement) Statement {
	if stmt == nil {
		return nil
	}

	stmt.Parameter = w.Visitor.Identifier(stmt.Parameter)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) DebuggerStatement(stmt *DebuggerStatement) Statement {
	return stmt
}

func (w *Walker) DoWhileStatement(stmt *DoWhileStatement) Statement {
	stmt.Body = w.Visitor.Statement(stmt.Body)
	stmt.Test = w.Visitor.Expression(stmt.Test)

	return stmt
}

func (w *Walker) EmptyStatement(stmt *EmptyStatement) Statement {
	return stmt
}

func (w *Walker) ExpressionStatement(stmt *ExpressionStatement) Statement {
	stmt.Expression = w.Visitor.Expression(stmt.Expression)

	return stmt
}

func (w *Walker) IfStatement(stmt *IfStatement) Statement {
	stmt.Test = w.Visitor.Expression(stmt.Test)
	stmt.Consequent = w.Visitor.Statement(stmt.Consequent)
	stmt.Alternate = w.Visitor.Statement(stmt.Alternate)

	return stmt
}

func (w *Walker) LabelledStatement(stmt *LabelledStatement) Statement {
	stmt.Label = w.Visitor.Identifier(stmt.Label)
	stmt.Statement = w.Visitor.Statement(stmt.Statement)

	return stmt
}

func (w *Walker) ReturnStatement(stmt *ReturnStatement) Statement {
	stmt.Argument = w.Visitor.Expression(stmt.Argument)

	return stmt
}

func (w *Walker) SwitchStatement(stmt *SwitchStatement) Statement {
	stmt.Discriminant = w.Visitor.Expression(stmt.Discriminant)

	for index, c := range stmt.Body {
		stmt.Body[index] = w.Visitor.Statement(c)
	}

	return stmt
}

func (w *Walker) CaseStatement(s *CaseStatement) Statement {
	s.Test = w.Visitor.Expression(s.Test)

	for index, st := range s.Consequent {
		s.Consequent[index] = w.Visitor.Statement(st)
	}

	return s
}

func (w *Walker) ThrowStatement(stmt *ThrowStatement) Statement {
	stmt.Argument = w.Visitor.Expression(stmt.Argument)

	return stmt
}

func (w *Walker) TryStatement(stmt *TryStatement) Statement {
	stmt.Body = w.Visitor.Statement(stmt.Body)
	stmt.Catch = w.Visitor.Statement(stmt.Catch)
	stmt.Finally = w.Visitor.Statement(stmt.Finally)

	return stmt
}

func (w *Walker) VariableStatement(stmt *VariableStatement) Statement {
	for index, v := range stmt.List {
		stmt.List[index] = w.Visitor.VariableBinding(v)
	}

	return stmt
}

func (w *Walker) WhileStatement(stmt *WhileStatement) Statement {
	stmt.Test = w.Visitor.Expression(stmt.Test)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) WithStatement(stmt *WithStatement) Statement {
	stmt.Object = w.Visitor.Expression(stmt.Object)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) YieldStatement(stmt *YieldStatement) Statement {
	stmt.Expression = w.Visitor.YieldExpression(stmt.Expression)

	return stmt
}

func (w *Walker) Identifier(exp *Identifier) *Identifier {
	if exp == nil {
		return nil
	}

	return exp
}

func (w *Walker) FunctionLiteral(stmt *FunctionLiteral) *FunctionLiteral {
	stmt.Id = w.Visitor.Identifier(stmt.Id)
	stmt.Parameters = w.Visitor.FunctionParameters(stmt.Parameters)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) StringLiteral(exp *StringLiteral) *StringLiteral {
	return exp
}

func (w *Walker) BooleanLiteral(exp *BooleanLiteral) *BooleanLiteral {
	return exp
}

func (w *Walker) ObjectLiteral(exp *ObjectLiteral) *ObjectLiteral {
	for index, prop := range exp.Properties {
		exp.Properties[index] = w.Visitor.ObjectProperty(prop)
	}

	return exp
}

func (w *Walker) ObjectProperty(ap ObjectProperty) ObjectProperty {
	if ap == nil {
		return nil
	}

	switch p := ap.(type) {
	case *ObjectPropertySetter:
		return w.Visitor.ObjectPropertySetter(p)
	case *ObjectPropertyGetter:
		return w.Visitor.ObjectPropertyGetter(p)
	case *ObjectPropertyValue:
		return w.Visitor.ObjectPropertyValue(p)
	case *ObjectSpread:
		return w.Visitor.ObjectSpread(p)

	default:
		panic("Unknown object property type")
	}
}

func (w *Walker) ObjectPropertySetter(p *ObjectPropertySetter) *ObjectPropertySetter {
	p.PropertyName = w.Visitor.ObjectPropertyName(p.PropertyName)
	p.Setter = w.Visitor.FunctionLiteral(p.Setter)

	return p
}

func (w *Walker) ObjectPropertyGetter(p *ObjectPropertyGetter) *ObjectPropertyGetter {
	p.PropertyName = w.Visitor.ObjectPropertyName(p.PropertyName)
	p.Getter = w.Visitor.FunctionLiteral(p.Getter)

	return p
}

func (w *Walker) ObjectPropertyValue(p *ObjectPropertyValue) *ObjectPropertyValue {
	p.PropertyName = w.Visitor.ObjectPropertyName(p.PropertyName)
	p.Value = w.Visitor.Expression(p.Value)

	return p
}

func (w *Walker) ObjectPropertyName(an ObjectPropertyName) ObjectPropertyName {
	switch n := an.(type) {
	case *ComputedName:
		return w.Visitor.ComputedName(n)
	case *Identifier:
		return w.Visitor.Identifier(n)

	default:
		panic("Unknown object property name type")
	}
}

func (w *Walker) ArrayLiteral(exp *ArrayLiteral) *ArrayLiteral {
	for index, item := range exp.List {
		exp.List[index] = w.Visitor.Expression(item)
	}

	return exp
}

func (w *Walker) NullLiteral(exp *NullLiteral) *NullLiteral {
	return exp
}

func (w *Walker) NumberLiteral(exp *NumberLiteral) *NumberLiteral {
	return exp
}

func (w *Walker) RegExpLiteral(exp *RegExpLiteral) *RegExpLiteral {
	return exp
}

func (w *Walker) ObjectSpread(exp *ObjectSpread) *ObjectSpread {
	exp.Expression = w.Visitor.Expression(exp.Expression)

	return exp
}

func (w *Walker) ArraySpread(exp *ArraySpread) *ArraySpread {
	exp.Expression = w.Visitor.Expression(exp.Expression)

	return exp
}

func (w *Walker) ArrayBinding(exp *ArrayBinding) *ArrayBinding {
	for index, b := range exp.List {
		exp.List[index] = w.Visitor.PatternBinder(b)
	}

	return exp
}

func (w *Walker) ObjectBinding(exp *ObjectBinding) *ObjectBinding {
	for index, b := range exp.List {
		exp.List[index] = w.Visitor.PatternBinder(b)
	}

	return exp
}

func (w *Walker) VariableBinding(exp *VariableBinding) *VariableBinding {
	exp.Binder = w.Visitor.PatternBinder(exp.Binder)
	exp.Initializer = w.Visitor.Expression(exp.Initializer)

	return exp
}

func (w *Walker) ClassExpression(exp *ClassExpression) *ClassExpression {
	exp.Name = w.Visitor.Identifier(exp.Name)
	exp.SuperClass = w.Visitor.Expression(exp.SuperClass)
	exp.Body = w.Visitor.Statement(exp.Body)

	return exp
}

func (w *Walker) MemberExpression(exp *MemberExpression) Expression {
	return exp
}

func (w *Walker) ClassSuperExpression(exp *ClassSuperExpression) *ClassSuperExpression {
	for index, arg := range exp.Arguments {
		exp.Arguments[index] = w.Visitor.Expression(arg)
	}

	return exp
}

func (w *Walker) CoalesceExpression(exp *CoalesceExpression) *CoalesceExpression {
	exp.Head = w.Visitor.Expression(exp.Head)
	exp.Consequent = w.Visitor.Expression(exp.Consequent)

	return exp
}

func (w *Walker) ConditionalExpression(exp *ConditionalExpression) *ConditionalExpression {
	exp.Test = w.Visitor.Expression(exp.Test)
	exp.Consequent = w.Visitor.Expression(exp.Consequent)
	exp.Alternate = w.Visitor.Expression(exp.Alternate)

	return exp
}

func (w *Walker) JsxElement(exp *JSXElement) *JSXElement {
	return exp
}

func (w *Walker) JsxFragment(exp *JSXFragment) *JSXFragment {
	return exp
}

func (w *Walker) NewTargetExpression(exp *NewTargetExpression) *NewTargetExpression {
	return exp
}

func (w *Walker) AssignExpression(exp *AssignExpression) *AssignExpression {
	exp.Left = w.Visitor.Expression(exp.Left)
	exp.Right = w.Visitor.Expression(exp.Right)

	return exp
}

func (w *Walker) BinaryExpression(exp *BinaryExpression) *BinaryExpression {
	exp.Left = w.Visitor.Expression(exp.Left)
	exp.Right = w.Visitor.Expression(exp.Right)

	return exp
}

func (w *Walker) CallExpression(exp *CallExpression) *CallExpression {
	exp.Callee = w.Visitor.Expression(exp.Callee)

	for index, arg := range exp.ArgumentList {
		exp.ArgumentList[index] = w.Visitor.Expression(arg)
	}

	return exp
}

func (w *Walker) SpreadExpression(exp *SpreadExpression) *SpreadExpression {
	exp.Value = w.Visitor.Expression(exp.Value)

	return exp
}

func (w *Walker) NewExpression(exp *NewExpression) *NewExpression {
	exp.Callee = w.Visitor.Expression(exp.Callee)

	for index, arg := range exp.ArgumentList {
		exp.ArgumentList[index] = w.Visitor.Expression(arg)
	}

	return exp
}

func (w *Walker) SequenceExpression(exp *SequenceExpression) *SequenceExpression {
	for index, e := range exp.Sequence {
		exp.Sequence[index] = w.Visitor.Expression(e)
	}

	return exp
}

func (w *Walker) ThisExpression(exp *ThisExpression) *ThisExpression {
	return exp
}

func (w *Walker) UnaryExpression(exp *UnaryExpression) *UnaryExpression {
	exp.Operand = w.Visitor.Expression(exp.Operand)

	return exp
}

func (w *Walker) ArrowFunctionExpression(exp *ArrowFunctionExpression) *ArrowFunctionExpression {
	if exp == nil {
		return nil
	}

	exp.Body = w.Visitor.Statement(exp.Body)

	for index, param := range exp.Parameters {
		exp.Parameters[index] = w.Visitor.FunctionParameter(param)
	}

	return exp
}

func (w *Walker) AwaitExpression(exp *AwaitExpression) *AwaitExpression {
	exp.Expression = w.Visitor.Expression(exp.Expression)

	return exp
}

func (w *Walker) OptionalObjectMemberAccessExpression(exp *OptionalObjectMemberAccessExpression) *OptionalObjectMemberAccessExpression {
	exp.Left = w.Visitor.Expression(exp.Left)
	exp.Identifier = w.Visitor.Identifier(exp.Identifier)

	return exp
}

func (w *Walker) OptionalArrayMemberAccessExpression(exp *OptionalArrayMemberAccessExpression) *OptionalArrayMemberAccessExpression {
	exp.Left = w.Visitor.Expression(exp.Left)
	exp.Index = w.Visitor.Expression(exp.Index)

	return exp
}

func (w *Walker) OptionalCallExpression(exp *OptionalCallExpression) *OptionalCallExpression {
	exp.Left = w.Visitor.Expression(exp.Left)

	for index, arg := range exp.Arguments {
		exp.Arguments[index] = w.Visitor.Expression(arg)
	}

	return exp
}

func (w *Walker) TemplateExpression(exp *TemplateExpression) *TemplateExpression {
	for index, s := range exp.Substitutions {
		exp.Substitutions[index] = w.Visitor.Expression(s)
	}

	return exp
}

func (w *Walker) TaggedTemplateExpression(exp *TaggedTemplateExpression) *TaggedTemplateExpression {
	exp.Tag = w.Visitor.Expression(exp.Tag)
	exp.Template = w.Visitor.TemplateExpression(exp.Template)

	return exp
}

func (w *Walker) YieldExpression(exp *YieldExpression) *YieldExpression {
	exp.Argument = w.Visitor.Expression(exp.Argument)

	return exp
}

func (w *Walker) FunctionParameters(params *FunctionParameters) *FunctionParameters {
	for index, param := range params.List {
		params.List[index] = w.Visitor.FunctionParameter(param)
	}

	return params
}

func (w *Walker) FunctionParameter(param FunctionParameter) FunctionParameter {
	switch p := param.(type) {
	case *IdentifierParameter:
		return w.Visitor.IdentifierParameter(p)
	case *RestParameter:
		return w.Visitor.RestParameter(p)
	case *ObjectPatternParameter:
		return w.Visitor.ObjectPatternParameter(p)
	case *ArrayPatternParameter:
		return w.Visitor.ArrayPatternParameter(p)

	default:
		panic("Unknown function parameter type")
	}
}

func (w *Walker) IdentifierParameter(p *IdentifierParameter) *IdentifierParameter {
	p.Id = w.Visitor.Identifier(p.Id)
	p.DefaultValue = w.Visitor.Expression(p.DefaultValue)

	return p
}

func (w *Walker) RestParameter(p *RestParameter) *RestParameter {
	p.Binder = w.Visitor.PatternBinder(p.Binder)

	return p
}

func (w *Walker) ObjectPatternParameter(p *ObjectPatternParameter) *ObjectPatternParameter {
	p.DefaultValue = w.Visitor.Expression(p.DefaultValue)

	for index, prop := range p.List {
		p.List[index] = w.Visitor.ObjectPatternIdentifierParameter(prop)
	}

	return p
}

func (w *Walker) ObjectPatternIdentifierParameter(p *ObjectPatternIdentifierParameter) *ObjectPatternIdentifierParameter {
	p.Parameter = w.Visitor.FunctionParameter(p.Parameter)

	return p
}

func (w *Walker) ArrayPatternParameter(p *ArrayPatternParameter) *ArrayPatternParameter {
	p.DefaultValue = w.Visitor.Expression(p.DefaultValue)

	for index, item := range p.List {
		p.List[index] = w.Visitor.FunctionParameter(item)
	}

	return p
}

func (w *Walker) PatternBinder(p PatternBinder) PatternBinder {
	if p == nil {
		return nil
	}

	switch b := p.(type) {
	case *IdentifierBinder:
		return w.Visitor.IdentifierBinder(b)
	case *ObjectRestBinder:
		return w.Visitor.ObjectRestBinder(b)
	case *ArrayRestBinder:
		return w.Visitor.ArrayRestBinder(b)
	case *ObjectPropertyBinder:
		return w.Visitor.ObjectPropertyBinder(b)
	case *ArrayItemBinder:
		return w.Visitor.ArrayItemBinder(b)
	case *ObjectBinding:
		return w.Visitor.ObjectBinding(b)
	case *ArrayBinding:
		return w.Visitor.ArrayBinding(b)

	default:
		panic("Unknown pattern binder type")
	}
}

func (w *Walker) IdentifierBinder(b *IdentifierBinder) *IdentifierBinder {
	b.Id = w.Visitor.Identifier(b.Id)

	return b
}

func (w *Walker) ObjectRestBinder(b *ObjectRestBinder) *ObjectRestBinder {
	b.Binder = w.Visitor.PatternBinder(b.Binder)

	return b
}

func (w *Walker) ArrayRestBinder(b *ArrayRestBinder) *ArrayRestBinder {
	b.Binder = w.Visitor.PatternBinder(b.Binder)

	return b
}

func (w *Walker) ObjectPropertyBinder(b *ObjectPropertyBinder) *ObjectPropertyBinder {
	b.Id = w.Visitor.ObjectPropertyName(b.Id)
	b.Binder = w.Visitor.PatternBinder(b.Binder)
	b.DefaultValue = w.Visitor.Expression(b.DefaultValue)

	return b
}

func (w *Walker) ArrayItemBinder(b *ArrayItemBinder) *ArrayItemBinder {
	b.Binder = w.Visitor.PatternBinder(b.Binder)
	b.DefaultValue = w.Visitor.Expression(b.DefaultValue)

	return b
}
