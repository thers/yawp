package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseOptionalExpression(left ast.IExpr) ast.IExpr {
	loc := p.loc()
	p.consumeExpected(token.OPTIONAL_CHAINING)
	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		switch p.token {
		case token.LEFT_BRACKET:
			p.consumeExpected(token.LEFT_BRACKET)

			index := p.parseAssignmentExpression()
			loc.End(p.consumeExpected(token.RIGHT_BRACKET))

			return &ast.OptionalArrayMemberAccessExpression{
				ExprNode: p.exprNodeAt(loc),
				Left:     left,
				Index:    index,
			}
		case token.LEFT_PARENTHESIS:
			arguments, _, end := p.parseArgumentList()
			loc.End(end)

			return &ast.OptionalCallExpression{
				ExprNode:  p.exprNodeAt(loc),
				Left:      left,
				Arguments: arguments,
			}
		default:
			p.error(loc, "Unexpected token")
		}
	}

	return &ast.OptionalObjectMemberAccessExpression{
		ExprNode:   p.exprNodeAt(loc),
		Left:       left,
		Identifier: identifier,
	}
}
