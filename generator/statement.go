package generator

import "yawp/parser/ast"


func (g *Generator) Statement(s ast.Statement) ast.Statement {
	s = g.DefaultVisitor.Statement(s)

	g.nl()

	return s
}

func (g *Generator) BlockStatement(bs *ast.BlockStatement) *ast.BlockStatement {
	g.rune('{').indentInc().nl()

	for index, stmt := range bs.List {
		if index > 0 {
			g.nl()
		}

		g.Statement(stmt)
	}

	g.indentDec().nl().rune('}')

	return bs
}

func (g *Generator) ExpressionStatement(stmt *ast.ExpressionStatement) *ast.ExpressionStatement {
	g.Expression(stmt.Expression)

	return stmt
}

func (g *Generator) WhileStatement(stmt *ast.WhileStatement) *ast.WhileStatement {
	g.str("while(")
	g.Expression(stmt.Test)
	g.str("){").indentInc().nl()
	g.Statement(stmt.Body)
	g.indentDec().nl().rune('}')

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
	g.str("){").indentInc().nl()
	g.Statement(stmt.Consequent)
	g.indentDec().nl().rune('}')

	if stmt.Alternate != nil {
		g.str("else{").indentInc().nl()
		g.Statement(stmt.Alternate)
		g.indentDec().nl().rune('}')
	}

	return stmt
}
