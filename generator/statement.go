package generator

import "yawp/parser/ast"


func (g *Generator) Statement(s ast.Statement) ast.Statement {
	s = g.DefaultVisitor.Statement(s)

	return s
}

func (g *Generator) BlockStatement(bs *ast.BlockStatement) *ast.BlockStatement {
	g.rune('{').indentInc().nl()

	for index, stmt := range bs.List {
		if index > 0 {
			g.semicolon()
		}

		g.Statement(stmt)
	}

	g.rune('}')

	return bs
}

func (g *Generator) ExpressionStatement(stmt *ast.ExpressionStatement) *ast.ExpressionStatement {
	g.Expression(stmt.Expression)

	return stmt
}

func (g *Generator) WhileStatement(stmt *ast.WhileStatement) *ast.WhileStatement {
	g.str("while(")
	g.Expression(stmt.Test)
	g.str("){")
	g.Statement(stmt.Body)
	g.rune('}')

	return stmt
}

func (g *Generator) DebuggerStatement(ds *ast.DebuggerStatement) *ast.DebuggerStatement {
	if !g.options.Minify {
		g.str("debugger")
		g.semicolon()
	}

	return ds
}

func (g *Generator) IfStatement(stmt *ast.IfStatement) *ast.IfStatement {
	g.str("if(")
	g.Expression(stmt.Test)
	g.str("){")
	g.Statement(stmt.Consequent)
	g.rune('}')

	if stmt.Alternate != nil {
		g.str("else{")
		g.Statement(stmt.Alternate)
		g.rune('}')
	}

	return stmt
}