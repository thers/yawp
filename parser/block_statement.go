package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	node := &ast.BlockStatement{}
	node.LeftBrace = p.consumeExpected(token.LEFT_BRACE)
	node.List = p.parseStatementList()
	node.RightBrace = p.consumeExpected(token.RIGHT_BRACE)

	return node
}

func (p *Parser) parseBlockStatementOrObjectPatternBinding() ast.Statement {
	start := p.idx
	partialState := p.captureState()

	objectBinding, success := p.maybeParseObjectBinding()

	if success && p.is(token.ASSIGN) {
		p.consumeExpected(token.ASSIGN)

		return &ast.ExpressionStatement{
			Expression: &ast.VariableBinding{
				Start:       start,
				Binder:      objectBinding,
				Initializer: p.parseAssignmentExpression(),
			},
		}
	}

	p.rewindStateTo(partialState)

	return p.parseBlockStatement()
}

