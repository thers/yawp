package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parsePrimaryExpression() ast.Expression {
	literal := p.literal
	loc := p.loc()

	switch p.token {
	case token.SUPER:
		if !p.scope.inClass && !p.scope.inFunction {
			p.error(loc, "illegal use of super keyword")

			return nil
		}

		p.next()
		arguments, _, end := p.parseArgumentList()

		return &ast.ClassSuperExpression{
			Loc:       loc.End(end),
			Arguments: arguments,
		}
	case token.CLASS:
		return p.parseClassExpression()
	case token.AWAIT:
		p.next()

		return &ast.AwaitExpression{
			Loc:        loc,
			Expression: p.parseAssignmentExpression(),
		}
	case token.ASYNC:
		st := p.snapshot()
		p.next()

		if p.is(token.FUNCTION) {
			return p.parseFunction(false, loc, true)
		} else {
			return p.tryParseAsyncArrowFunction(loc, st)
		}
	case token.IDENTIFIER:
		return p.parseIdentifierOrSingleArgumentArrowFunction(false)
	case token.NULL:
		p.next()
		return &ast.NullLiteral{
			Loc:     loc,
			Literal: literal,
		}
	case token.BOOLEAN:
		p.next()

		if literal != ast.LBooleanTrue && literal != ast.LBooleanFalse {
			p.error(loc, "Illegal boolean literal")
		}

		return &ast.BooleanLiteral{
			Loc:     loc,
			Literal: literal,
		}
	case token.TEMPLATE_QUOTE:
		return p.parseTemplateExpression()
	case token.STRING:
		p.next()

		return &ast.StringLiteral{
			Loc:     loc,
			Literal: literal,
		}
	case token.NUMBER:
		p.next()

		return &ast.NumberLiteral{
			Loc:     loc,
			Literal: literal,
		}
	case token.SLASH, token.QUOTIENT_ASSIGN:
		return p.parseRegExpLiteral()
	case token.LEFT_BRACE:
		return p.parseObjectLiteralOrObjectPatternAssignment()
	case token.LEFT_BRACKET:
		return p.parseArrayLiteralOrArrayBinding()
	case token.LEFT_PARENTHESIS:
		return p.parseArrowFunctionOrSequenceExpression(false)
	case token.THIS:
		p.next()
		return &ast.ThisExpression{
			Loc: loc,
		}
	case token.FUNCTION:
		return p.parseFunction(false, loc, false)
	case token.YIELD:
		return p.parseYieldExpression()
	case token.JSX_FRAGMENT_START:
		return p.parseJSXFragment()
	case token.LESS:
		return p.parseJSXElementOrGenericArrowFunction()
	case token.IMPORT:
		return p.parseImportCall()
	}

	p.errorUnexpectedToken(p.token)

	return nil
}
