package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseArrayLiteral() ast.IExpr {
	var value []ast.IExpr

	loc := p.loc()

	p.consumeExpected(token.LEFT_BRACKET)

	for p.until(token.RIGHT_BRACKET) {
		// [,,,]
		if p.is(token.COMMA) {
			p.next()
			value = append(value, nil)
			continue
		}

		if p.is(token.DOTDOTDOT) {
			p.consumeExpected(token.DOTDOTDOT)

			value = append(value, &ast.ArraySpread{
				Expression: p.parseAssignmentExpression(),
			})
		} else {
			value = append(value, p.parseAssignmentExpression())
		}

		p.consumePossible(token.COMMA)
	}

	loc.End(p.consumeExpected(token.RIGHT_BRACKET))

	return &ast.ArrayLiteral{
		ExprNode: p.exprNodeAt(loc),
		List:     value,
	}
}

func (p *Parser) maybeParseArrayBinding() (*ast.ArrayBinding, bool) {
	wasLeftHandSideAllowed := p.allowPatternBindingLeftHandSideExpressions
	p.allowPatternBindingLeftHandSideExpressions = true

	defer func() {
		_ = recover()

		p.allowPatternBindingLeftHandSideExpressions = wasLeftHandSideAllowed
	}()

	return p.parseArrayBinding(), true
}

func (p *Parser) parseArrayLiteralOrArrayPatternAssignment() ast.IExpr {
	loc := p.loc()
	snapshot := p.snapshot()

	arrayBinding, success := p.maybeParseArrayBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.VariableBinding{
			ExprNode:    p.exprNodeAt(loc),
			Binder:      arrayBinding,
			Initializer: p.parseAssignmentExpression(),
		}
	}

	p.toSnapshot(snapshot)

	return p.parseArrayLiteral()
}

func (p *Parser) parseArrayBindingStatementOrArrayLiteral() ast.IStmt {
	loc := p.loc()
	snapshot := p.snapshot()

	arrayBinding, success := p.maybeParseArrayBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.ExpressionStatement{
			StmtNode: p.stmtNodeAt(loc),
			Expression: &ast.VariableBinding{
				ExprNode:    p.exprNodeAt(loc),
				Binder:      arrayBinding,
				Initializer: p.parseAssignmentExpression(),
			},
		}
	}

	p.toSnapshot(snapshot)

	return &ast.ExpressionStatement{
		StmtNode:   p.stmtNode(),
		Expression: p.parseAssignmentExpression(),
	}
}
