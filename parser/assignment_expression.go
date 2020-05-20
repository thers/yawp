package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseAssignmentExpression() ast.IExpr {
	var operator token.Token

	left := p.parseConditionalExpression()

	switch p.token {
	case token.ASSIGN:
		operator = p.token
	case token.ADD_ASSIGN:
		operator = token.PLUS
	case token.SUBTRACT_ASSIGN:
		operator = token.MINUS
	case token.MULTIPLY_ASSIGN:
		operator = token.MULTIPLY
	case token.QUOTIENT_ASSIGN:
		operator = token.SLASH
	case token.REMAINDER_ASSIGN:
		operator = token.REMAINDER
	case token.EXPONENTIATION_ASSIGN:
		operator = token.EXPONENTIATION_ASSIGN
	case token.AND_ASSIGN:
		operator = token.AND
	case token.AND_NOT_ASSIGN:
		operator = token.AND_NOT
	case token.OR_ASSIGN:
		operator = token.OR
	case token.EXCLUSIVE_OR_ASSIGN:
		operator = token.EXCLUSIVE_OR
	case token.SHIFT_LEFT_ASSIGN:
		operator = token.SHIFT_LEFT
	case token.SHIFT_RIGHT_ASSIGN:
		operator = token.SHIFT_RIGHT
	case token.UNSIGNED_SHIFT_RIGHT_ASSIGN:
		operator = token.UNSIGNED_SHIFT_RIGHT
	}

	if operator != 0 {
		p.next()

		switch l := left.(type) {
		case *ast.Identifier:
			p.symbol(l, ast.SWrite.Add(l.Symbol.Flags), ast.SRUnknown)
		case *ast.MemberExpression:
			// these are the only valid types of left expression in this context
			// pattern binding assignments are handled outside by
			// parseArrayPatternAssignmentOrLiteral
			// and
			// parseObjectPatternAssignmentOrLiteral
		default:
			p.error(left.GetLoc(), "Invalid left-hand side in assignment")
		}

		restoreSymbolFlags := p.useSymbolFlags(ast.SRead)
		right := p.parseAssignmentExpression()
		restoreSymbolFlags()

		return &ast.AssignmentExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}
