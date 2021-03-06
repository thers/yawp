package ast

type Visitor interface {
	Body(stmts []IStmt) []IStmt

	// statements
	Statement(stmt IStmt) IStmt
	Statements(stmt Statements) Statements
	ExportDeclaration(stmt *ExportStatement) IStmt
	ImportDeclaration(stmt *ImportStatement) IStmt
	FlowTypeStatement(stmt *FlowTypeStatement) *FlowTypeStatement
	FlowInterfaceStatement(stmt *FlowInterfaceStatement) *FlowInterfaceStatement
	BlockStatement(stmt *BlockStatement) IStmt
	ClassStatement(stmt *ClassStatement) IStmt
	ClassFieldStatement(stmt *ClassFieldStatement) IStmt
	ClassAccessorStatement(stmt *ClassAccessorStatement) IStmt
	ClassMethodStatement(stmt *ClassMethodStatement) IStmt
	LegacyDecoratorStatement(stmt *LegacyDecoratorStatement) IStmt
	ForInStatement(stmt *ForInStatement) IStmt
	ForOfStatement(stmt *ForOfStatement) IStmt
	ForStatement(stmt *ForStatement) IStmt
	BranchStatement(stmt *BranchStatement) IStmt
	CatchStatement(stmt *CatchStatement) IStmt
	DebuggerStatement(stmt *DebuggerStatement) IStmt
	DoWhileStatement(stmt *DoWhileStatement) IStmt
	EmptyStatement(stmt *EmptyStatement) IStmt
	ExpressionStatement(stmt *ExpressionStatement) IStmt
	IfStatement(stmt *IfStatement) IStmt
	LabelledStatement(stmt *LabelledStatement) IStmt
	ReturnStatement(stmt *ReturnStatement) IStmt
	SwitchStatement(stmt *SwitchStatement) IStmt
	CaseStatement(stmt *CaseStatement) IStmt
	ThrowStatement(stmt *ThrowStatement) IStmt
	TryStatement(stmt *TryStatement) IStmt
	VariableStatement(stmt *VariableStatement) IStmt
	WhileStatement(stmt *WhileStatement) IStmt
	WithStatement(stmt *WithStatement) IStmt

	// expressions
	Expression(exp IExpr) IExpr
	Expressions(exps Expressions) Expressions
	Identifier(exp *Identifier) *Identifier
	MemberExpression(exp *MemberExpression) IExpr

	ImportClause(exp *ImportClause) *ImportClause
	ImportCall(exp *ImportCall) *ImportCall

	FlowTypeAssertion(exp *FlowTypeAssertionExpression) *FlowTypeAssertionExpression

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
	SuperExpression(exp *SuperExpression) IExpr
	CoalesceExpression(exp *CoalesceExpression) *CoalesceExpression
	ConditionalExpression(exp *ConditionalExpression) *ConditionalExpression
	JsxElement(exp *JSXElement) *JSXElement
	JsxFragment(exp *JSXFragment) *JSXFragment
	NewTargetExpression(exp *NewTargetExpression) *NewTargetExpression
	AssignExpression(exp *AssignmentExpression) *AssignmentExpression
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
	PatternBinder(binder PatternBinder) PatternBinder
	IdentifierBinder(b *IdentifierBinder) *IdentifierBinder
	ObjectRestBinder(b *ObjectRestBinder) *ObjectRestBinder
	ArrayRestBinder(b *ArrayRestBinder) *ArrayRestBinder
	ObjectPropertyBinder(b *ObjectPropertyBinder) *ObjectPropertyBinder
	ArrayItemBinder(b *ArrayItemBinder) *ArrayItemBinder
	ArrayBinding(b *ArrayBinding) *ArrayBinding
	ObjectBinding(b *ObjectBinding) *ObjectBinding

	FunctionBody(fb *FunctionBody) *FunctionBody
	FunctionParameters(params *FunctionParameters) *FunctionParameters
	FunctionParameter(param FunctionParameter) FunctionParameter
	IdentifierParameter(ip *IdentifierParameter) FunctionParameter
	RestParameter(rp *RestParameter) FunctionParameter
	PatternParameter(pp *PatternParameter) FunctionParameter

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

	ReplacementStatement  IStmt
	ReplacementExpression IExpr
}

func (w *Walker) Body(stmts []IStmt) []IStmt {
	for index, stmt := range stmts {
		stmts[index] = w.Visitor.Statement(stmt)
	}

	return stmts
}

func (w *Walker) Statements(statements Statements) Statements {
	for index, statement := range statements {
		statements[index] = w.Visitor.Statement(statement)
	}

	return statements
}

