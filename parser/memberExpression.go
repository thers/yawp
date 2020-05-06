package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseMemberExpression() ast.Expression {
	var ok bool
	var left ast.Expression

	left = p.parseIdentifier()

	for {
		if p.is(token.PERIOD) {
			dotMember := p.parseDotMember(left)
			left, ok = dotMember.(*ast.MemberExpression)

			if !ok {
				p.error(left.GetLoc(), "Invalid member expression dot prefix")

				return nil
			}
		} else if p.is(token.LEFT_BRACKET) {
			bracketMember := p.parseBracketMember(left)
			left, ok = bracketMember.(*ast.MemberExpression)

			if !ok {
				p.error(bracketMember.GetLoc(), "Invalid member expression")

				return nil
			}
		} else {
			break
		}
	}

	return left
}
