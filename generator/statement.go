package generator

import "yawp/parser/ast"


func (g *Generator) BlockStatement(blck *ast.BlockStatement) *ast.BlockStatement {
	for index, stmt := range blck.List {
		if index > 0 {
			g.nl()
		}

		g.Statement(stmt)
	}

	return blck
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
