package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseOptionalExpression(left ast.Expression) ast.Expression {
	optionalChaining := p.consumeExpected(token.OPTIONAL_CHAINING)

	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		switch p.token {
		case token.LEFT_BRACKET:
			p.consumeExpected(token.LEFT_BRACKET)

			index := p.parseAssignmentExpression()
			end := p.consumeExpected(token.RIGHT_BRACKET)

			return &ast.OptionalArrayMemberAccessExpression{
				Left:  left,
				Index: index,
				End:   end,
			}
		case token.LEFT_PARENTHESIS:
			arguments, _, end := p.parseArgumentList()

			return &ast.OptionalCallExpression{
				Left:      left,
				Arguments: arguments,
				End:       end,
			}
		default:
			p.consumeExpected(token.IDENTIFIER)
			p.nextStatement()

			return &ast.BadExpression{
				From: optionalChaining,
				To:   p.idx,
			}
		}
	}

	p.next()

	return &ast.OptionalObjectMemberAccessExpression{
		Left:       left,
		Identifier: identifier,
	}
}
