package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parsePrimaryExpression() ast.Expression {
	literal := p.literal
	idx := p.idx

	switch p.token {
	case token.SUPER:
		start := p.idx

		if !p.scope.inClass && !p.scope.inFunction {
			p.error(start, "illegal use of super keyword")

			return &ast.BadExpression{
				From: start,
				To:   start,
			}
		}

		p.next()
		arguments, _, end := p.parseArgumentList()
		p.semicolon()

		return &ast.ClassSuperExpression{
			Start:     start,
			End:       end,
			Arguments: arguments,
		}
	case token.CLASS:
		return p.parseClassExpression()
	case token.AWAIT:
		idx := p.idx
		p.next()

		return &ast.AwaitExpression{
			Start:      idx,
			Expression: p.parseAssignmentExpression(),
		}
	case token.ASYNC:
		idx := p.idx
		p.next()

		if p.is(token.FUNCTION) {
			return p.parseFunction(false, idx, true)
		} else {
			return p.tryParseAsyncArrowFunction(idx)
		}
	case token.IDENTIFIER:
		return p.parseIdentifierOrSingleArgumentArrowFunction(false)
	case token.NULL:
		p.next()
		return &ast.NullLiteral{
			Start:   idx,
			Literal: literal,
		}
	case token.BOOLEAN:
		p.next()
		value := false
		switch literal {
		case "true":
			value = true
		case "false":
			value = false
		default:
			p.error(idx, "Illegal boolean literal")
		}
		return &ast.BooleanLiteral{
			Start:   idx,
			Literal: literal,
			Value:   value,
		}
	case token.TEMPLATE_QUOTE:
		return p.parseTemplateExpression()
	case token.STRING:
		p.next()
		value, err := parseStringLiteral(literal[1 : len(literal)-1])
		if err != nil {
			p.error(idx, err.Error())
		}
		return &ast.StringLiteral{
			Start:   idx,
			Literal: literal,
			Value:   value,
		}
	case token.NUMBER:
		p.next()
		value, err := parseNumberLiteral(literal)
		if err != nil {
			p.error(idx, err.Error())
			value = 0
		}
		return &ast.NumberLiteral{
			Start:   idx,
			Literal: literal,
			Value:   value,
		}
	case token.SLASH, token.QUOTIENT_ASSIGN:
		return p.parseRegExpLiteral()
	case token.LEFT_BRACE:
		return p.parseObjectLiteralOrObjectPatternBinding()
	case token.LEFT_BRACKET:
		return p.parseArrayLiteralOrArrayBinding()
	case token.LEFT_PARENTHESIS:
		return p.parseArrowFunctionOrSequenceExpression(false)
	case token.THIS:
		p.next()
		return &ast.ThisExpression{
			Start: idx,
		}
	case token.FUNCTION:
		return p.parseFunction(false, p.idx, false)
	case token.YIELD:
		return p.parseYieldExpression()
	case token.JSX_FRAGMENT_START, token.LESS:
		return p.parseJSX()
	}

	p.errorUnexpectedToken(p.token)
	p.nextStatement()
	return &ast.BadExpression{From: idx, To: p.idx}
}
