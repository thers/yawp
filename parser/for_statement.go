package parser

import (
	"fmt"
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseIterationStatement() ast.IStmt {
	inIteration := p.scope.inIteration
	p.scope.inIteration = true
	defer func() {
		p.scope.inIteration = inIteration
	}()
	return p.parseStatement()
}

func (p *Parser) parseForIn(loc *file.Loc, into ast.IStmt) *ast.ForInStatement {
	// Already have consumed "<into> in"

	source := p.parseExpression()
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForInStatement{
		StmtNode: p.stmtNodeAt(loc),
		Left:     into,
		Right:    source,
		Body:     p.parseIterationStatement(),
	}
}

func (p *Parser) parseForOf(loc *file.Loc, into ast.IStmt) *ast.ForOfStatement {
	// Already have consumed "<into> of"

	source := p.parseExpression()
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForOfStatement{
		StmtNode: p.stmtNodeAt(loc),
		Left:     into,
		Right:    source,
		Body:     p.parseIterationStatement(),
	}
}

func (p *Parser) parseFor(loc *file.Loc, initializer ast.IStmt) *ast.ForStatement {

	// Already have consumed "<initializer> ;"

	var test, update ast.IExpr

	if !p.is(token.SEMICOLON) {
		test = p.parseExpression()
	}
	p.consumeExpected(token.SEMICOLON)

	if !p.is(token.RIGHT_PARENTHESIS) {
		update = p.parseExpression()
	}
	p.consumeExpected(token.RIGHT_PARENTHESIS)

	return &ast.ForStatement{
		StmtNode:    p.stmtNodeAt(loc),
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        p.parseIterationStatement(),
	}
}

func (p *Parser) maybeParseExpression() ast.IExpr {
	wasAllowIn := p.scope.allowIn
	p.scope.allowIn = false
	defer func() {
		_ = recover()
		p.scope.allowIn = wasAllowIn
	}()

	return p.parseExpression()
}

func (p *Parser) getForKind() (isForIn, isForOf, isFor bool) {
	isForIn = p.is(token.IN)
	isForOf = p.isAllowed(token.OF)
	isFor = p.is(token.SEMICOLON)

	return
}

func (p *Parser) parseForStatement() ast.IStmt {
	start := p.loc()
	p.consumeExpected(token.FOR)
	p.consumeExpected(token.LEFT_PARENTHESIS)

	p.useSymbolsScope(ast.SSTBlock)
	defer p.restoreSymbolsScope()

	// for (;
	if p.is(token.SEMICOLON) {
		p.next()
		return p.parseFor(start, nil)
	}

	// for(const/var/let
	if p.isVariableStatementStart() {
		loc := p.loc()
		kind := p.token

		p.next()

		left := &ast.VariableStatement{
			StmtNode: p.stmtNodeAt(loc),
			Kind:     kind,
			List:     p.parseVariableDeclarationList(kind),
		}

		isForIn, isForOf, isFor := p.getForKind()

		if isForIn || isForOf {
			if len(left.List) > 1 {
				p.error(start, fmt.Sprintf("for-%s can not declare multiple variables", token.IN.String()))

				return nil
			}

			p.next()

			if isForIn {
				return p.parseForIn(start, left)
			} else {
				return p.parseForOf(start, left)
			}
		}

		if isFor {
			p.next()

			return p.parseFor(start, left)
		}

		p.unexpectedToken()

		return nil
	}

	// for (anything
	left := &ast.ExpressionStatement{}
	snapshot := p.snapshot()

	lhs := p.maybeParseExpression()
	if lhs == nil {
		// already found binding pattern it seems
		p.toSnapshot(snapshot)

		switch p.token {
		case token.LEFT_BRACE:
			left.Expression = &ast.VariableBinding{
				ExprNode: p.exprNode(),
				Binder:   p.parseObjectBindingAllowLHS(),
			}
		case token.LEFT_BRACKET:
			left.Expression = &ast.VariableBinding{
				ExprNode: p.exprNode(),
				Binder:   p.parseArrayBindingAllowLHS(),
			}
		default:
			p.unexpectedToken()

			return nil
		}
	} else {
		left.Expression = lhs
	}

	isForIn, isForOf, isFor := p.getForKind()

	// consume in/for/;
	p.next()

	if isForIn || isForOf {
		// it could be that we've parsed id,
		// array literal or object literal
		// when these are in fact binding patterns
		switch leftExp := left.Expression.(type) {
		case *ast.Identifier:
			left.Expression = &ast.VariableBinding{
				ExprNode: leftExp.ExprNode,
				Binder: &ast.IdentifierBinder{
					Id: leftExp,
				},
			}
		case *ast.ArrayLiteral:
			p.toSnapshot(snapshot)
			left.Expression = &ast.VariableBinding{
				ExprNode: leftExp.ExprNode,
				Binder:   p.parseArrayBindingAllowLHS(),
			}
			p.next()
		case *ast.ObjectLiteral:
			p.toSnapshot(snapshot)
			left.Expression = &ast.VariableBinding{
				ExprNode: leftExp.ExprNode,
				Binder:   p.parseObjectBindingAllowLHS(),
			}
			p.next()
		}
	}

	if isForIn {
		return p.parseForIn(start, left)
	}

	if isForOf {
		return p.parseForOf(start, left)
	}

	if isFor {
		return p.parseFor(start, left)
	}

	p.unexpectedToken()

	return nil
}
