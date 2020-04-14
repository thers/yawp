package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	var value []ast.Expression

	start := p.consumeExpected(token.LEFT_BRACKET)
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
	end := p.consumeExpected(token.RIGHT_BRACKET)

	return &ast.ArrayLiteral{
		Start: start,
		End:   end,
		Value: value,
	}
}

func (p *Parser) maybeParseArrayBinding() (*ast.ArrayBinding, bool) {
	defer func() {
		err := recover()
		if err != nil {
			return
		}
	}()

	return p.parseArrayBinding(), true
}

func (p *Parser) parseArrayLiteralOrArrayBinding() ast.Expression {
	start := p.loc
	partialState := p.captureState()

	arrayBinding, success := p.maybeParseArrayBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.VariableBinding{
			Start:       start,
			Binder:      arrayBinding,
			Initializer: p.parseAssignmentExpression(),
		}
	}

	p.rewindStateTo(partialState)

	return p.parseArrayLiteral()
}

func (p *Parser) parseArrayBindingStatementOrArrayLiteral() ast.Statement {
	return &ast.ExpressionStatement{
		Expression: p.parseAssignmentExpression(),
	}
}
