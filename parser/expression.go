package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseIdentifier() *ast.Identifier {
	if !p.is(token.IDENTIFIER) {
		p.unexpectedToken()
		return nil
	}

	id := p.currentIdentifier()
	p.next()

	return id
}

func (p *Parser) parseString() *ast.StringLiteral {
	defer p.next()

	return &ast.StringLiteral{
		ExprNode: p.exprNode(),
		Literal:  p.literal,
	}
}

func (p *Parser) parseNumber() *ast.NumberLiteral {
	defer p.next()

	return &ast.NumberLiteral{
		ExprNode: p.exprNode(),
		Literal:  p.literal,
	}
}

func (p *Parser) currentIdentifier() *ast.Identifier {
	return &ast.Identifier{
		ExprNode: p.exprNode(),
		Name:     p.literal,
	}
}

func (p *Parser) parseIdentifierIncludingKeywords() *ast.Identifier {
	if matchIdentifier.MatchString(p.literal) {
		literal := p.literal
		loc := p.loc()

		p.next()

		return &ast.Identifier{
			ExprNode: p.exprNodeAt(loc),
			Name:     literal,
		}
	}

	return nil
}

func (p *Parser) parseRegExpLiteral() *ast.RegExpLiteral {
	loc := p.loc()
	loc.From--
	loc.Col--

	offset := p.chrOffset - 1 // Opening slash already gotten
	if p.is(token.QUOTIENT_ASSIGN) {
		offset -= 1 // =
		loc.From--
		loc.Col--
	}

	pattern, err := p.scanString(offset)
	endOffset := p.chrOffset

	if err == nil {
		pattern = pattern[1 : len(pattern)-1]
	}

	flags := ""
	if !isLineTerminator(p.chr) && !isLineWhiteSpace(p.chr) {
		p.next()

		if p.is(token.IDENTIFIER) { // gim

			flags = p.literal
			p.next()
			endOffset = p.chrOffset - 1
		}
	} else {
		p.next()
	}

	literal := p.src[offset:endOffset]
	loc.End(file.Idx(endOffset))

	return &ast.RegExpLiteral{
		ExprNode: p.exprNodeAt(loc),
		Literal:  literal,
		Pattern:  pattern,
		Flags:    flags,
	}
}

func (p *Parser) parseArgumentList() (argumentList []ast.IExpr, start, end file.Idx) {
	start = p.consumeExpected(token.LEFT_PARENTHESIS)

	for p.until(token.RIGHT_PARENTHESIS) {
		if p.is(token.DOTDOTDOT) {
			loc := p.loc()
			p.consumeExpected(token.DOTDOTDOT)

			argumentList = append(argumentList, &ast.SpreadExpression{
				ExprNode: p.exprNodeAt(loc),
				Value:    p.parseAssignmentExpression(),
			})
		} else {
			argumentList = append(argumentList, p.parseAssignmentExpression())
		}

		p.consumePossible(token.COMMA)
	}

	end = p.consumeExpected(token.RIGHT_PARENTHESIS)
	return
}

func (p *Parser) parseCallExpression(left ast.IExpr, typeArguments []ast.FlowType) ast.IExpr {
	argumentList, _, _ := p.parseArgumentList()

	return &ast.CallExpression{
		ExprNode:      p.exprNodeAt(left.GetLoc()),
		Callee:        left,
		TypeArguments: typeArguments,
		ArgumentList:  argumentList,
	}
}

func (p *Parser) parsePostfixExpression() ast.IExpr {
	operand := p.parseLeftHandSideExpressionAllowCall()

	switch p.token {
	case token.INCREMENT, token.DECREMENT:
		// Make sure there is no line terminator here
		if p.implicitSemicolon {
			break
		}

		loc := p.loc()
		tkn := p.token

		p.next()

		switch operand.(type) {
		case *ast.Identifier, *ast.MemberExpression:
		default:
			p.error(loc, "Invalid left-hand side in assignment")
		}

		return &ast.UnaryExpression{
			ExprNode: p.exprNodeAt(loc),
			Operator: tkn,
			Operand:  operand,
			Postfix:  true,
		}
	}

	return operand
}

