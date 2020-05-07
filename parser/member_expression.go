package parser

import (
	"yawp/parser/ast"
)

func (p *Parser) parseMemberExpressionOrIdentifier() ast.Expression {
	loc := p.loc()
	lhs := p.parseLeftHandSideExpressionAllowCall()

	switch lhs.(type) {
	case *ast.Identifier:
	case *ast.MemberExpression:

	default:
		p.error(loc, "Invalid member expression")

		return nil
	}

	return lhs
}
