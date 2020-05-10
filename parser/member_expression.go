package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
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

func (p *Parser) parseDotMember(left ast.Expression) ast.Expression {
	wasAllowIn := p.scope.allowIn
	p.scope.allowIn = false

	p.consumeExpected(token.PERIOD)

	p.scope.allowIn = wasAllowIn

	// this.#bla
	if p.is(token.HASH) && p.scope.inClass {
		if leftThisExp, ok := left.(*ast.ThisExpression); ok {
			p.consumeExpected(token.HASH)

			leftThisExp.Private = true
		} else {
			p.unexpectedToken()

			return nil
		}
	}

	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		p.unexpectedToken()

		return nil
	}

	return &ast.MemberExpression{
		Left:  left,
		Right: identifier,
		Kind:  ast.MKObject,
	}
}

func (p *Parser) parseBracketMember(left ast.Expression) ast.Expression {
	p.consumeExpected(token.LEFT_BRACKET)
	member := p.parseExpression()
	p.consumeExpected(token.RIGHT_BRACKET)

	return &ast.MemberExpression{
		Left:  left,
		Right: member,
		Kind:  ast.MKArray,
	}
}
