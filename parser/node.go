package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
)

func (p *Parser) node() ast.Node {
	return p.nodeAt(p.loc())
}

func (p *Parser) nodeAt(loc *file.Loc) ast.Node {
	return ast.Node{
		Loc: loc,
	}
}

func (p *Parser) exprNode() ast.ExprNode {
	return p.exprNodeAt(p.loc())
}

func (p *Parser) exprNodeAt(loc *file.Loc) ast.ExprNode {
	return ast.ExprNode{
		Loc: loc,
	}
}

func (p *Parser) stmtNode() ast.StmtNode {
	return p.stmtNodeAt(p.loc())
}

func (p *Parser) stmtNodeAt(loc *file.Loc) ast.StmtNode {
	return ast.StmtNode{
		Loc: loc,
	}
}
