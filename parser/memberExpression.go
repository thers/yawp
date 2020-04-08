package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseMemberExpression() ast.MemberExpression {
	var ok bool
	var left ast.MemberExpression

	left = p.parseIdentifier()

	for {
		if p.is(token.PERIOD) {
			dotMember := p.parseDotMember(left)
			left, ok = dotMember.(ast.MemberExpression)

			if !ok {
				return &ast.BadExpression{
					From: left.StartAt(),
					To:   left.EndAt(),
				}
			}
		} else if p.is(token.LEFT_BRACKET) {
			bracketMember := p.parseBracketMember(left)
			left, ok = bracketMember.(ast.MemberExpression)

			if !ok {
				return &ast.BadExpression{
					From: left.StartAt(),
					To:   left.EndAt(),
				}
			}
		} else {
			break
		}
	}

	return left
}
