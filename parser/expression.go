package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseIdentifier() *ast.Identifier {
	literal := p.literal
	idx := p.idx
	p.next()
	return &ast.Identifier{
		Name:  literal,
		Start: idx,
	}
}

func (p *Parser) parseIdentifierIncludingKeywords() *ast.Identifier {
	if matchIdentifier.MatchString(p.literal) {
		literal := p.literal
		idx := p.idx

		p.next()

		return &ast.Identifier{
			Name:  literal,
			Start: idx,
		}
	}

	return nil
}

func (p *Parser) parseRegExpLiteral() *ast.RegExpLiteral {

	offset := p.chrOffset - 1 // Opening slash already gotten
	if p.is(token.QUOTIENT_ASSIGN) {
		offset -= 1 // =
	}
	idx := p.idxOf(offset)

	pattern, err := p.scanString(offset)
	endOffset := p.chrOffset

	if err == nil {
		pattern = pattern[1 : len(pattern)-1]
	}

	flags := ""
	if !isLineTerminator(p.nextChr) && !isLineWhiteSpace(p.nextChr) {
		p.next()

		if p.is(token.IDENTIFIER) { // gim

			flags = p.literal
			p.next()
			endOffset = p.chrOffset - 1
		}
	} else {
		p.next()
	}

	literal := p.str[offset:endOffset]

	return &ast.RegExpLiteral{
		Start:   idx,
		Literal: literal,
		Pattern: pattern,
		Flags:   flags,
	}
}

func (p *Parser) parseArgumentList() (argumentList []ast.Expression, idx0, idx1 file.Idx) {
	idx0 = p.consumeExpected(token.LEFT_PARENTHESIS)
	if !p.is(token.RIGHT_PARENTHESIS) {
		for {
			argumentList = append(argumentList, p.parseAssignmentExpression())
			if !p.is(token.COMMA) {
				break
			}
			p.next()
		}
	}
	idx1 = p.consumeExpected(token.RIGHT_PARENTHESIS)
	return
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	argumentList, idx0, idx1 := p.parseArgumentList()
	return &ast.CallExpression{
		Callee:           left,
		LeftParenthesis:  idx0,
		ArgumentList:     argumentList,
		RightParenthesis: idx1,
	}
}

func (p *Parser) parseDotMember(left ast.Expression) ast.Expression {
	period := p.consumeExpected(token.PERIOD)

	literal := p.literal
	idx := p.idx

	// this.#bla
	if p.is(token.HASH) {
		if leftThisExp, ok := left.(*ast.ThisExpression); ok {
			p.consumeExpected(token.HASH)

			literal = p.literal
			leftThisExp.Private = true
		} else {
			p.unexpectedToken()
			p.next()
		}
	}

	if !matchIdentifier.MatchString(literal) {
		p.consumeExpected(token.IDENTIFIER)
		p.nextStatement()
		return &ast.BadExpression{From: period, To: p.idx}
	}

	p.next()

	return &ast.DotExpression{
		Left: left,
		Identifier: ast.Identifier{
			Start: idx,
			Name:  literal,
		},
	}
}

func (p *Parser) parseBracketMember(left ast.Expression) ast.Expression {
	idx0 := p.consumeExpected(token.LEFT_BRACKET)
	member := p.parseExpression()
	idx1 := p.consumeExpected(token.RIGHT_BRACKET)
	return &ast.BracketExpression{
		LeftBracket:  idx0,
		Left:         left,
		Member:       member,
		RightBracket: idx1,
	}
}

func (p *Parser) parseNewExpression() ast.Expression {
	idx := p.consumeExpected(token.NEW)
	callee := p.parseLeftHandSideExpression()
	node := &ast.NewExpression{
		Start:  idx,
		Callee: callee,
	}
	if p.is(token.LEFT_PARENTHESIS) {
		argumentList, idx0, idx1 := p.parseArgumentList()
		node.ArgumentList = argumentList
		node.LeftParenthesis = idx0
		node.RightParenthesis = idx1
	}
	return node
}