func (p *Parser) parseUnaryExpression() ast.IExpr {
	switch p.token {
	case token.PLUS, token.MINUS, token.NOT, token.BITWISE_NOT:
		fallthrough
	case token.DELETE, token.VOID, token.TYPEOF:
		tkn := p.token
		loc := p.loc()

		p.next()

		return &ast.UnaryExpression{
			ExprNode: p.exprNodeAt(loc),
			Operator: tkn,
			Operand:  p.parseUnaryExpression(),
		}
	case token.INCREMENT, token.DECREMENT:
		tkn := p.token
		loc := p.loc()

		p.next()

		operand := p.parseUnaryExpression()
		switch operand.(type) {
		case *ast.Identifier, *ast.MemberExpression:
		default:
			p.error(loc, "Invalid left-hand side in assignment")
		}
		return &ast.UnaryExpression{
			ExprNode: p.exprNodeAt(loc),
			Operator: tkn,
			Operand:  operand,
		}
	}

	return p.parsePostfixExpression()
}

func (p *Parser) parseShiftExpression() ast.IExpr {
	next := p.parseAdditiveExpression
	left := next()

	for p.is(token.SHIFT_LEFT) || p.is(token.SHIFT_RIGHT) ||
		p.is(token.UNSIGNED_SHIFT_RIGHT) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseRelationalExpression() ast.IExpr {
	next := p.parseShiftExpression
	left := next()

	allowIn := p.scope.allowIn
	p.scope.allowIn = true
	defer func() {
		p.scope.allowIn = allowIn
	}()

	switch p.token {
	case token.LESS, token.LESS_OR_EQUAL, token.GREATER, token.GREATER_OR_EQUAL:
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
			ExprNode:   p.exprNodeAt(left.GetLoc()),
			Operator:   tkn,
			Left:       left,
			Right:      p.parseRelationalExpression(),
			Comparison: true,
		}
	case token.INSTANCEOF:
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    p.parseRelationalExpression(),
		}
	case token.IN:
		if !allowIn {
			return left
		}
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    p.parseRelationalExpression(),
		}
	}

	return left
}

func (p *Parser) parseEqualityExpression() ast.IExpr {
	next := p.parseRelationalExpression
	left := next()

	for p.is(token.EQUAL) || p.is(token.NOT_EQUAL) ||
		p.is(token.STRICT_EQUAL) || p.is(token.STRICT_NOT_EQUAL) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode:   p.exprNodeAt(left.GetLoc()),
			Operator:   tkn,
			Left:       left,
			Right:      next(),
			Comparison: true,
		}
	}

	return left
}

func (p *Parser) parseBitwiseAndExpression() ast.IExpr {
	next := p.parseEqualityExpression
	left := next()

	for p.is(token.AND) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseBitwiseExclusiveOrExpression() ast.IExpr {
	next := p.parseBitwiseAndExpression
	left := next()

	for p.is(token.EXCLUSIVE_OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseBitwiseOrExpression() ast.IExpr {
	next := p.parseBitwiseExclusiveOrExpression
	left := next()

	for p.is(token.OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseLogicalAndExpression() ast.IExpr {
	next := p.parseBitwiseOrExpression
	left := next()

	for p.is(token.LOGICAL_AND) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseLogicalOrExpression() ast.IExpr {
	next := p.parseLogicalAndExpression
	left := next()

	for p.is(token.LOGICAL_OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			ExprNode: p.exprNodeAt(left.GetLoc()),
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseExpression() ast.IExpr {
	loc := p.loc()
	next := p.parseAssignmentExpression
	left := next()

	if p.is(token.COMMA) {
		sequence := []ast.IExpr{left}
		for {
			if !p.is(token.COMMA) {
				break
			}
			p.next()

			exp := next()
			loc = loc.Add(exp.GetLoc())
			sequence = append(sequence, exp)
		}
		return &ast.SequenceExpression{
			ExprNode: p.exprNodeAt(loc),
			Sequence: sequence,
		}
	} else if p.is(token.COLON) && p.scope.allowTypeAssertion {
		typeAssertion := p.parseFlowTypeAnnotation()

		return &ast.FlowTypeAssertionExpression{
			ExprNode: p.exprNodeAt(loc),
			Left:     left,
			FlowType: typeAssertion,
		}
	}

	return left
}
