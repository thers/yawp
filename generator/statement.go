package generator

import "yawp/parser/ast"

func (g *Generator) statement(astmt ast.Statement) {
	if astmt == nil {
		return
	}

	switch stmt := astmt.(type) {
	case *ast.WhileStatement:
		g.whileStatement(stmt)
	case *ast.ExpressionStatement:
		g.expressionStatement(stmt)
	case *ast.BlockStatement:
		g.blockStatement(stmt)
	case *ast.EmptyStatement:
		return
	case *ast.FlowInterfaceStatement, *ast.FlowTypeStatement:
		return
	case *ast.DebuggerStatement:
		g.debuggerStatement()
	case *ast.IfStatement:
		g.ifStatement(stmt)
	case *ast.ClassStatement:
		g.classStatement(stmt)
	default:
		g.str("'unknown statement';\n")
	}
}

func (g *Generator) statements(list []ast.Statement) {
	for _, astmt := range list {
		g.statement(astmt)
		g.nl()
	}
}

func (g *Generator) blockStatement(blck *ast.BlockStatement) {
	for index, stmt := range blck.List {
		if index > 0 {
			g.semicolon()
		}

		g.statement(stmt)
	}
}

func (g *Generator) expressionStatement(stmt *ast.ExpressionStatement) {
	g.expression(stmt.Expression)
}

func (g *Generator) whileStatement(stmt *ast.WhileStatement) {
	g.str("while(")
	g.expression(stmt.Test)
	g.str("){")
	g.statement(stmt.Body)
	g.rune('}')
}

func (g *Generator) debuggerStatement() {
	if !g.opt.Minify {
		g.str("debugger")
		g.semicolon()
	}
}

func (g *Generator) ifStatement(stmt *ast.IfStatement) {
	g.str("if(")
	g.expression(stmt.Test)
	g.str("){")
	g.statement(stmt.Consequent)
	g.rune('}')

	if stmt.Alternate != nil {
		g.str("else{")
		g.statement(stmt.Alternate)
		g.rune('}')
	}
}