func (w *Walker) Statement(stmt IStmt) IStmt {
	if stmt == nil {
		return nil
	}

	switch s := stmt.(type) {
	case Statements:
		stmt = w.Visitor.Statements(s)
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

func (w *Walker) Expressions(exps Expressions) Expressions {
	for index, exp := range exps {
		exps[index] = w.Visitor.Expression(exp)
	}

	return exps
}

func (w *Walker) Expression(exp IExpr) IExpr {
	if exp == nil {
		return nil
	}

	switch s := exp.(type) {
	case Expressions:
		exp = w.Visitor.Expressions(s)
	case *Identifier:
		exp = w.Visitor.Identifier(s)
	case *MemberExpression:
		exp = w.Visitor.MemberExpression(s)
	case *ImportClause:
		exp = w.Visitor.ImportClause(s)
	case *ImportCall:
		exp = w.Visitor.ImportCall(s)
	case *FlowTypeAssertionExpression:
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
	case *SuperExpression:
		exp = w.Visitor.SuperExpression(s)
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
	case *AssignmentExpression:
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

func (w *Walker) ExportDeclaration(stmt *ExportStatement) IStmt {
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

func (w *Walker) ImportDeclaration(stmt *ImportStatement) IStmt {
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

func (w *Walker) FlowTypeAssertion(exp *FlowTypeAssertionExpression) *FlowTypeAssertionExpression {
	return exp
}

func (w *Walker) BlockStatement(block *BlockStatement) IStmt {
	for index, stmt := range block.List {
		block.List[index] = w.Visitor.Statement(stmt)
	}

	return block
}

func (w *Walker) ClassStatement(stmt *ClassStatement) IStmt {
	stmt.Expression = w.Visitor.ClassExpression(stmt.Expression)

	return stmt
}

func (w *Walker) ClassFieldStatement(stmt *ClassFieldStatement) IStmt {
	stmt.Initializer = w.Visitor.Expression(stmt.Initializer)
	stmt.Name = w.Visitor.ObjectPropertyName(stmt.Name)

	return stmt
}

func (w *Walker) ComputedName(n *ComputedName) *ComputedName {
	n.Expression = w.Visitor.Expression(n.Expression)
	return n
}

func (w *Walker) ClassAccessorStatement(stmt *ClassAccessorStatement) IStmt {
	stmt.Body = w.Visitor.FunctionLiteral(stmt.Body)
	stmt.Field = w.Visitor.ObjectPropertyName(stmt.Field)

	return stmt
}

func (w *Walker) ClassMethodStatement(stmt *ClassMethodStatement) IStmt {
	stmt.Body = w.Visitor.FunctionBody(stmt.Body)
	stmt.Name = w.Visitor.ObjectPropertyName(stmt.Name)
	stmt.Parameters = w.Visitor.FunctionParameters(stmt.Parameters)

	return stmt
}

func (w *Walker) LegacyDecoratorStatement(stmt *LegacyDecoratorStatement) IStmt {
	for index, dec := range stmt.Decorators {
		stmt.Decorators[index] = w.Visitor.Expression(dec)
	}

	stmt.Subject = w.Visitor.Statement(stmt.Subject)

	return stmt
}

func (w *Walker) ForInStatement(stmt *ForInStatement) IStmt {
	stmt.Left = w.Visitor.Statement(stmt.Left)
	stmt.Right = w.Visitor.Expression(stmt.Right)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) ForOfStatement(stmt *ForOfStatement) IStmt {
	stmt.Left = w.Visitor.Statement(stmt.Left)
	stmt.Right = w.Visitor.Expression(stmt.Right)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) ForStatement(stmt *ForStatement) IStmt {
	stmt.Initializer = w.Visitor.Statement(stmt.Initializer)
	stmt.Test = w.Visitor.Expression(stmt.Test)
	stmt.Update = w.Visitor.Expression(stmt.Update)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) BranchStatement(stmt *BranchStatement) IStmt {
	stmt.Label = w.Visitor.Identifier(stmt.Label)

	return stmt
}

func (w *Walker) CatchStatement(stmt *CatchStatement) IStmt {
	if stmt == nil {
		return nil
	}

	stmt.Parameter = w.Visitor.PatternBinder(stmt.Parameter)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) DebuggerStatement(stmt *DebuggerStatement) IStmt {
	return stmt
}

func (w *Walker) DoWhileStatement(stmt *DoWhileStatement) IStmt {
	stmt.Body = w.Visitor.Statement(stmt.Body)
	stmt.Test = w.Visitor.Expression(stmt.Test)

	return stmt
}

func (w *Walker) EmptyStatement(stmt *EmptyStatement) IStmt {
	return stmt
}

func (w *Walker) ExpressionStatement(stmt *ExpressionStatement) IStmt {
	stmt.Expression = w.Visitor.Expression(stmt.Expression)

	return stmt
}

func (w *Walker) IfStatement(stmt *IfStatement) IStmt {
	stmt.Test = w.Visitor.Expression(stmt.Test)
	stmt.Consequent = w.Visitor.Statement(stmt.Consequent)
	stmt.Alternate = w.Visitor.Statement(stmt.Alternate)

	return stmt
}

func (w *Walker) LabelledStatement(stmt *LabelledStatement) IStmt {
	stmt.Label = w.Visitor.Identifier(stmt.Label)
	stmt.Statement = w.Visitor.Statement(stmt.Statement)

	return stmt
}

func (w *Walker) ReturnStatement(stmt *ReturnStatement) IStmt {
	stmt.Argument = w.Visitor.Expression(stmt.Argument)

	return stmt
}

func (w *Walker) SwitchStatement(stmt *SwitchStatement) IStmt {
	stmt.Discriminant = w.Visitor.Expression(stmt.Discriminant)

	for index, c := range stmt.Body {
		stmt.Body[index] = w.Visitor.Statement(c)
	}

	return stmt
}

func (w *Walker) CaseStatement(s *CaseStatement) IStmt {
	s.Test = w.Visitor.Expression(s.Test)

	for index, st := range s.Consequent {
		s.Consequent[index] = w.Visitor.Statement(st)
	}

	return s
}

func (w *Walker) ThrowStatement(stmt *ThrowStatement) IStmt {
	stmt.Argument = w.Visitor.Expression(stmt.Argument)

	return stmt
}

func (w *Walker) TryStatement(stmt *TryStatement) IStmt {
	stmt.Body = w.Visitor.Statement(stmt.Body)
	stmt.Catch = w.Visitor.Statement(stmt.Catch)
	stmt.Finally = w.Visitor.Statement(stmt.Finally)

	return stmt
}

func (w *Walker) VariableStatement(stmt *VariableStatement) IStmt {
	for index, v := range stmt.List {
		stmt.List[index] = w.Visitor.VariableBinding(v)
	}

	return stmt
}

func (w *Walker) WhileStatement(stmt *WhileStatement) IStmt {
	stmt.Test = w.Visitor.Expression(stmt.Test)
	stmt.Body = w.Visitor.Statement(stmt.Body)

	return stmt
}

func (w *Walker) WithStatement(stmt *WithStatement) IStmt {
	stmt.Object = w.Visitor.Expression(stmt.Object)
	stmt.Body = w.Visitor.Statement(stmt.Body)

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
	stmt.Body = w.Visitor.FunctionBody(stmt.Body)

	return stmt
}

func (w *Walker) FunctionBody(fb *FunctionBody) *FunctionBody {
	for index, statement := range fb.List {
		fb.List[index] = w.Visitor.Statement(statement)
	}

	return fb
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

func (w *Walker) MemberExpression(exp *MemberExpression) IExpr {
	return exp
}

func (w *Walker) SuperExpression(exp *SuperExpression) IExpr {
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

func (w *Walker) AssignExpression(exp *AssignmentExpression) *AssignmentExpression {
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

	exp.Body = w.Visitor.FunctionBody(exp.Body)

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
	case *PatternParameter:
		return w.Visitor.PatternParameter(p)

	default:
		panic("Unknown function parameter type")
	}
}

func (w *Walker) IdentifierParameter(ip *IdentifierParameter) FunctionParameter {
	ip.Id = w.Visitor.Identifier(ip.Id)
	ip.DefaultValue = w.Visitor.Expression(ip.DefaultValue)

	return ip
}

func (w *Walker) PatternParameter(pp *PatternParameter) FunctionParameter {
	pp.Binder = w.Visitor.PatternBinder(pp.Binder)
	pp.DefaultValue = w.Visitor.Expression(pp.DefaultValue)

	return pp
}

func (w *Walker) RestParameter(rp *RestParameter) FunctionParameter {
	rp.Binder = w.Visitor.PatternBinder(rp.Binder)

	return rp
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
	b.PropertyName = w.Visitor.ObjectPropertyName(b.PropertyName)
	b.Binder = w.Visitor.PatternBinder(b.Binder)
	b.DefaultValue = w.Visitor.Expression(b.DefaultValue)

	return b
}

func (w *Walker) ArrayItemBinder(b *ArrayItemBinder) *ArrayItemBinder {
	b.Binder = w.Visitor.PatternBinder(b.Binder)
	b.DefaultValue = w.Visitor.Expression(b.DefaultValue)

	return b
}