func (p *Parser) parsePostfixExpression() ast.Expression {
	operand := p.parseLeftHandSideExpressionAllowCall()

	switch p.token {
	case token.INCREMENT, token.DECREMENT:
		// Make sure there is no line terminator here
		if p.implicitSemicolon {
			break
		}
		tkn := p.token
		idx := p.idx
		p.next()
		switch operand.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			p.error(idx, "Invalid left-hand side in assignment")
			p.nextStatement()
			return &ast.BadExpression{From: idx, To: p.idx}
		}
		return &ast.UnaryExpression{
			Operator: tkn,
			Start:    idx,
			Operand:  operand,
			Postfix:  true,
		}
	}

	return operand
}

func (p *Parser) parseUnaryExpression() ast.Expression {

	switch p.token {
	case token.PLUS, token.MINUS, token.NOT, token.BITWISE_NOT:
		fallthrough
	case token.DELETE, token.VOID, token.TYPEOF:
		tkn := p.token
		idx := p.idx
		p.next()
		return &ast.UnaryExpression{
			Operator: tkn,
			Start:    idx,
			Operand:  p.parseUnaryExpression(),
		}
	case token.INCREMENT, token.DECREMENT:
		tkn := p.token
		idx := p.idx
		p.next()
		operand := p.parseUnaryExpression()
		switch operand.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			p.error(idx, "Invalid left-hand side in assignment")
			p.nextStatement()
			return &ast.BadExpression{From: idx, To: p.idx}
		}
		return &ast.UnaryExpression{
			Operator: tkn,
			Start:    idx,
			Operand:  operand,
		}
	}

	return p.parsePostfixExpression()
}

func (p *Parser) parseShiftExpression() ast.Expression {
	next := p.parseAdditiveExpression
	left := next()

	for p.is(token.SHIFT_LEFT) || p.is(token.SHIFT_RIGHT) ||
		p.is(token.UNSIGNED_SHIFT_RIGHT) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseRelationalExpression() ast.Expression {
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
			Operator:   tkn,
			Left:       left,
			Right:      p.parseRelationalExpression(),
			Comparison: true,
		}
	case token.INSTANCEOF:
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
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
			Operator: tkn,
			Left:     left,
			Right:    p.parseRelationalExpression(),
		}
	}

	return left
}

func (p *Parser) parseEqualityExpression() ast.Expression {
	next := p.parseRelationalExpression
	left := next()

	for p.is(token.EQUAL) || p.is(token.NOT_EQUAL) ||
		p.is(token.STRICT_EQUAL) || p.is(token.STRICT_NOT_EQUAL) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      next(),
			Comparison: true,
		}
	}

	return left
}

func (p *Parser) parseBitwiseAndExpression() ast.Expression {
	next := p.parseEqualityExpression
	left := next()

	for p.is(token.AND) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseBitwiseExclusiveOrExpression() ast.Expression {
	next := p.parseBitwiseAndExpression
	left := next()

	for p.is(token.EXCLUSIVE_OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseBitwiseOrExpression() ast.Expression {
	next := p.parseBitwiseExclusiveOrExpression
	left := next()

	for p.is(token.OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseLogicalAndExpression() ast.Expression {
	next := p.parseBitwiseOrExpression
	left := next()

	for p.is(token.LOGICAL_AND) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseLogicalOrExpression() ast.Expression {
	next := p.parseLogicalAndExpression
	left := next()

	for p.is(token.LOGICAL_OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseExpression() ast.Expression {
	next := p.parseAssignmentExpression
	left := next()

	if p.is(token.COMMA) {
		sequence := []ast.Expression{left}
		for {
			if !p.is(token.COMMA) {
				break
			}
			p.next()
			sequence = append(sequence, next())
		}
		return &ast.SequenceExpression{
			Sequence: sequence,
		}
	}

	return left
}
